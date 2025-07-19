package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/auth"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
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
	var request contracts.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid request body",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	input := request.ToUseCaseInput()

	user, err := h.userUseCase.Register(input)
	if err != nil {
		h.logger.Error("Failed to register user: %v", err)
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, expirationTime, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to generate token",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create login response
	loginResponse := contracts.CreateUserLoginResponse(user.ToUserDTO(), token, "", expirationTime)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loginResponse)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request contracts.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := contracts.ResponseDTO[any]{
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
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Invalid email or password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, expiresIn, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to generate token",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert domain user to DTO
	response := contracts.CreateUserLoginResponse(
		user.ToUserDTO(),
		token,
		"", // Refresh token not implemented
		expiresIn,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile handles getting the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in CreateProduct")
		response := contracts.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userUseCase.GetUserByID(userID)
	if err != nil {
		h.logger.Error("Failed to get user profile: %v", err)
		response := contracts.ErrorResponse("Failed to get user profile")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.SuccessResponse(user.ToUserDTO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile handles updating the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Unauthorized",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var request contracts.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := contracts.ResponseDTO[any]{
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
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to update user profile",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.SuccessResponse(user.ToUserDTO())

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
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to list users",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// TODO: Get total count from repository
	total := len(users)

	response := contracts.CreateUserListResponse(users, total, page, pageSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ChangePassword handles changing the user's password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Unauthorized",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	var request contracts.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := contracts.ResponseDTO[any]{
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
		response := contracts.ResponseDTO[any]{
			Success: false,
			Error:   "Failed to change password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.ResponseDTO[any]{
		Success: true,
		Message: "Password changed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
