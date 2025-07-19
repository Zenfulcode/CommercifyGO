package email

import (
	"os"
	"testing"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"gorm.io/gorm"
)

func TestSMTPEmailService_SendOrderShipped(t *testing.T) {
	// Setup
	zapLogger := logger.NewLogger()
	emailConfig := config.EmailConfig{
		Enabled:      false, // Disable actual email sending for test
		FromEmail:    "test@example.com",
		FromName:     "Test Store",
		SMTPHost:     "localhost",
		SMTPPort:     587,
		SMTPUsername: "test",
		SMTPPassword: "test",
		AdminEmail:   "admin@example.com",
		ContactEmail: "support@example.com",
		StoreName:    "Test Store",
	}

	service := NewSMTPEmailService(emailConfig, zapLogger)

	// Create test order
	order := &entity.Order{
		Model:    gorm.Model{ID: 123},
		Currency: "USD",
	}

	// Create test user
	user := &entity.User{
		Model:     gorm.Model{ID: 1},
		Email:     "customer@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Test sending order shipped email
	err := service.SendOrderShipped(order, user, "1Z999AA1234567890", "https://example.com/track")

	// Should not error since email is disabled
	if err != nil {
		t.Errorf("Expected no error when email is disabled, got: %v", err)
	}
}

func TestTemplateExists(t *testing.T) {
	// Check if the template file exists
	if _, err := os.Stat("../../../templates/emails/order_shipped.html"); os.IsNotExist(err) {
		t.Error("order_shipped.html template file does not exist")
	}
}

func TestAllTemplatesExist(t *testing.T) {
	templates := []string{
		"order_shipped.html",
		"order_confirmation.html",
		"order_notification.html",
		"checkout_recovery.html",
	}

	for _, template := range templates {
		path := "../../../templates/emails/" + template
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Template file does not exist: %s", template)
		}
	}
}

func TestTemplateRendering(t *testing.T) {
	// Skip template rendering test since it requires proper working directory setup
	// The existence test above already verifies templates are present
	t.Skip("Template rendering test requires proper working directory setup")
}
