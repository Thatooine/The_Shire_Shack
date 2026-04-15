package users

import (
	"context"

	pkgUsers "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
)

type UserReaderServiceMock struct {
	GetUserFn     func(ctx context.Context, request pkgUsers.GetUserRequest) (*pkgUsers.GetUserResponse, error)
	ListUsersFn   func(ctx context.Context, request pkgUsers.ListUsersRequest) (*pkgUsers.ListUsersResponse, error)
	SearchUsersFn func(ctx context.Context, request pkgUsers.SearchUsersRequest) (*pkgUsers.SearchUsersResponse, error)
}

func (m *UserReaderServiceMock) GetUser(ctx context.Context, request pkgUsers.GetUserRequest) (*pkgUsers.GetUserResponse, error) {
	return m.GetUserFn(ctx, request)
}

func (m *UserReaderServiceMock) ListUsers(ctx context.Context, request pkgUsers.ListUsersRequest) (*pkgUsers.ListUsersResponse, error) {
	return m.ListUsersFn(ctx, request)
}

func (m *UserReaderServiceMock) SearchUsers(ctx context.Context, request pkgUsers.SearchUsersRequest) (*pkgUsers.SearchUsersResponse, error) {
	return m.SearchUsersFn(ctx, request)
}
