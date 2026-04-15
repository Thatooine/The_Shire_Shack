package restaurants

import "context"

// RatingReaderService defines the interface for listing ratings for a dish.
type RatingReaderService interface {
	ListRatings(ctx context.Context, request ListRatingsRequest) (*ListRatingsResponse, error)
}

// ListRatings

type ListRatingsRequest struct {
	DishID string
	Offset int
	Limit  int
}

type ListRatingsResponse struct {
	Ratings []Rating
	Total   int64
}
