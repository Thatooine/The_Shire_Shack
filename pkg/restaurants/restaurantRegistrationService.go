package restaurants

import "context"

// RestaurantRegistrationService defines the interface for registering a restaurant.
// It creates the restaurant and promotes the calling user to RestaurantOwner.
type RestaurantRegistrationService interface {
	RegisterRestaurant(ctx context.Context, request RegisterRestaurantRequest) (*RegisterRestaurantResponse, error)
}

type RegisterRestaurantRequest struct {
	UserID string
	Name   string
	City   string
	Image  string
}

type RegisterRestaurantResponse struct {
	Restaurant Restaurant
}
