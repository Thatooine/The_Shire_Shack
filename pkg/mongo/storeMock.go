package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// StoreMock is a mock implementation of the Storer interface for testing.
type StoreMock struct {
	PutFn    func(ctx context.Context, document interface{}) (string, error)
	GetFn    func(ctx context.Context, id string, result interface{}) error
	ListFn   func(ctx context.Context, filter bson.M, offset int64, limit int64, results interface{}) (int64, error)
	UpdateFn func(ctx context.Context, id string, update bson.M, result interface{}) error
	RemoveFn func(ctx context.Context, id string) error
}

func (m *StoreMock) Put(ctx context.Context, document interface{}) (string, error) {
	return m.PutFn(ctx, document)
}

func (m *StoreMock) Get(ctx context.Context, id string, result interface{}) error {
	return m.GetFn(ctx, id, result)
}

func (m *StoreMock) List(ctx context.Context, filter bson.M, offset int64, limit int64, results interface{}) (int64, error) {
	return m.ListFn(ctx, filter, offset, limit, results)
}

func (m *StoreMock) Update(ctx context.Context, id string, update bson.M, result interface{}) error {
	return m.UpdateFn(ctx, id, update, result)
}

func (m *StoreMock) Remove(ctx context.Context, id string) error {
	return m.RemoveFn(ctx, id)
}
