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

type RestaurantRegistrationServiceImpl struct {
	restaurantStore pkgMongo.Storer
	userStore       pkgMongo.Storer
}

func NewRestaurantRegistrationServiceImpl(restaurantStore pkgMongo.Storer, userStore pkgMongo.Storer) *RestaurantRegistrationServiceImpl {
	return &RestaurantRegistrationServiceImpl{
		restaurantStore: restaurantStore,
		userStore:       userStore,
	}
}

func (s *RestaurantRegistrationServiceImpl) RegisterRestaurant(ctx context.Context, request pkgRestaurants.RegisterRestaurantRequest) (*pkgRestaurants.RegisterRestaurantResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for RegisterRestaurant: %w", err)
	}

	// Fetch the current user.
	var user users.User
	if err := s.userStore.Get(ctx, request.UserID, &user); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get user")
		return nil, fmt.Errorf("RegisterRestaurant failed: %w", err)
	}

	hasRole := slices.Contains(user.Roles, users.RoleRestaurantOwner)

	// Check if a restaurant already exists for this owner.
	var existing []pkgRestaurants.Restaurant
	_, err := s.restaurantStore.List(ctx, bson.M{"ownerID": request.UserID}, 0, 1, &existing)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to check existing restaurants")
		return nil, fmt.Errorf("RegisterRestaurant failed: %w", err)
	}

	// Both restaurant and role exist — already fully registered.
	if len(existing) > 0 && hasRole {
		return nil, fmt.Errorf("user already has a registered restaurant: %w", errs.ErrConflict)
	}

	// Restaurant exists but role is missing — recover from a previous partial failure.
	if len(existing) > 0 && !hasRole {
		if err := s.addRestaurantOwnerRole(ctx, user); err != nil {
			return nil, err
		}
		return &pkgRestaurants.RegisterRestaurantResponse{
			Restaurant: existing[0],
		}, nil
	}

	// Create the restaurant.
	restaurant := pkgRestaurants.Restaurant{
		ID:      uuid.New().String(),
		OwnerID: request.UserID,
		Name:    request.Name,
		City:    request.City,
		Image:   request.Image,
	}

	_, err = s.restaurantStore.Put(ctx, restaurant)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to store restaurant")
		return nil, fmt.Errorf("RegisterRestaurant failed: %w", err)
	}

	// Update the user's roles to include RestaurantOwner.
	if err := s.addRestaurantOwnerRole(ctx, user); err != nil {
		return nil, err
	}

	return &pkgRestaurants.RegisterRestaurantResponse{
		Restaurant: restaurant,
	}, nil
}

func (s *RestaurantRegistrationServiceImpl) addRestaurantOwnerRole(ctx context.Context, user users.User) error {
	newRoles := append(user.Roles, users.RoleRestaurantOwner)
	var updatedUser users.User
	if err := s.userStore.Update(ctx, user.ID, bson.M{"roles": newRoles}, &updatedUser); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to update user roles")
		return fmt.Errorf("RegisterRestaurant failed: %w", err)
	}
	return nil
}
