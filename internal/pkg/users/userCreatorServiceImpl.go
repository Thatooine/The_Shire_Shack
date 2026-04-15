package users

import (
	"context"
	"fmt"

	pkgMongo "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
	"github.com/rs/zerolog/log"
)

type UserCreatorServiceImpl struct {
	store *pkgMongo.Store
}

func NewUserCreatorServiceImpl(store *pkgMongo.Store) *UserCreatorServiceImpl {
	return &UserCreatorServiceImpl{
		store: store,
	}
}

func (r *UserCreatorServiceImpl) CreateUser(ctx context.Context, request users.CreateUserRequest) (*users.CreateUserResponse, error) {
	if err := request.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("request validation failed")
		return nil, fmt.Errorf("invalid request for CreateUser: %w", err)
	}

	user := users.User{
		ID:    users.NewID(),
		Name:  request.Name,
		Email: request.Email,
		Roles: []users.Role{users.RoleCustomer},
	}

	_, err := r.store.Put(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to store user")
		return nil, fmt.Errorf("CreateUser failed: %w", err)
	}

	return &users.CreateUserResponse{
		User: user,
	}, nil
}
