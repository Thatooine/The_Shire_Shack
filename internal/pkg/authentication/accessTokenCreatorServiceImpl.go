package authentication

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
	"github.com/go-jose/go-jose/v4"
	"github.com/rs/zerolog/log"
)

// AccessTokenCreatorServiceImpl creates signed JWT access tokens using go-jose.
type AccessTokenCreatorServiceImpl struct {
	tokenSigner jose.Signer
}

// NewAccessTokenCreatorServiceImpl returns a new AccessTokenCreatorServiceImpl
// with the provided jose.Signer for signing tokens.
func NewAccessTokenCreatorServiceImpl(tokenSigner jose.Signer) *AccessTokenCreatorServiceImpl {
	return &AccessTokenCreatorServiceImpl{
		tokenSigner: tokenSigner,
	}
}

// CreateAccessToken marshals the login claims, signs the payload, and returns
// the compact-serialized JWT as the access token.
func (a *AccessTokenCreatorServiceImpl) CreateAccessToken(ctx context.Context, request authentication.CreateAccessTokenRequest) (*authentication.CreateAccessTokenResponse, error) {
	// marshal claims to JSON
	claimsPayload, err := json.Marshal(request.LoginClaim)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("could not marshal claims for token")
		return nil, fmt.Errorf("CreateAccessToken failed: could not marshal claims: %w", err)
	}

	// sign the marshalled payload
	signedObj, err := a.tokenSigner.Sign(claimsPayload)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("could not sign payload")
		return nil, fmt.Errorf("CreateAccessToken failed: could not sign payload: %w", err)
	}

	// serialize the signed object into a compact JWT string
	signedJWT, err := signedObj.CompactSerialize()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("could not serialize signed token")
		return nil, fmt.Errorf("CreateAccessToken failed: could not serialize token: %w", err)
	}

	return &authentication.CreateAccessTokenResponse{
		AccessToken: signedJWT,
	}, nil
}
