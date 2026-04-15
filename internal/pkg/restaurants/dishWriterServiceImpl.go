package restaurants

import (
	"context"
	"fmt"
	"slices"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	pkgMongo "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	pkgRestaurants "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DishWriterServiceImpl struct {
	store           *pkgMongo.Store
	userStore       *pkgMongo.Store
	restaurantStore *pkgMongo.Store
}

func NewDishWriterServiceImpl(store *pkgMongo.Store, userStore *pkgMongo.Store, restaurantStore *pkgMongo.Store) *DishWriterServiceImpl {
	return &DishWriterServiceImpl{
		store:           store,
		userStore:       userStore,
		restaurantStore: restaurantStore,
	}
}

// verifyRestaurantOwnership checks that the user has the RestaurantOwner role
// and owns the restaurant identified by restaurantID.
func (d *DishWriterServiceImpl) verifyRestaurantOwnership(ctx context.Context, userID string, restaurantID string) error {
	var user users.User
	if err := d.userStore.Get(ctx, userID, &user); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get user")
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !slices.Contains(user.Roles, users.RoleRestaurantOwner) {
		return fmt.Errorf("user is not a restaurant owner: %w", errs.ErrForbidden)
	}

	var restaurant pkgRestaurants.Restaurant
	if err := d.restaurantStore.Get(ctx, restaurantID, &restaurant); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get restaurant")
		return fmt.Errorf("failed to get restaurant: %w", err)
	}

	if restaurant.OwnerID != userID {
		return fmt.Errorf("user does not own this restaurant: %w", errs.ErrForbidden)
	}

	return nil
}

// CreateDish validates the request, verifies restaurant ownership, then inserts a new dish.
func (d *DishWriterServiceImpl) CreateDish(ctx context.Context, request pkgRestaurants.CreateDishRequest) (*pkgRestaurants.CreateDishResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for CreateDish: %w", err)
	}

	if err := d.verifyRestaurantOwnership(ctx, request.UserID, request.RestaurantID); err != nil {
		return nil, err
	}

	dish := pkgRestaurants.Dish{
		ID:           uuid.New().String(),
		Name:         request.Name,
		Description:  request.Description,
		Price:        request.Price,
		RestaurantID: request.RestaurantID,
		Image:        request.Image,
	}

	_, err := d.store.Put(ctx, dish)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to store dish")
		return nil, fmt.Errorf("CreateDish failed: %w", err)
	}

	return &pkgRestaurants.CreateDishResponse{
		Dish: dish,
	}, nil
}

// UpdateDish validates the request, verifies the user owns the dish's restaurant,
// then updates the dish.
func (d *DishWriterServiceImpl) UpdateDish(ctx context.Context, request pkgRestaurants.UpdateDishRequest) (*pkgRestaurants.UpdateDishResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for UpdateDish: %w", err)
	}

	// Look up the dish to find its restaurant.
	var existing pkgRestaurants.Dish
	if err := d.store.Get(ctx, request.ID, &existing); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get dish")
		return nil, fmt.Errorf("UpdateDish failed: %w", err)
	}

	if err := d.verifyRestaurantOwnership(ctx, request.UserID, existing.RestaurantID); err != nil {
		return nil, err
	}

	update := bson.M{
		"name":        request.Name,
		"description": request.Description,
		"price":       request.Price,
		"image":       request.Image,
	}

	var dish pkgRestaurants.Dish
	if err := d.store.Update(ctx, request.ID, update, &dish); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to update dish in store")
		return nil, fmt.Errorf("UpdateDish failed: %w", err)
	}

	return &pkgRestaurants.UpdateDishResponse{
		Dish: dish,
	}, nil
}

// DeleteDish validates the request, verifies the user owns the dish's restaurant,
// then removes the dish.
func (d *DishWriterServiceImpl) DeleteDish(ctx context.Context, request pkgRestaurants.DeleteDishRequest) error {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return fmt.Errorf("invalid request for DeleteDish: %w", err)
	}

	// Look up the dish to find its restaurant.
	var existing pkgRestaurants.Dish
	if err := d.store.Get(ctx, request.ID, &existing); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get dish")
		return fmt.Errorf("DeleteDish failed: %w", err)
	}

	if err := d.verifyRestaurantOwnership(ctx, request.UserID, existing.RestaurantID); err != nil {
		return err
	}

	if err := d.store.Remove(ctx, request.ID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to remove dish from store")
		return fmt.Errorf("DeleteDish failed: %w", err)
	}

	return nil
}
