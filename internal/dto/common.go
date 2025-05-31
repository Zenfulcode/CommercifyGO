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
	Data       []T           `json:"data,omitempty"`
	Pagination PaginationDTO `json:"pagination,omitempty"`
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
