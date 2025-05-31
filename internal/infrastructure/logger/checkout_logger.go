package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// CheckoutLogger provides logging for checkout operations
type CheckoutLogger struct {
	logFile *os.File
}

// CheckoutEventType defines the type of checkout event
type CheckoutEventType string

const (
	// CheckoutCreated event occurs when a new checkout is created
	CheckoutCreated CheckoutEventType = "CREATED"
	// CheckoutUpdated event occurs when a checkout is updated
	CheckoutUpdated CheckoutEventType = "UPDATED"
	// CheckoutItemAdded event occurs when an item is added to a checkout
	CheckoutItemAdded CheckoutEventType = "ITEM_ADDED"
	// CheckoutItemRemoved event occurs when an item is removed from a checkout
	CheckoutItemRemoved CheckoutEventType = "ITEM_REMOVED"
	// CheckoutItemUpdated event occurs when an item in a checkout is updated
	CheckoutItemUpdated CheckoutEventType = "ITEM_UPDATED"
	// CheckoutAddressSet event occurs when an address is set for a checkout
	CheckoutAddressSet CheckoutEventType = "ADDRESS_SET"
	// CheckoutShippingMethodSet event occurs when a shipping method is set for a checkout
	CheckoutShippingMethodSet CheckoutEventType = "SHIPPING_METHOD_SET"
	// CheckoutDiscountApplied event occurs when a discount is applied to a checkout
	CheckoutDiscountApplied CheckoutEventType = "DISCOUNT_APPLIED"
	// CheckoutDiscountRemoved event occurs when a discount is removed from a checkout
	CheckoutDiscountRemoved CheckoutEventType = "DISCOUNT_REMOVED"
	// CheckoutCompleted event occurs when a checkout is completed
	CheckoutCompleted CheckoutEventType = "COMPLETED"
	// CheckoutAbandoned event occurs when a checkout is marked as abandoned
	CheckoutAbandoned CheckoutEventType = "ABANDONED"
	// CheckoutExpired event occurs when a checkout is marked as expired
	CheckoutExpired CheckoutEventType = "EXPIRED"
	// CheckoutConverted event occurs when a guest checkout is converted to a user checkout
	CheckoutConverted CheckoutEventType = "CONVERTED"
	// CheckoutRecoveryEmailSent event occurs when a recovery email is sent for an abandoned checkout
	CheckoutRecoveryEmailSent CheckoutEventType = "RECOVERY_EMAIL_SENT"
)

// NewCheckoutLogger creates a new instance of CheckoutLogger
func NewCheckoutLogger(logPath string) (*CheckoutLogger, error) {
	// Create log file with append mode
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open checkout log file: %w", err)
	}

	return &CheckoutLogger{
		logFile: logFile,
	}, nil
}

// Close closes the log file
func (l *CheckoutLogger) Close() error {
	return l.logFile.Close()
}

// Log logs a checkout event
func (l *CheckoutLogger) Log(eventType CheckoutEventType, checkout *entity.Checkout, additionalInfo string) error {
	timestamp := time.Now().Format(time.RFC3339)
	userIdentifier := fmt.Sprintf("user:%d", checkout.UserID)
	if checkout.UserID == 0 {
		userIdentifier = fmt.Sprintf("session:%s", checkout.SessionID)
	}

	logEntry := fmt.Sprintf(
		"%s [CHECKOUT:%d] [%s] [%s] [STATUS:%s] [TOTAL:%d] [CURRENCY:%s] %s\n",
		timestamp,
		checkout.ID,
		eventType,
		userIdentifier,
		checkout.Status,
		checkout.FinalAmount,
		checkout.Currency,
		additionalInfo,
	)

	_, err := l.logFile.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("failed to write to checkout log: %w", err)
	}

	return nil
}

// LogItemEvent logs a checkout item event
func (l *CheckoutLogger) LogItemEvent(eventType CheckoutEventType, checkout *entity.Checkout, productID uint, variantID uint, quantity int) error {
	additionalInfo := fmt.Sprintf("Product:%d Variant:%d Quantity:%d", productID, variantID, quantity)
	return l.Log(eventType, checkout, additionalInfo)
}

// LogDiscountEvent logs a discount event
func (l *CheckoutLogger) LogDiscountEvent(eventType CheckoutEventType, checkout *entity.Checkout, discountCode string, discountAmount int64) error {
	additionalInfo := fmt.Sprintf("DiscountCode:%s Amount:%d", discountCode, discountAmount)
	return l.Log(eventType, checkout, additionalInfo)
}

// LogShippingEvent logs a shipping method event
func (l *CheckoutLogger) LogShippingEvent(eventType CheckoutEventType, checkout *entity.Checkout, shippingMethodID uint, shippingCost int64) error {
	additionalInfo := fmt.Sprintf("ShippingMethodID:%d Cost:%d", shippingMethodID, shippingCost)
	return l.Log(eventType, checkout, additionalInfo)
}

// LogAddressEvent logs an address event
func (l *CheckoutLogger) LogAddressEvent(eventType CheckoutEventType, checkout *entity.Checkout, addressType string) error {
	additionalInfo := fmt.Sprintf("AddressType:%s", addressType)
	return l.Log(eventType, checkout, additionalInfo)
}

// LogCompletionEvent logs a checkout completion event
func (l *CheckoutLogger) LogCompletionEvent(eventType CheckoutEventType, checkout *entity.Checkout, orderID uint) error {
	additionalInfo := fmt.Sprintf("OrderID:%d", orderID)
	return l.Log(eventType, checkout, additionalInfo)
}

// LogConversionEvent logs a checkout conversion event
func (l *CheckoutLogger) LogConversionEvent(checkout *entity.Checkout, sessionID string, userID uint) error {
	additionalInfo := fmt.Sprintf("FromSession:%s ToUser:%d", sessionID, userID)
	return l.Log(CheckoutConverted, checkout, additionalInfo)
}

// LogRecoveryEmailEvent logs a recovery email event
func (l *CheckoutLogger) LogRecoveryEmailEvent(checkout *entity.Checkout, email string) error {
	additionalInfo := fmt.Sprintf("Email:%s", email)
	return l.Log(CheckoutRecoveryEmailSent, checkout, additionalInfo)
}
