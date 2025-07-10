package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// SMTPEmailService implements the email service interface using SMTP
type SMTPEmailService struct {
	config config.EmailConfig
	logger logger.Logger
}

// NewSMTPEmailService creates a new SMTPEmailService
func NewSMTPEmailService(config config.EmailConfig, logger logger.Logger) *SMTPEmailService {
	return &SMTPEmailService{
		config: config,
		logger: logger,
	}
}

// SendEmail sends an email with the given data
func (s *SMTPEmailService) SendEmail(data service.EmailData) error {
	s.logger.Info("Attempting to send email to: %s, Subject: %s, Enabled: %t", data.To, data.Subject, s.config.Enabled)

	// If email service is disabled, log and return
	if !s.config.Enabled {
		s.logger.Info("Email service is disabled. Would have sent email to: %s, Subject: %s", data.To, data.Subject)
		return nil
	}

	// Prepare email content
	var body string
	var err error

	if data.Template != "" && data.Data != nil {
		// Use template if provided
		body, err = s.renderTemplate(data.Template, data.Data)
		if err != nil {
			s.logger.Error("Failed to render email template %s: %v", data.Template, err)
			return err
		}
		s.logger.Info("Email template rendered successfully")
	} else {
		// Use provided body
		body = data.Body
	}

	// Set up authentication information
	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUsername,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	// Prepare email headers
	contentType := "text/plain"
	if data.IsHTML {
		contentType = "text/html"
	}

	// Format email message
	msg := fmt.Appendf(nil, "From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", s.config.FromName, s.config.FromEmail, data.To, data.Subject, contentType, body)

	// Send email
	s.logger.Info("Attempting to send email via SMTP to %s:%d", s.config.SMTPHost, s.config.SMTPPort)
	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.FromEmail,
		[]string{data.To},
		msg,
	)

	if err != nil {
		s.logger.Error("Failed to send email to %s: %v", data.To, err)
		return err
	}

	s.logger.Info("Email sent successfully to %s", data.To)
	return nil
}

// SendOrderConfirmation sends an order confirmation email to the customer
func (s *SMTPEmailService) SendOrderConfirmation(order *entity.Order, user *entity.User) error {
	s.logger.Info("Sending order confirmation email for Order ID: %d to User: %s", order.ID, user.Email)

	// Prepare data for the template
	shippingAddr := order.GetShippingAddress()
	billingAddr := order.GetBillingAddress()
	appliedDiscount := order.GetAppliedDiscount()

	// Debug logging
	s.logger.Info("Email template data - Order ID: %d", order.ID)
	s.logger.Info("Shipping Address: %+v", shippingAddr)
	s.logger.Info("Billing Address: %+v", billingAddr)
	s.logger.Info("Applied Discount: %+v", appliedDiscount)

	data := map[string]any{
		"Order":           order,
		"User":            user,
		"StoreName":       s.config.FromName,
		"ContactEmail":    s.config.FromEmail,
		"AppliedDiscount": appliedDiscount,
		"ShippingAddr":    shippingAddr,
		"BillingAddr":     billingAddr,
		"Currency":        order.Currency,
	}

	// Send email
	return s.SendEmail(service.EmailData{
		To:       user.Email,
		Subject:  fmt.Sprintf("Order Confirmation #%d", order.ID),
		IsHTML:   true,
		Template: "order_confirmation.html",
		Data:     data,
	})
}

// SendOrderNotification sends an order notification email to the admin
func (s *SMTPEmailService) SendOrderNotification(order *entity.Order, user *entity.User) error {
	s.logger.Info("Sending order notification email for Order ID: %d to Admin: %s", order.ID, s.config.AdminEmail)

	// Prepare data for the template
	shippingAddr := order.GetShippingAddress()
	billingAddr := order.GetBillingAddress()
	appliedDiscount := order.GetAppliedDiscount()

	// Debug logging
	s.logger.Info("Email template data - Order ID: %d", order.ID)

	data := map[string]any{
		"Order":           order,
		"User":            user,
		"StoreName":       s.config.FromName,
		"AppliedDiscount": appliedDiscount,
		"ShippingAddr":    shippingAddr,
		"BillingAddr":     billingAddr,
		"Currency":        order.Currency,
	}

	// Send email
	return s.SendEmail(service.EmailData{
		To:       s.config.AdminEmail,
		Subject:  fmt.Sprintf("New Order #%d Received", order.ID),
		IsHTML:   true,
		Template: "order_notification.html",
		Data:     data,
	})
}

// SendOrderShipped sends an order shipped notification email to the customer
func (s *SMTPEmailService) SendOrderShipped(order *entity.Order, user *entity.User, trackingNumber, trackingURL string) error {
	s.logger.Info("Sending order shipped email for Order ID: %d to User: %s", order.ID, user.Email)

	// Prepare data for the template
	shippingAddr := order.GetShippingAddress()
	billingAddr := order.GetBillingAddress()
	appliedDiscount := order.GetAppliedDiscount()

	// Debug logging
	s.logger.Info("Email template data - Order ID: %d", order.ID)
	s.logger.Info("Tracking Number: %s", trackingNumber)
	s.logger.Info("Tracking URL: %s", trackingURL)
	s.logger.Info("Shipping Address: %+v", shippingAddr)

	data := map[string]any{
		"Order":           order,
		"User":            user,
		"StoreName":       s.config.FromName,
		"ContactEmail":    s.config.FromEmail,
		"AppliedDiscount": appliedDiscount,
		"ShippingAddr":    shippingAddr,
		"BillingAddr":     billingAddr,
		"Currency":        order.Currency,
		"TrackingNumber":  trackingNumber,
		"TrackingURL":     trackingURL,
	}

	// Send email
	return s.SendEmail(service.EmailData{
		To:       user.Email,
		Subject:  fmt.Sprintf("Your Order #%d Has Been Shipped! 📦", order.ID),
		IsHTML:   true,
		Template: "order_shipped.html",
		Data:     data,
	})
}

// renderTemplate renders an HTML template with the given data
func (s *SMTPEmailService) renderTemplate(templateName string, data map[string]any) (string, error) {
	// Get template path
	templatePath := filepath.Join("templates", "emails", templateName)

	// Create template with helper functions
	tmpl := template.New(templateName).Funcs(template.FuncMap{
		"centsToDollars": func(cents int64) float64 {
			return money.FromCents(cents)
		},
		"formatPrice": func(cents int64) string {
			return fmt.Sprintf("%.2f", money.FromCents(cents))
		},
		"formatPriceWithCurrency": func(cents int64, currency string) string {
			return s.formatCurrency(cents, currency)
		},
	})

	// Parse template
	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatCurrency formats a cents amount with the currency code at the end
func (s *SMTPEmailService) formatCurrency(amount int64, currency string) string {
	// Format amount as decimal
	decimal := money.FromCents(amount)

	// Format with currency code at the end for all currencies
	if currency == "JPY" {
		// JPY typically doesn't use decimals
		return fmt.Sprintf("%.0f %s", decimal*100, currency) // Convert back to whole yen
	}
	return fmt.Sprintf("%.2f %s", decimal, currency)
}
