package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/auth"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/middleware"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	jwtService  *auth.JWTService
	logger      logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase, jwtService *auth.JWTService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid request body",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	input := usecase.RegisterInput{
		Email:     request.Email,
		Password:  request.Password,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	user, err := h.userUseCase.Register(input)
	if err != nil {
		h.logger.Error("Failed to register user: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to generate token",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert domain user to DTO
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Create login response
	loginResponse := dto.UserLoginResponse{
		User:         userDTO,
		AccessToken:  token,
		RefreshToken: "",   // TODO: Implement refresh token
		ExpiresIn:    3600, // TODO: Make this configurable
	}

	response := dto.ResponseDTO[dto.UserLoginResponse]{
		Success: true,
		Data:    loginResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid request body",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	input := usecase.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	}

	user, err := h.userUseCase.Login(input)
	if err != nil {
		h.logger.Error("Login failed: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid email or password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to generate token",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert domain user to DTO
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Create login response
	loginResponse := dto.UserLoginResponse{
		User:         userDTO,
		AccessToken:  token,
		RefreshToken: "",   // TODO: Implement refresh token
		ExpiresIn:    3600, // TODO: Make this configurable
	}

	response := dto.ResponseDTO[dto.UserLoginResponse]{
		Success: true,
		Data:    loginResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile handles getting the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in CreateProduct")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userUseCase.GetUserByID(userID)
	if err != nil {
		h.logger.Error("Failed to get user profile: %v", err)
		response := dto.ErrorResponse("Failed to get user profile")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert domain user to DTO
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := dto.SuccessResponse(userDTO)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile handles updating the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Unauthorized",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var request dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid request body",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	input := usecase.UpdateUserInput{
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	user, err := h.userUseCase.UpdateUser(userID, input)
	if err != nil {
		h.logger.Error("Failed to update user profile: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to update user profile",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert domain user to DTO
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := dto.ResponseDTO[dto.UserDTO]{
		Success: true,
		Data:    userDTO,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListUsers handles listing all users (admin only)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	offset := (page - 1) * pageSize
	users, err := h.userUseCase.ListUsers(offset, pageSize)
	if err != nil {
		h.logger.Error("Failed to list users: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to list users",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert domain users to DTOs
	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	// TODO: Get total count from repository
	total := len(users)

	response := dto.UserListResponse{
		ListResponseDTO: dto.ListResponseDTO[dto.UserDTO]{
			Success: true,
			Data:    userDTOs,
			Pagination: dto.PaginationDTO{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ChangePassword handles changing the user's password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Unauthorized",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var request dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid request body",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	input := usecase.ChangePasswordInput{
		CurrentPassword: request.CurrentPassword,
		NewPassword:     request.NewPassword,
	}

	err := h.userUseCase.ChangePassword(userID, input)
	if err != nil {
		h.logger.Error("Failed to change password: %v", err)
		response := dto.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to change password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.ResponseDTO[any]{
		Success: true,
		Message: "Password changed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
