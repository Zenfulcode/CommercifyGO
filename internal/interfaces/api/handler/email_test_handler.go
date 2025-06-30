package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"gorm.io/gorm"
)

// EmailTestHandler handles email testing endpoints
type EmailTestHandler struct {
	emailSvc service.EmailService
	logger   logger.Logger
	config   config.EmailConfig
}

// NewEmailTestHandler creates a new EmailTestHandler
func NewEmailTestHandler(emailSvc service.EmailService, logger logger.Logger, emailConfig config.EmailConfig) *EmailTestHandler {
	return &EmailTestHandler{
		emailSvc: emailSvc,
		logger:   logger,
		config:   emailConfig,
	}
}

// TestEmail sends test order confirmation and notification emails
func (h *EmailTestHandler) TestEmail(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Test email endpoint called")

	// Create a mock user (but we'll send emails to admin address)
	mockUser := &entity.User{
		Email:     "customer@example.com", // This is just for the mock data
		FirstName: "John",
		LastName:  "Doe",
	}

	// Create a mock order
	mockOrder := &entity.Order{
		Model: gorm.Model{
			ID:        12345, // Mock order ID
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		OrderNumber:       "ORD-12345",
		UserID:            mockUser.ID,
		Status:            entity.OrderStatusCompleted,
		PaymentStatus:     entity.PaymentStatusCaptured,
		TotalAmount:       9950, // $99.50 in cents (subtotal before shipping/discounts)
		ShippingCost:      850,  // $8.50 shipping cost
		DiscountAmount:    1500, // $15.00 discount
		FinalAmount:       8300, // $83.00 final amount (99.50 + 8.50 - 15.00)
		Currency:          "USD",
		CheckoutSessionID: "test-checkout-session-12345", // Add checkout session ID for testing
		CustomerDetails: &entity.CustomerDetails{
			Email:    mockUser.Email,
			Phone:    "+1234567890",
			FullName: mockUser.FirstName + " " + mockUser.LastName,
		},
		IsGuestOrder:    false,
		PaymentProvider: "stripe",
		PaymentMethod:   "card",
		Items: []entity.OrderItem{
			{
				ProductID:   1,
				Quantity:    2,
				Price:       2500, // $25.00 in cents
				Subtotal:    5000, // $50.00 in cents
				ProductName: "Test Product 1",
				SKU:         "TEST-001",
			},
			{
				ProductID:   2,
				Quantity:    1,
				Price:       4950, // $49.50 in cents
				Subtotal:    4950, // $49.50 in cents
				ProductName: "Test Product 2",
				SKU:         "TEST-002",
			},
		},
	}

	var errors []string

	// Override email addresses to send both emails to admin for testing
	adminUser := &entity.User{
		Email:     h.config.AdminEmail, // Send to admin email
		FirstName: mockUser.FirstName,
		LastName:  mockUser.LastName,
	}

	// Also update the order's customer details to use admin email for testing
	testOrder := *mockOrder
	testOrder.CustomerDetails = &entity.CustomerDetails{
		Email:    h.config.AdminEmail, // Send to admin email
		Phone:    mockOrder.CustomerDetails.Phone,
		FullName: mockOrder.CustomerDetails.FullName,
	}

	// Send order confirmation email to admin (instead of customer)
	h.logger.Info("Sending test order confirmation email to admin: %s", h.config.AdminEmail)
	if err := h.emailSvc.SendOrderConfirmation(&testOrder, adminUser); err != nil {
		h.logger.Error("Failed to send order confirmation email: %v", err)
		errors = append(errors, "Order confirmation: "+err.Error())
	} else {
		h.logger.Info("Order confirmation email sent successfully")
	}

	// Send order notification email to admin
	h.logger.Info("Sending test order notification email to admin: %s", h.config.AdminEmail)
	if err := h.emailSvc.SendOrderNotification(&testOrder, adminUser); err != nil {
		h.logger.Error("Failed to send order notification email: %v", err)
		errors = append(errors, "Order notification: "+err.Error())
	} else {
		h.logger.Info("Order notification email sent successfully")
	}

	// Set addresses using JSON helper methods
	shippingAddr := entity.Address{
		Street1:    "123 Test Street",
		Street2:    "street 2",
		City:       "Test City",
		State:      "Test State",
		PostalCode: "12345",
		Country:    "US",
	}
	billingAddr := entity.Address{
		Street1:    "123 Test Street",
		Street2:    "",
		City:       "Test City",
		State:      "Test State",
		PostalCode: "12345",
		Country:    "US",
	}

	mockOrder.SetShippingAddressJSON(&shippingAddr)
	mockOrder.SetBillingAddressJSON(&billingAddr)

	// Set applied discount using JSON helper method
	appliedDiscount := &entity.AppliedDiscount{
		DiscountID:     1,
		DiscountCode:   "SUMMER25",
		DiscountAmount: 1500, // $15.00 discount
	}
	mockOrder.SetAppliedDiscountJSON(appliedDiscount)

	w.Header().Set("Content-Type", "application/json")

	if len(errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"errors":  errors,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Both order confirmation and notification emails sent successfully",
		"details": map[string]string{
			"customer_email": mockUser.Email,
			"order_id":       "12345",
		},
	})
}
