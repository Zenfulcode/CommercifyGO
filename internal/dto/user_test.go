package dto

import (
	"testing"
	"time"
)

func TestUserDTO(t *testing.T) {
	now := time.Now()
	user := UserDTO{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got %s", user.Email)
	}
	if user.FirstName != "John" {
		t.Errorf("Expected FirstName 'John', got %s", user.FirstName)
	}
	if user.LastName != "Doe" {
		t.Errorf("Expected LastName 'Doe', got %s", user.LastName)
	}
	if user.Role != "admin" {
		t.Errorf("Expected Role 'admin', got %s", user.Role)
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt %v, got %v", now, user.CreatedAt)
	}
	if !user.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt %v, got %v", now, user.UpdatedAt)
	}
}

func TestCreateUserRequest(t *testing.T) {
	request := CreateUserRequest{
		Email:     "newuser@example.com",
		Password:  "securepassword",
		FirstName: "Jane",
		LastName:  "Smith",
	}

	if request.Email != "newuser@example.com" {
		t.Errorf("Expected Email 'newuser@example.com', got %s", request.Email)
	}
	if request.Password != "securepassword" {
		t.Errorf("Expected Password 'securepassword', got %s", request.Password)
	}
	if request.FirstName != "Jane" {
		t.Errorf("Expected FirstName 'Jane', got %s", request.FirstName)
	}
	if request.LastName != "Smith" {
		t.Errorf("Expected LastName 'Smith', got %s", request.LastName)
	}
}

func TestUpdateUserRequest(t *testing.T) {
	request := UpdateUserRequest{
		FirstName: "UpdatedFirst",
		LastName:  "UpdatedLast",
	}

	if request.FirstName != "UpdatedFirst" {
		t.Errorf("Expected FirstName 'UpdatedFirst', got %s", request.FirstName)
	}
	if request.LastName != "UpdatedLast" {
		t.Errorf("Expected LastName 'UpdatedLast', got %s", request.LastName)
	}
}

func TestUpdateUserRequestEmpty(t *testing.T) {
	request := UpdateUserRequest{}

	if request.FirstName != "" {
		t.Errorf("Expected FirstName empty, got %s", request.FirstName)
	}
	if request.LastName != "" {
		t.Errorf("Expected LastName empty, got %s", request.LastName)
	}
}

func TestUserLoginRequest(t *testing.T) {
	request := UserLoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	}

	if request.Email != "user@example.com" {
		t.Errorf("Expected Email 'user@example.com', got %s", request.Email)
	}
	if request.Password != "password123" {
		t.Errorf("Expected Password 'password123', got %s", request.Password)
	}
}

func TestUserLoginResponse(t *testing.T) {
	now := time.Now()
	user := UserDTO{
		ID:        1,
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := UserLoginResponse{
		User:         user,
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_456",
		ExpiresIn:    3600,
	}

	if response.User.ID != 1 {
		t.Errorf("Expected User.ID 1, got %d", response.User.ID)
	}
	if response.User.Email != "user@example.com" {
		t.Errorf("Expected User.Email 'user@example.com', got %s", response.User.Email)
	}
	if response.AccessToken != "access_token_123" {
		t.Errorf("Expected AccessToken 'access_token_123', got %s", response.AccessToken)
	}
	if response.RefreshToken != "refresh_token_456" {
		t.Errorf("Expected RefreshToken 'refresh_token_456', got %s", response.RefreshToken)
	}
	if response.ExpiresIn != 3600 {
		t.Errorf("Expected ExpiresIn 3600, got %d", response.ExpiresIn)
	}
}

func TestUserListResponse(t *testing.T) {
	users := []UserDTO{
		{
			ID:        1,
			Email:     "user1@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Role:      "user",
		},
		{
			ID:        2,
			Email:     "user2@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
			Role:      "admin",
		},
	}

	pagination := PaginationDTO{
		Page:     1,
		PageSize: 10,
		Total:    2,
	}

	response := UserListResponse{
		ListResponseDTO: ListResponseDTO[UserDTO]{
			Success:    true,
			Message:    "Users retrieved successfully",
			Data:       users,
			Pagination: pagination,
		},
	}

	if !response.Success {
		t.Errorf("Expected Success true, got %t", response.Success)
	}
	if len(response.Data) != 2 {
		t.Errorf("Expected Data length 2, got %d", len(response.Data))
	}
	if response.Data[0].Email != "user1@example.com" {
		t.Errorf("Expected Data[0].Email 'user1@example.com', got %s", response.Data[0].Email)
	}
	if response.Data[1].Role != "admin" {
		t.Errorf("Expected Data[1].Role 'admin', got %s", response.Data[1].Role)
	}
	if response.Pagination.Total != 2 {
		t.Errorf("Expected Pagination.Total 2, got %d", response.Pagination.Total)
	}
}

func TestChangePasswordRequest(t *testing.T) {
	request := ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword123",
	}

	if request.CurrentPassword != "oldpassword" {
		t.Errorf("Expected CurrentPassword 'oldpassword', got %s", request.CurrentPassword)
	}
	if request.NewPassword != "newpassword123" {
		t.Errorf("Expected NewPassword 'newpassword123', got %s", request.NewPassword)
	}
}

func TestChangePasswordRequestEmpty(t *testing.T) {
	request := ChangePasswordRequest{}

	if request.CurrentPassword != "" {
		t.Errorf("Expected CurrentPassword empty, got %s", request.CurrentPassword)
	}
	if request.NewPassword != "" {
		t.Errorf("Expected NewPassword empty, got %s", request.NewPassword)
	}
}
