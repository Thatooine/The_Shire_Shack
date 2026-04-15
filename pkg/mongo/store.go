package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Storer defines the interface for persistent storage operations.
type Storer interface {
	Put(ctx context.Context, document interface{}) (string, error)
	Get(ctx context.Context, id string, result interface{}) error
	List(ctx context.Context, filter bson.M, offset int64, limit int64, results interface{}) (int64, error)
	Update(ctx context.Context, id string, update bson.M, result interface{}) error
	Remove(ctx context.Context, id string) error
}

// Store provides persistent storage operations for a single MongoDB collection.
type Store struct {
	collection *mongo.Collection
}

// NewStore returns a new Store for the given database and collection name.
func NewStore(client *mongo.Client, database string, collection string) *Store {
	return &Store{
		collection: client.Database(database).Collection(collection),
	}
}

// Put inserts a single document and returns the inserted ID.
func (s *Store) Put(ctx context.Context, document interface{}) (string, error) {
	result, err := s.collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to store document: %w", err)
	}

	id, ok := result.InsertedID.(string)
	if !ok {
		return fmt.Sprintf("%v", result.InsertedID), nil
	}
	return id, nil
}

// Get retrieves a single document by its _id field. Pass a pointer to decode into.
func (s *Store) Get(ctx context.Context, id string, result interface{}) error {
	filter := bson.M{"id": id}

	if err := s.collection.FindOne(ctx, filter).Decode(result); err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	return nil
}

// List retrieves documents matching the given filter with pagination.
// Pass a pointer to a slice to decode into.
func (s *Store) List(ctx context.Context, filter bson.M, offset int64, limit int64, results interface{}) (int64, error) {
	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to list documents: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return 0, fmt.Errorf("failed to decode documents: %w", err)
	}

	return total, nil
}

// Update updates a single document by its _id field with the given fields.
// Pass a pointer to decode the updated document into.
func (s *Store) Update(ctx context.Context, id string, update bson.M, result interface{}) error {
	filter := bson.M{"id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := s.collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, opts).Decode(result); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	return nil
}

// EnsureIndex creates an ascending index on the given field if it does not already exist.
// If unique is true, the index enforces that no two documents share the same field value.
func (s *Store) EnsureIndex(ctx context.Context, field string, unique bool) error {
	model := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(unique),
	}
	_, err := s.collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to create index on %q: %w", field, err)
	}
	return nil
}

// Remove deletes a single document by its _id field.
func (s *Store) Remove(ctx context.Context, id string) error {
	filter := bson.M{"id": id}

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to remove document: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("document with id %s not found", id)
	}
	return nil
}
