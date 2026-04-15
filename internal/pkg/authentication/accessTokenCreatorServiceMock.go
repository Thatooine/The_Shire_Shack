package authentication

import (
	"context"

	pkgAuth "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
)

type AccessTokenCreatorServiceMock struct {
	CreateAccessTokenFn func(ctx context.Context, request pkgAuth.CreateAccessTokenRequest) (*pkgAuth.CreateAccessTokenResponse, error)
}

func (m *AccessTokenCreatorServiceMock) CreateAccessToken(ctx context.Context, request pkgAuth.CreateAccessTokenRequest) (*pkgAuth.CreateAccessTokenResponse, error) {
	return m.CreateAccessTokenFn(ctx, request)
}
