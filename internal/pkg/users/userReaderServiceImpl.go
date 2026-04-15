package users

import (
	"context"
	"fmt"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// UserReaderServiceImpl implements users.UserReaderServiceImpl using a mongo Store for the "users" collection.
type UserReaderServiceImpl struct {
	store *mongo.Store
}

func NewUserReaderServiceImpl(store *mongo.Store) *UserReaderServiceImpl {
	return &UserReaderServiceImpl{
		store: store,
	}
}

// GetUser fetches a single user by Email.
func (u *UserReaderServiceImpl) GetUser(ctx context.Context, request users.GetUserRequest) (*users.GetUserResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for GetUser: %w", err)
	}

	filter := bson.M{"email": request.Email}
	var userList []users.User
	_, err := u.store.List(ctx, filter, 0, 1, &userList)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get user from store")
		return nil, fmt.Errorf("GetUser failed: %w", err)
	}

	if len(userList) == 0 {
		return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
	}

	return &users.GetUserResponse{
		User: userList[0],
	}, nil
}

// ListUsers returns a paginated list of users.
func (u *UserReaderServiceImpl) ListUsers(ctx context.Context, request users.ListUsersRequest) (*users.ListUsersResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for ListUsers: %w", err)
	}

	filter := bson.M{}

	var userList []users.User
	total, err := u.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &userList)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to list users from store")
		return nil, fmt.Errorf("ListUsers failed: %w", err)
	}

	return &users.ListUsersResponse{
		Users: userList,
		Total: total,
	}, nil
}

// SearchUsers performs a case-insensitive regex search across user name and email.
func (u *UserReaderServiceImpl) SearchUsers(ctx context.Context, request users.SearchUsersRequest) (*users.SearchUsersResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for SearchUsers: %w", err)
	}

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": request.Query, "$options": "i"}},
			{"email": bson.M{"$regex": request.Query, "$options": "i"}},
		},
	}

	var userList []users.User
	total, err := u.store.List(ctx, filter, int64(request.Offset), int64(request.Limit), &userList)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to search users from store")
		return nil, fmt.Errorf("SearchUsers failed: %w", err)
	}

	return &users.SearchUsersResponse{
		Users: userList,
		Total: total,
	}, nil
}
