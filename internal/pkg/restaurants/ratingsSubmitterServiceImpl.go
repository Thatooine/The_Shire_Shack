package restaurants

import (
	"context"
	"fmt"
	"time"

	pkgMongo "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	pkgRestaurants "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RatingSubmitterServiceImpl struct {
	store *pkgMongo.Store
}

func NewRatingSubmitterServiceImpl(store *pkgMongo.Store) *RatingSubmitterServiceImpl {
	return &RatingSubmitterServiceImpl{
		store: store,
	}
}

// SubmitRating validates and stores a new rating for a dish.
func (r *RatingSubmitterServiceImpl) SubmitRating(ctx context.Context, request pkgRestaurants.SubmitRatingRequest) (*pkgRestaurants.SubmitRatingResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for SubmitRating: %w", err)
	}

	rating := pkgRestaurants.Rating{
		ID:        uuid.New().String(),
		DishID:    request.DishID,
		UserID:    request.UserID,
		Score:     request.Score,
		Review:    request.Review,
		CreatedAt: time.Now(),
	}

	id, err := r.store.Put(ctx, rating)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to store rating")
		return nil, fmt.Errorf("SubmitRating failed: %w", err)
	}

	rating.ID = id
	return &pkgRestaurants.SubmitRatingResponse{
		Rating: rating,
	}, nil
}
