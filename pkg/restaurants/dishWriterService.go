package restaurants

import "context"

// DishWriterService defines the interface for creating, updating, and deleting dishes.
type DishWriterService interface {
	CreateDish(ctx context.Context, request CreateDishRequest) (*CreateDishResponse, error)
	UpdateDish(ctx context.Context, request UpdateDishRequest) (*UpdateDishResponse, error)
	DeleteDish(ctx context.Context, request DeleteDishRequest) error
}

// CreateDish

type CreateDishRequest struct {
	UserID       string
	Name         string
	Description  string
	Price        float64
	RestaurantID string
	Image        string
}

type CreateDishResponse struct {
	Dish Dish
}

// UpdateDish

type UpdateDishRequest struct {
	UserID      string
	ID          string
	Name        string
	Description string
	Price       float64
	Image       string
}

type UpdateDishResponse struct {
	Dish Dish
}

// DeleteDish

type DeleteDishRequest struct {
	UserID string
	ID     string
}
