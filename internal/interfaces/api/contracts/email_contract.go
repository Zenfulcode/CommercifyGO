package contracts

import "github.com/zenfulcode/commercify/internal/infrastructure/validation"

// EmailTestRequest represents the request body for testing emails
type EmailTestRequest struct {
	Email string `json:"email"`
}

// Validate validates the email test request
func (r *EmailTestRequest) Validate() error {
	return validation.ValidateEmail(r.Email)
}
