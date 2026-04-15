package restaurants

import (
	"context"
	"fmt"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// RestaurantReaderServiceImpl implements restaurants.RestaurantReaderService
// using a mongo Store for the "restaurants" collection.
type RestaurantReaderServiceImpl struct {
	store *mongo.Store
}

func NewRestaurantReaderServiceImpl(store *mongo.Store) *RestaurantReaderServiceImpl {
	return &RestaurantReaderServiceImpl{
		store: store,
	}
}

// GetRestaurant fetches a single restaurant by ID.
func (r *RestaurantReaderServiceImpl) GetRestaurant(ctx context.Context, request restaurants.GetRestaurantRequest) (*restaurants.GetRestaurantResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for GetRestaurant: %w", err)
	}

	var restaurant restaurants.Restaurant
	if err := r.store.Get(ctx, request.ID, &restaurant); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get restaurant from store")
		return nil, fmt.Errorf("GetRestaurant failed: %w", err)
	}

	return &restaurants.GetRestaurantResponse{
		Restaurant: restaurant,
	}, nil
}

// GetMyRestaurant fetches the restaurant owned by the given user.
func (r *RestaurantReaderServiceImpl) GetMyRestaurant(ctx context.Context, request restaurants.GetMyRestaurantRequest) (*restaurants.GetRestaurantResponse, error) {
	filter := bson.M{"ownerID": request.OwnerID}
	var result []restaurants.Restaurant
	_, err := r.store.List(ctx, filter, 0, 1, &result)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get restaurant from store")
		return nil, fmt.Errorf("GetMyRestaurant failed: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no restaurant found for this user: %w", errs.ErrNotFound)
	}

	return &restaurants.GetRestaurantResponse{
		Restaurant: result[0],
	}, nil
}

// ListRestaurants returns a paginated list of restaurants.
func (r *RestaurantReaderServiceImpl) ListRestaurants(ctx context.Context, request restaurants.ListRestaurantsRequest) (*restaurants.ListRestaurantsResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for ListRestaurants: %w", err)
	}

	filter := bson.M{}

	var result []restaurants.Restaurant
	total, err := r.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &result)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to list restaurants from store")
		return nil, fmt.Errorf("ListRestaurants failed: %w", err)
	}

	return &restaurants.ListRestaurantsResponse{
		Restaurants: result,
		Total:       total,
	}, nil
}

// SearchRestaurants performs a case-insensitive regex search across restaurant name and city.
func (r *RestaurantReaderServiceImpl) SearchRestaurants(ctx context.Context, request restaurants.SearchRestaurantsRequest) (*restaurants.SearchRestaurantsResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for SearchRestaurants: %w", err)
	}

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": request.Query, "$options": "i"}},
			{"city": bson.M{"$regex": request.Query, "$options": "i"}},
		},
	}

	var result []restaurants.Restaurant
	total, err := r.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &result)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to search restaurants from store")
		return nil, fmt.Errorf("SearchRestaurants failed: %w", err)
	}

	return &restaurants.SearchRestaurantsResponse{
		Restaurants: result,
		Total:       total,
	}, nil
}
