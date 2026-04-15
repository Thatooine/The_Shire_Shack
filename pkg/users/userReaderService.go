package users

import "context"

// UserReaderService defines the interface for viewing, listing, and searching users.
type UserReaderService interface {
	GetUser(ctx context.Context, request GetUserRequest) (*GetUserResponse, error)
	ListUsers(ctx context.Context, request ListUsersRequest) (*ListUsersResponse, error)
	SearchUsers(ctx context.Context, request SearchUsersRequest) (*SearchUsersResponse, error)
}

// GetUser

type GetUserRequest struct {
	Email string
}

type GetUserResponse struct {
	User User
}

// ListUsers

type ListUsersRequest struct {
	Offset int
	Limit  int
}

type ListUsersResponse struct {
	Users []User
	Total int64
}

// SearchUsers

type SearchUsersRequest struct {
	Query  string
	Offset int
	Limit  int
}

type SearchUsersResponse struct {
	Users []User
	Total int64
}
