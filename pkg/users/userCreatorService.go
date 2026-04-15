package users

import "context"

// UserCreatorService defines the interface for creating a user.
type UserCreatorService interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, error)
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type CreateUserResponse struct {
	User User
}
