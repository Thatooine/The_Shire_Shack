package restaurants

import "context"

// DishReaderService defines the interface for viewing, listing, and searching dishes.
type DishReaderService interface {
	GetDish(ctx context.Context, request GetDishRequest) (*GetDishResponse, error)
	ListDishes(ctx context.Context, request ListDishesRequest) (*ListDishesResponse, error)
	SearchDishes(ctx context.Context, request SearchDishesRequest) (*SearchDishesResponse, error)
}

// GetDish

type GetDishRequest struct {
	ID string
}

type GetDishResponse struct {
	Dish Dish
}

// ListDishes

type ListDishesRequest struct {
	RestaurantID string
	Offset       int
	Limit        int
}

type ListDishesResponse struct {
	Dishes []Dish
	Total  int64
}

// SearchDishes

type SearchDishesRequest struct {
	Query  string
	Offset int
	Limit  int
}

type SearchDishesResponse struct {
	Dishes []Dish
	Total  int64
}
