package entity

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/common"
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
type PaymentTransaction struct {
	gorm.Model
	OrderID       uint              `gorm:"index;not null"`
	Order         Order             `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TransactionID string            `gorm:"uniqueIndex;not null;size:255"` // External transaction ID from payment provider
	Type          TransactionType   `gorm:"not null;size:50"`              // Type of transaction (authorize, capture, refund, cancel)
	Status        TransactionStatus `gorm:"not null;size:50"`              // Status of the transaction
	Amount        int64             `gorm:"not null"`                      // Amount of the transaction
	Currency      string            `gorm:"not null;size:3"`               // Currency of the transaction
	Provider      string            `gorm:"not null;size:100"`             // Payment provider (stripe, paypal, etc.)
	RawResponse   string            `gorm:"type:text"`                     // Raw response from payment provider (JSON)
	Metadata      common.StringMap  `gorm:"type:text"`                     // Additional metadata stored as JSON
}

// NewPaymentTransaction creates a new payment transaction
func NewPaymentTransaction(
	orderID uint,
	transactionID string,
	transactionType TransactionType,
	status TransactionStatus,
	amount int64,
	currency string,
	provider string,
) (*PaymentTransaction, error) {
	if orderID == 0 {
		return nil, errors.New("orderID cannot be zero")
	}
	if transactionID == "" {
		return nil, errors.New("transactionID cannot be empty")
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
		OrderID:       orderID,
		TransactionID: transactionID,
		Type:          transactionType,
		Status:        status,
		Amount:        amount,
		Currency:      currency,
		Provider:      provider,
		Metadata:      make(common.StringMap),
	}, nil
}

// AddMetadata adds metadata to the transaction
func (pt *PaymentTransaction) AddMetadata(key, value string) {
	if pt.Metadata == nil {
		pt.Metadata = make(common.StringMap)
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
