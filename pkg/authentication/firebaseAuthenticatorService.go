package authentication

import "context"

// FirebaseAuthenticatorService defines the interface for authenticating a user via a Firebase token.
// The flow: receive a Firebase ID token, verify it against the Firebase client to retrieve the user,
// then produce and return a JWT containing the user's claims.
type FirebaseAuthenticatorService interface {
	AuthenticateWithFirebaseToken(ctx context.Context, request FirebaseAuthRequest) (*FirebaseAuthResponse, error)
}

type FirebaseAuthRequest struct {
	FirebaseToken string
}

type FirebaseAuthResponse struct {
	Token  string
	UserID string
	Email  string
}
