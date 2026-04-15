package users

import (
	"context"

	pkgUsers "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
)

type UserCreatorServiceMock struct {
	CreateUserFn func(ctx context.Context, request pkgUsers.CreateUserRequest) (*pkgUsers.CreateUserResponse, error)
}

func (m *UserCreatorServiceMock) CreateUser(ctx context.Context, request pkgUsers.CreateUserRequest) (*pkgUsers.CreateUserResponse, error) {
	return m.CreateUserFn(ctx, request)
}
