package contracts

import "github.com/zenfulcode/commercify/internal/dto"

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateUserRequest represents the data needed to update an existing user
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// UserLoginRequest represents the data needed for user login
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginResponse represents the response after successful login
type UserLoginResponse struct {
	User         dto.UserDTO `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	ListResponseDTO[dto.UserDTO]
}

// ChangePasswordRequest represents the data needed to change a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
