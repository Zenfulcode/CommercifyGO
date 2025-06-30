package contracts

import (
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/dto"
)

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
	ExpiresIn    int         `json:"expires_in"`
}

func (r *UserLoginRequest) ToUseCaseInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

// ChangePasswordRequest represents the data needed to change a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (r *CreateUserRequest) ToUseCaseInput() usecase.RegisterInput {
	return usecase.RegisterInput{
		Email:     r.Email,
		Password:  r.Password,
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}
}

func CreateUserLoginResponse(user *dto.UserDTO, accessToken, refreshToken string, expiresIn int) UserLoginResponse {
	return UserLoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

func CreateUserListResponse(users []*entity.User, totalCount, page, pageSize int) ListResponseDTO[dto.UserDTO] {
	var userDTOs []dto.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, *user.ToUserDTO())
	}

	if len(userDTOs) == 0 {
		return ListResponseDTO[dto.UserDTO]{
			Data: []dto.UserDTO{},
			Pagination: PaginationDTO{
				Total:    totalCount,
				Page:     page,
				PageSize: pageSize,
			},
			Success: true,
			Message: "No users found",
		}
	}

	return ListResponseDTO[dto.UserDTO]{
		Success: true,
		Data:    userDTOs,
		Pagination: PaginationDTO{
			Total:    totalCount,
			Page:     page,
			PageSize: pageSize,
		},
		Message: "Users retrieved successfully",
	}
}
