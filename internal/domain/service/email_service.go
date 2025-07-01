package service

import "github.com/zenfulcode/commercify/internal/domain/entity"

// EmailData represents the data needed to send an email
type EmailData struct {
	To       string
	Subject  string
	Body     string
	IsHTML   bool
	Template string
	Data     map[string]any
}

// EmailService defines the interface for email operations
type EmailService interface {
	// SendEmail sends an email with the given data
	SendEmail(data EmailData) error

	// SendOrderConfirmation sends an order confirmation email to the customer
	SendOrderConfirmation(order *entity.Order, user *entity.User) error

	// SendOrderNotification sends an order notification email to the admin
	SendOrderNotification(order *entity.Order, user *entity.User) error
}
