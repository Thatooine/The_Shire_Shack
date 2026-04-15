package restaurants

import "context"

// RatingSubmitterService defines the interface for submitting a rating for a dish.
type RatingSubmitterService interface {
	SubmitRating(ctx context.Context, request SubmitRatingRequest) (*SubmitRatingResponse, error)
}

// SubmitRating

type SubmitRatingRequest struct {
	DishID string
	UserID string
	Score  int
	Review string
}

type SubmitRatingResponse struct {
	Rating Rating
}
