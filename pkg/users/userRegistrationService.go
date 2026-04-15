package users

import "context"

// UserRegistrationService defines the interface for registering new users.
type UserRegistrationService interface {
	RegisterWithEmailAndPassword(ctx context.Context, request RegisterWithEmailAndPasswordRequest) (*RegisterResponse, error)
	RegisterWithFirebaseToken(ctx context.Context, request RegisterWithFirebaseTokenRequest) (*RegisterResponse, error)
}

type RegisterWithEmailAndPasswordRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterWithFirebaseTokenRequest struct {
	Name          string
	FirebaseToken string
}

type RegisterResponse struct {
	Token  string
	UserID string
	Email  string
}
