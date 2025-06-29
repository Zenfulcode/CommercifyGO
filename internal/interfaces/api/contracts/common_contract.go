package contracts

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
