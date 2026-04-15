package authentication

import "context"

type AccessTokenValidatorService interface {
	ValidateAccessToken(ctx context.Context, request ValidateAccessTokenRequest) (*ValidateAccessTokenResponse, error)
}

type ValidateAccessTokenRequest struct {
	AccessToken string
}

type ValidateAccessTokenResponse struct {
	LoginClaim LoginClaim
}
