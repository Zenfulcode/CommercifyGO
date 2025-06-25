package dto

// PaginationDTO represents pagination parameters
type PaginationDTO struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// ResponseDTO is a generic response wrapper
type ResponseDTO[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListResponseDTO is a generic list response wrapper
type ListResponseDTO[T any] struct {
	Success    bool          `json:"success"`
	Message    string        `json:"message,omitempty"`
	Data       []T           `json:"data"`
	Pagination PaginationDTO `json:"pagination"`
	Error      string        `json:"error,omitempty"`
}

// AddressDTO represents a shipping or billing address
type AddressDTO struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
}

// CustomerDetailsDTO represents customer information for a checkout
type CustomerDetailsDTO struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

func ErrorResponse(message string) ResponseDTO[any] {
	return ResponseDTO[any]{
		Success: false,
		Error:   message,
	}
}

func SuccessResponseWithMessage[T any](data T, message string) ResponseDTO[T] {
	return ResponseDTO[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func SuccessResponseMessage(message string) ResponseDTO[any] {
	return ResponseDTO[any]{
		Success: true,
		Message: message,
		Data:    nil,
	}
}

func SuccessResponse[T any](data T) ResponseDTO[T] {
	response := ResponseDTO[T]{
		Success: true,
		Data:    data,
	}

	return response
}
