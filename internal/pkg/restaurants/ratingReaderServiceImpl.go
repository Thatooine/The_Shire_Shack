package restaurants

import (
	"context"
	"fmt"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RatingReaderServiceImpl struct {
	store *mongo.Store
}

func NewRatingReaderServiceImpl(store *mongo.Store) *RatingReaderServiceImpl {
	return &RatingReaderServiceImpl{
		store: store,
	}
}

// ListRatings returns a paginated list of ratings for a given dish.
func (r *RatingReaderServiceImpl) ListRatings(ctx context.Context, request restaurants.ListRatingsRequest) (*restaurants.ListRatingsResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for ListRatings: %w", err)
	}

	filter := bson.M{"dish_id": request.DishID}

	var ratings []restaurants.Rating
	total, err := r.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &ratings)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to list ratings from store")
		return nil, fmt.Errorf("ListRatings failed: %w", err)
	}

	return &restaurants.ListRatingsResponse{
		Ratings: ratings,
		Total:   total,
	}, nil
}
