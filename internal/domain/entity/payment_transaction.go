package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TransactionType represents the type of payment transaction
type TransactionType string

const (
	TransactionTypeAuthorize TransactionType = "authorize"
	TransactionTypeCapture   TransactionType = "capture"
	TransactionTypeRefund    TransactionType = "refund"
	TransactionTypeCancel    TransactionType = "cancel"
)

// TransactionStatus represents the status of a payment transaction
type TransactionStatus string

const (
	TransactionStatusSuccessful TransactionStatus = "successful"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusPending    TransactionStatus = "pending"
)

// PaymentTransaction represents a payment transaction record
// Each order can have multiple transactions per type (for scenarios like partial captures, retries, webhooks, etc.)
// Each transaction represents a specific event in the payment lifecycle
type PaymentTransaction struct {
	gorm.Model
	OrderID       uint              `gorm:"index;not null"`              // Foreign key to order (indexed for performance)
	Order         Order             `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TransactionID string            `gorm:"uniqueIndex;not null;size:100"`                     // Human-readable transaction number (e.g., "TXN-AUTH-2025-001")
	ExternalID    string            `gorm:"index;size:255"`                                    // External transaction ID from payment provider (can be empty for some providers)
	Type          TransactionType   `gorm:"not null;size:50;index:idx_order_type"`            // Type of transaction (authorize, capture, refund, cancel)
	Status        TransactionStatus `gorm:"not null;size:50"`                                  // Status of the transaction (pending -> successful/failed)
	Amount        int64             `gorm:"not null"`                                          // Amount of the transaction
	Currency      string            `gorm:"not null;size:3"`                                   // Currency of the transaction
	Provider      string            `gorm:"not null;size:100"`                                 // Payment provider (stripe, paypal, etc.)
	RawResponse   string            `gorm:"type:text"`                                         // Raw response from payment provider (JSON)
	Metadata      datatypes.JSONMap `gorm:"type:text"`                                         // Additional metadata stored as JSON
}

// NewPaymentTransaction creates a new payment transaction
func NewPaymentTransaction(
	orderID uint,
	externalID string,
	transactionType TransactionType,
	status TransactionStatus,
	amount int64,
	currency string,
	provider string,
) (*PaymentTransaction, error) {
	if orderID == 0 {
		return nil, errors.New("orderID cannot be zero")
	}
	if string(transactionType) == "" {
		return nil, errors.New("transactionType cannot be empty")
	}
	if string(status) == "" {
		return nil, errors.New("status cannot be empty")
	}
	if provider == "" {
		return nil, errors.New("provider cannot be empty")
	}
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	return &PaymentTransaction{
		OrderID:    orderID,
		ExternalID: externalID, // Can be empty for some providers
		Type:       transactionType,
		Status:     status,
		Amount:     amount,
		Currency:   currency,
		Provider:   provider,
		Metadata:   make(datatypes.JSONMap),
		// TransactionID will be set when the transaction is saved to get the sequence number
	}, nil
}

// AddMetadata adds metadata to the transaction
func (pt *PaymentTransaction) AddMetadata(key, value string) {
	if pt.Metadata == nil {
		pt.Metadata = make(datatypes.JSONMap)
	}
	pt.Metadata[key] = value
}

// SetRawResponse sets the raw response from the payment provider
func (pt *PaymentTransaction) SetRawResponse(response string) {
	pt.RawResponse = response

}

// UpdateStatus updates the status of the transaction
func (pt *PaymentTransaction) UpdateStatus(status TransactionStatus) {
	pt.Status = status

}

// SetTransactionID sets the friendly number for the transaction
func (pt *PaymentTransaction) SetTransactionID(sequence int) {
	pt.TransactionID = generateTransactionID(pt.Type, sequence)
}

// GetDisplayName returns a user-friendly display name for the transaction
func (pt *PaymentTransaction) GetDisplayName() string {
	if pt.TransactionID != "" {
		return pt.TransactionID
	}
	// Fallback to external ID if transaction ID is not set
	return pt.ExternalID
}

// GetTypeDisplayName returns a user-friendly name for the transaction type
func (pt *PaymentTransaction) GetTypeDisplayName() string {
	switch pt.Type {
	case TransactionTypeAuthorize:
		return "Authorization"
	case TransactionTypeCapture:
		return "Capture"
	case TransactionTypeRefund:
		return "Refund"
	case TransactionTypeCancel:
		return "Cancellation"
	default:
		return string(pt.Type)
	}
}

// generateTransactionID generates a human-readable transaction ID
// This becomes the primary TransactionID field in the database
// Format: TXN-{TYPE}-{YEAR}-{SEQUENCE}
// Examples: TXN-AUTH-2025-001, TXN-CAPT-2025-002, TXN-REFUND-2025-001
func generateTransactionID(transactionType TransactionType, sequence int) string {
	year := time.Now().Year()
	typeCode := strings.ToUpper(string(transactionType))

	// Create shorter type codes for better readability
	switch transactionType {
	case TransactionTypeAuthorize:
		typeCode = "AUTH"
	case TransactionTypeCapture:
		typeCode = "CAPT"
	case TransactionTypeRefund:
		typeCode = "REFUND"
	case TransactionTypeCancel:
		typeCode = "CANCEL"
	}

	return fmt.Sprintf("TXN-%s-%d-%03d", typeCode, year, sequence)
}

// SetExternalID sets the external payment provider ID
func (pt *PaymentTransaction) SetExternalID(externalID string) {
	pt.ExternalID = externalID
}
