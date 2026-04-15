package restaurants

import "context"

// RestaurantReaderService defines the interface for viewing, listing, and searching restaurants.
type RestaurantReaderService interface {
	GetRestaurant(ctx context.Context, request GetRestaurantRequest) (*GetRestaurantResponse, error)
	GetMyRestaurant(ctx context.Context, request GetMyRestaurantRequest) (*GetRestaurantResponse, error)
	ListRestaurants(ctx context.Context, request ListRestaurantsRequest) (*ListRestaurantsResponse, error)
	SearchRestaurants(ctx context.Context, request SearchRestaurantsRequest) (*SearchRestaurantsResponse, error)
}

// GetRestaurant

type GetRestaurantRequest struct {
	ID string
}

type GetRestaurantResponse struct {
	Restaurant Restaurant
}

// GetMyRestaurant

type GetMyRestaurantRequest struct {
	OwnerID string
}

// ListRestaurants

type ListRestaurantsRequest struct {
	Offset int
	Limit  int
}

type ListRestaurantsResponse struct {
	Restaurants []Restaurant
	Total       int64
}

// SearchRestaurants

type SearchRestaurantsRequest struct {
	Query  string
	Offset int
	Limit  int
}

type SearchRestaurantsResponse struct {
	Restaurants []Restaurant
	Total       int64
}
