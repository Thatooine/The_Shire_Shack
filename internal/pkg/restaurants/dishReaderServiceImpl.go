package restaurants

import (
	"context"
	"fmt"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// DishReaderServiceImpl implements pkgRestaurants.DishReaderService
// using a mongo Store for the "dishes" collection.
type DishReaderServiceImpl struct {
	store *mongo.Store
}

func NewDishReaderServiceImpl(store *mongo.Store) *DishReaderServiceImpl {
	return &DishReaderServiceImpl{
		store: store,
	}
}

// GetDish fetches a single dish by ID.
func (d *DishReaderServiceImpl) GetDish(ctx context.Context, request restaurants.GetDishRequest) (*restaurants.GetDishResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for GetDish: %w", err)
	}

	var dish restaurants.Dish
	if err := d.store.Get(ctx, request.ID, &dish); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get dish from store")
		return nil, fmt.Errorf("GetDish failed: %w", err)
	}

	return &restaurants.GetDishResponse{
		Dish: dish,
	}, nil
}

// ListDishes returns a paginated list of dishes, optionally filtered by restaurant.
func (d *DishReaderServiceImpl) ListDishes(ctx context.Context, request restaurants.ListDishesRequest) (*restaurants.ListDishesResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for ListDishes: %w", err)
	}

	// Only filter by restaurant if provided
	filter := bson.M{}
	if request.RestaurantID != "" {
		filter["restaurant_id"] = request.RestaurantID
	}

	var dishes []restaurants.Dish
	total, err := d.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &dishes)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to list dishes from store")
		return nil, fmt.Errorf("ListDishes failed: %w", err)
	}

	return &restaurants.ListDishesResponse{
		Dishes: dishes,
		Total:  total,
	}, nil
}

// SearchDishes performs a case-insensitive regex search across dish name and description.
func (d *DishReaderServiceImpl) SearchDishes(ctx context.Context, request restaurants.SearchDishesRequest) (*restaurants.SearchDishesResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for SearchDishes: %w", err)
	}

	// Match dishes where name or description contains the query (case-insensitive)
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": request.Query, "$options": "i"}},
			{"description": bson.M{"$regex": request.Query, "$options": "i"}},
		},
	}

	var dishes []restaurants.Dish
	total, err := d.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &dishes)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to search dishes from store")
		return nil, fmt.Errorf("SearchDishes failed: %w", err)
	}

	return &restaurants.SearchDishesResponse{
		Dishes: dishes,
		Total:  total,
	}, nil
}
