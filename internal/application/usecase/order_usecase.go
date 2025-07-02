package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
)

// OrderUseCase implements order-related use cases
type OrderUseCase struct {
	orderRepo          repository.OrderRepository
	productRepo        repository.ProductRepository
	productVariantRepo repository.ProductVariantRepository
	userRepo           repository.UserRepository
	paymentSvc         service.PaymentService
	emailSvc           service.EmailService
	paymentTxnRepo     repository.PaymentTransactionRepository
	currencyRepo       repository.CurrencyRepository
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	productVariantRepo repository.ProductVariantRepository,
	userRepo repository.UserRepository,
	paymentSvc service.PaymentService,
	emailSvc service.EmailService,
	paymentTxnRepo repository.PaymentTransactionRepository,
	currencyRepo repository.CurrencyRepository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:          orderRepo,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		userRepo:           userRepo,
		paymentSvc:         paymentSvc,
		emailSvc:           emailSvc,
		paymentTxnRepo:     paymentTxnRepo,
		currencyRepo:       currencyRepo,
	}
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (uc *OrderUseCase) GetAvailablePaymentProviders() []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProviders()
}

// GetAvailablePaymentProvidersForCurrency returns a list of available payment providers that support the given currency
func (uc *OrderUseCase) GetAvailablePaymentProvidersForCurrency(currency string) []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProvidersForCurrency(currency)
}

// UpdateOrderStatusInput contains the data needed to update an order status
type UpdateOrderStatusInput struct {
	OrderID uint               `json:"order_id"`
	Status  entity.OrderStatus `json:"status"`
}

// UpdateOrderStatus updates the status of an order
func (uc *OrderUseCase) UpdateOrderStatus(input UpdateOrderStatusInput) (*entity.Order, error) {
	// Get order
	order, err := uc.orderRepo.GetByID(input.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Update status
	if err := order.UpdateStatus(input.Status); err != nil {
		return nil, err
	}

	// Update order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (uc *OrderUseCase) GetOrderByID(id uint) (*entity.Order, error) {
	if id == 0 {
		return nil, errors.New("order ID cannot be 0")
	}

	order, err := uc.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByPaymentID retrieves an order by its payment ID
func (uc *OrderUseCase) GetOrderByPaymentID(paymentID string) (*entity.Order, error) {
	if paymentID == "" {
		return nil, errors.New("payment ID cannot be empty")
	}

	// Delegate to the order repository which has this functionality
	order, err := uc.orderRepo.GetByPaymentID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by payment ID: %w", err)
	}

	return order, nil
}

func (uc *OrderUseCase) GetOrderByExternalID(externalID string) (*entity.Order, error) {
	if externalID == "" {
		return nil, errors.New("external ID cannot be empty")
	}

	// Extract order ID from the reference
	var orderID uint
	_, err := fmt.Sscanf(externalID, "order-%d-", &orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid reference format in MobilePay webhook event: %s", externalID)
	}

	fmt.Printf("Extracted order ID from external ID: %d\n", orderID)

	// Delegate to the order repository which has this functionality
	order, err := uc.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by external ID: %w", err)
	}

	return order, nil
}

// GetUserOrders retrieves orders for a user
func (uc *OrderUseCase) GetUserOrders(userID uint, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.GetByUser(userID, offset, limit)
}

func (uc *OrderUseCase) ListOrdersByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.ListByStatus(status, offset, limit)
}

func (uc *OrderUseCase) FailOrder(order *entity.Order) error {
	// Update the payment status to failed, which will also update order status to cancelled
	if err := order.UpdatePaymentStatus(entity.PaymentStatusFailed); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Save the updated order in the repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save updated order: %w", err)
	}

	return nil
}

// CapturePayment captures an authorized payment
func (uc *OrderUseCase) CapturePayment(transactionID string, amount int64) error {
	// Find the order with this payment ID
	order, err := uc.orderRepo.GetByPaymentID(transactionID)
	if err != nil {
		return errors.New("order not found for payment ID")
	}

	// Check if the payment is already captured
	if order.PaymentStatus == entity.PaymentStatusCaptured {
		return errors.New("payment already captured for this order")
	}

	// Check if the payment is in authorized state and order is shipped (new rule)
	if order.PaymentStatus != entity.PaymentStatusAuthorized {
		return errors.New("payment must be authorized before capture")
	}

	if order.Status != entity.OrderStatusShipped {
		return errors.New("order must be shipped before payment can be captured")
	}

	// Check if the amount is valid
	if amount <= 0 {
		return errors.New("capture amount must be greater than zero")
	}

	// Check if amount is greater than the order amount
	if amount > order.FinalAmount {
		return errors.New("capture amount cannot exceed the original payment amount")
	}

	providerType := common.PaymentProviderType(order.PaymentProvider)

	// Call payment service to capture payment
	_, err = uc.paymentSvc.CapturePayment(transactionID, order.Currency, amount, providerType)
	if err != nil {
		// Record failed capture attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			"", // Idempotency key
			entity.TransactionTypeCapture,
			entity.TransactionStatusFailed,
			amount,
			order.Currency,
			string(providerType),
		)

		if txErr == nil {
			txn.AddMetadata("error", err.Error())
			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				log.Printf("Failed to save capture transaction: %v\n", err)
			}
		}

		return fmt.Errorf("failed to capture payment: %v", err)
	}

	// Update payment status to captured, which will also update order status to completed
	// if err := order.UpdatePaymentStatus(entity.PaymentStatusCaptured); err != nil {
	// 	return fmt.Errorf("failed to update payment status: %v", err)
	// }

	// Save the updated order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save order status: %v", err)
	}

	// Stock was already decreased when payment was authorized, no need to decrease again

	// Record successful capture transaction
	// Track if this is a full or partial capture
	isFullCapture := amount >= order.FinalAmount

	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		"", // Idempotency key
		entity.TransactionTypeCapture,
		entity.TransactionStatusSuccessful,
		amount,
		order.Currency,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("full_capture", fmt.Sprintf("%t", isFullCapture))

		// Record total authorized amount
		if isFullCapture {
			txn.AddMetadata("remaining_amount", "0")
		} else {
			remainingAmount := order.FinalAmount - amount
			txn.AddMetadata("remaining_amount", fmt.Sprintf("%.2f", money.FromCents(remainingAmount)))
		}

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save capture transaction: %v\n", err)
		}
	}

	return nil
}

// CancelPayment cancels a payment
func (uc *OrderUseCase) CancelPayment(transactionID string) error {
	// Find the order with this payment ID
	order, err := uc.orderRepo.GetByPaymentID(transactionID)
	if err != nil {
		return errors.New("order not found for payment ID")
	}

	// Check if the payment is already cancelled
	if order.PaymentStatus == entity.PaymentStatusCancelled {
		return errors.New("payment already canceled")
	}

	// Check if the payment is in authorized state (can only cancel authorized payments that aren't captured)
	if order.PaymentStatus != entity.PaymentStatusAuthorized {
		return errors.New("payment cancellation only allowed for authorized payments")
	}

	// Check if the transaction ID is valid
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	providerType := common.PaymentProviderType(order.PaymentProvider)

	_, err = uc.paymentSvc.CancelPayment(transactionID, providerType)
	if err != nil {
		// Record failed cancellation attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			"", // Idempotency key
			entity.TransactionTypeCancel,
			entity.TransactionStatusFailed,
			0, // No amount for cancellation
			order.Currency,
			string(providerType),
		)
		if txErr == nil {
			txn.AddMetadata("error", err.Error())
			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				log.Printf("Failed to save cancel transaction: %v\n", err)
			}
		}

		return fmt.Errorf("failed to cancel payment: %v", err)
	}

	// Update payment status to cancelled, which will also update order status to cancelled
	if err := order.UpdatePaymentStatus(entity.PaymentStatusCancelled); err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	// Save the updated order in the repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save order status: %v", err)
	}

	// Record successful cancellation transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		"", // Idempotency key
		entity.TransactionTypeCancel,
		entity.TransactionStatusSuccessful,
		0, // No amount for cancellation
		order.Currency,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("previous_order_status", string(order.Status))
		txn.AddMetadata("previous_payment_status", string(entity.PaymentStatusAuthorized))

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save cancel transaction: %v\n", err)
		}
	}

	return nil
}

// RefundPayment refunds a payment
func (uc *OrderUseCase) RefundPayment(transactionID string, amount int64) error {
	// Find the order with this payment ID
	order, err := uc.orderRepo.GetByPaymentID(transactionID)
	if err != nil {
		return errors.New("order not found for payment ID")
	}

	// Check if the payment is in a state that allows refund (authorized, captured, or partially refunded)
	if order.PaymentStatus != entity.PaymentStatusAuthorized &&
		order.PaymentStatus != entity.PaymentStatusCaptured &&
		order.PaymentStatus != entity.PaymentStatusRefunded {
		return errors.New("payment refund only allowed for authorized, captured, or partially refunded payments")
	}

	// Check if the amount is valid
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Check if the refund amount exceeds the original amount
	if amount > order.FinalAmount {
		return errors.New("refund amount cannot exceed the original payment amount")
	}

	providerType := common.PaymentProviderType(order.PaymentProvider)

	// Get the total captured amount (what's available to refund)
	totalCapturedAmount, err := uc.paymentTxnRepo.SumCapturedAmountByOrderID(order.ID)
	if err != nil {
		return fmt.Errorf("failed to get captured amount: %w", err)
	}

	// If no amount has been captured, we can't refund
	if totalCapturedAmount == 0 {
		return errors.New("no captured amount available for refund")
	}

	// Get total refunded amount so far (if any)
	totalRefundedSoFar, err := uc.paymentTxnRepo.SumRefundedAmountByOrderID(order.ID)
	if err != nil {
		return fmt.Errorf("failed to get refunded amount: %w", err)
	}

	// Check if the payment has already been fully refunded
	if totalRefundedSoFar >= totalCapturedAmount {
		return errors.New("payment has already been fully refunded")
	}

	// Check if we're trying to refund more than the remaining amount
	remainingAmount := totalCapturedAmount - totalRefundedSoFar
	if amount > remainingAmount {
		return fmt.Errorf("refund amount (%d) would exceed remaining refundable amount (%d)", amount, remainingAmount)
	}

	_, err = uc.paymentSvc.RefundPayment(transactionID, order.Currency, amount, providerType)
	if err != nil {
		// Record failed refund attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			"", // Idempotency key
			entity.TransactionTypeRefund,
			entity.TransactionStatusFailed,
			amount,
			order.Currency,
			string(providerType),
		)
		if txErr == nil {
			txn.AddMetadata("error", err.Error())
			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				log.Printf("Failed to save refund transaction: %v\n", err)
			}
		}

		return fmt.Errorf("failed to refund payment: %v", err)
	}

	// Calculate if this is a full refund (refunding all captured amount)
	isFullRefund := (totalRefundedSoFar + amount) >= totalCapturedAmount

	// Only update the payment status to refunded if it's a full refund
	if isFullRefund {
		if err := order.UpdatePaymentStatus(entity.PaymentStatusRefunded); err != nil {
			return fmt.Errorf("failed to update payment status: %v", err)
		}

		// Save the updated order in the repository
		if err := uc.orderRepo.Update(order); err != nil {
			return fmt.Errorf("failed to save order status: %v", err)
		}
	}

	// Record successful refund transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		"", // Idempotency key
		entity.TransactionTypeRefund,
		entity.TransactionStatusSuccessful,
		amount,
		order.Currency,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("full_refund", fmt.Sprintf("%t", isFullRefund))
		txn.AddMetadata("previous_payment_status", string(order.PaymentStatus))

		// Record total refunded amount including this transaction
		totalRefunded := totalRefundedSoFar + amount
		txn.AddMetadata("total_refunded", fmt.Sprintf("%.2f", money.FromCents(totalRefunded)))

		// Record remaining amount still available for refund
		remainingAmount := max(order.FinalAmount-totalRefunded, 0)
		txn.AddMetadata("remaining_available", fmt.Sprintf("%.2f", money.FromCents(remainingAmount)))

		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			log.Printf("Failed to save refund transaction: %v\n", err)
		}
	}

	return nil
}

func (uc *OrderUseCase) UpdatePaymentTransaction(transactionID string, status entity.TransactionStatus, metadata map[string]string) error {
	txn, err := uc.paymentTxnRepo.GetByTransactionID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get payment transaction: %w", err)
	}

	txn.UpdateStatus(status)

	for key, value := range metadata {
		txn.AddMetadata(key, value)
	}

	return uc.paymentTxnRepo.Update(txn)
}

// ForceApproveMobilePayPayment force approves a MobilePay payment
func (uc *OrderUseCase) ForceApproveMobilePayPayment(paymentID string, phoneNumber string) error {
	// Get the payment service
	paymentSvc, ok := uc.paymentSvc.(*payment.MultiProviderPaymentService)
	if !ok {
		return errors.New("invalid payment service")
	}

	// Force approve the payment
	return paymentSvc.ForceApprovePayment(paymentID, phoneNumber, common.PaymentProviderMobilePay)
}

// GetUserByID retrieves a user by ID
func (uc *OrderUseCase) GetUserByID(id uint) (*entity.User, error) {
	return uc.userRepo.GetByID(id)
}

// ListAllOrders lists all orders
func (uc *OrderUseCase) ListAllOrders(offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.ListAll(offset, limit)
}

// RecordPaymentTransaction records a payment transaction for an order
func (uc *OrderUseCase) RecordPaymentTransaction(transaction *entity.PaymentTransaction) error {
	if transaction == nil {
		return errors.New("payment transaction cannot be nil")
	}

	// Validate the order exists
	_, err := uc.orderRepo.GetByID(transaction.OrderID)
	if err != nil {
		return fmt.Errorf("failed to verify order existence: %w", err)
	}

	// Create transaction record
	return uc.paymentTxnRepo.Create(transaction)
}

// GetTransactionByTransactionID retrieves a payment transaction by its transaction ID
func (uc *OrderUseCase) GetTransactionByTransactionID(transactionID string) (*entity.PaymentTransaction, error) {
	return uc.paymentTxnRepo.GetByTransactionID(transactionID)
}

// GetTransactionByIdempotencyKey retrieves a payment transaction by its idempotency key
func (uc *OrderUseCase) GetTransactionByIdempotencyKey(idempotencyKey string) (*entity.PaymentTransaction, error) {
	return uc.paymentTxnRepo.GetByIdempotencyKey(idempotencyKey)
}

// GetLatestPendingTransactionByType retrieves the latest pending transaction of a specific type for an order
func (uc *OrderUseCase) GetLatestPendingTransactionByType(orderID uint, txnType entity.TransactionType) (*entity.PaymentTransaction, error) {
	// Get all transactions for the order
	transactions, err := uc.paymentTxnRepo.GetByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for order %d: %w", orderID, err)
	}

	// Find the latest pending transaction of the specified type
	var latestPending *entity.PaymentTransaction
	for _, txn := range transactions {
		if txn.Type == txnType && txn.Status == entity.TransactionStatusPending {
			if latestPending == nil || txn.CreatedAt.After(latestPending.CreatedAt) {
				latestPending = txn
			}
		}
	}

	if latestPending == nil {
		return nil, fmt.Errorf("no pending transaction of type %s found for order %d", txnType, orderID)
	}

	return latestPending, nil
}

// UpdatePaymentTransactionStatus updates an existing transaction's status and metadata
func (uc *OrderUseCase) UpdatePaymentTransactionStatus(transaction *entity.PaymentTransaction, status entity.TransactionStatus, rawResponse string, metadata map[string]string) error {
	// Update status using the proper method that handles amount field updates
	transaction.UpdateStatus(status)
	
	if rawResponse != "" {
		transaction.SetRawResponse(rawResponse)
	}

	// Add any new metadata
	for key, value := range metadata {
		transaction.AddMetadata(key, value)
	}

	// Save the updated transaction
	return uc.paymentTxnRepo.Update(transaction)
}

// UpdatePaymentStatusInput contains the data needed to update payment status
type UpdatePaymentStatusInput struct {
	OrderID       uint
	PaymentStatus entity.PaymentStatus
	TransactionID string // Optional, for logging purposes
}

// UpdatePaymentStatus updates the payment status of an order
func (uc *OrderUseCase) UpdatePaymentStatus(input UpdatePaymentStatusInput) (*entity.Order, error) {
	// Get order
	order, err := uc.orderRepo.GetByID(input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Store the previous payment status to determine if stock updates are needed
	previousPaymentStatus := order.PaymentStatus

	// Update payment status
	if err := order.UpdatePaymentStatus(input.PaymentStatus); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// Update order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Handle stock updates based on payment status transitions
	if err := uc.handleStockUpdatesForPaymentStatusChange(order, previousPaymentStatus, input.PaymentStatus); err != nil {
		// Log the error but don't fail the status update since the payment status change was successful
		log.Printf("Warning: Failed to update stock for order %d: %v", order.ID, err)
	}

	// Send emails for payment status changes
	if err := uc.handleEmailsForPaymentStatusChange(order, previousPaymentStatus, input.PaymentStatus); err != nil {
		// Log the error but don't fail the status update since the payment status change was successful
		log.Printf("Warning: Failed to send emails for order %d: %v", order.ID, err)
	}

	return order, nil
}

// handleStockUpdatesForPaymentStatusChange handles stock updates when payment status changes
func (uc *OrderUseCase) handleStockUpdatesForPaymentStatusChange(order *entity.Order, previousStatus, newStatus entity.PaymentStatus) error {
	// Only handle stock changes for specific transitions
	switch {
	case previousStatus != entity.PaymentStatusAuthorized && newStatus == entity.PaymentStatusAuthorized:
		// Payment was just authorized - decrease stock to reserve items
		return uc.decreaseStock(order)

	case previousStatus == entity.PaymentStatusAuthorized && newStatus == entity.PaymentStatusCancelled:
		// Payment was authorized but now cancelled - restore stock
		return uc.increaseStock(order)

	case previousStatus == entity.PaymentStatusAuthorized && newStatus == entity.PaymentStatusFailed:
		// Payment was authorized but now failed - restore stock
		return uc.increaseStock(order)

	case previousStatus == entity.PaymentStatusCaptured && newStatus == entity.PaymentStatusRefunded:
		// Payment was captured but now refunded - restore stock
		return uc.increaseStock(order)

	case previousStatus != entity.PaymentStatusCancelled && newStatus == entity.PaymentStatusCancelled && previousStatus != entity.PaymentStatusAuthorized:
		// Payment was cancelled without being authorized first - no stock change needed
		return nil

	case previousStatus != entity.PaymentStatusFailed && newStatus == entity.PaymentStatusFailed && previousStatus != entity.PaymentStatusAuthorized:
		// Payment failed without being authorized first - no stock change needed
		return nil

	default:
		// No stock change needed for other transitions (e.g., authorized -> captured)
		return nil
	}
}

// decreaseStock decreases stock for all items in an order
func (uc *OrderUseCase) decreaseStock(order *entity.Order) error {
	for _, item := range order.Items {
		// Skip items without variant ID (shouldn't happen, but safety check)
		if item.ProductVariantID == 0 {
			continue
		}

		// Get the variant
		variant, err := uc.productVariantRepo.GetByID(item.ProductVariantID)
		if err != nil {
			return fmt.Errorf("failed to get variant %d: %w", item.ProductVariantID, err)
		}

		// Check if there's enough stock
		if variant.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s (SKU: %s): available %d, required %d",
				item.ProductName, item.SKU, variant.Stock, item.Quantity)
		}

		// Update stock
		changeAmount := -item.Quantity // Negative because we're decreasing
		if err := variant.UpdateStock(changeAmount); err != nil {
			return fmt.Errorf("failed to update stock for variant %d: %w", item.ProductVariantID, err)
		}

		// Save the updated variant
		if err := uc.productVariantRepo.Update(variant); err != nil {
			return fmt.Errorf("failed to save variant %d: %w", item.ProductVariantID, err)
		}
	}
	return nil
}

// increaseStock increases stock for all items in an order (for cancellations/refunds)
func (uc *OrderUseCase) increaseStock(order *entity.Order) error {
	for _, item := range order.Items {
		// Skip items without variant ID (shouldn't happen, but safety check)
		if item.ProductVariantID == 0 {
			continue
		}

		// Get the variant
		variant, err := uc.productVariantRepo.GetByID(item.ProductVariantID)
		if err != nil {
			return fmt.Errorf("failed to get variant %d: %w", item.ProductVariantID, err)
		}

		// Update stock
		changeAmount := item.Quantity // Positive because we're increasing
		if err := variant.UpdateStock(changeAmount); err != nil {
			return fmt.Errorf("failed to update stock for variant %d: %w", item.ProductVariantID, err)
		}

		// Save the updated variant
		if err := uc.productVariantRepo.Update(variant); err != nil {
			return fmt.Errorf("failed to save variant %d: %w", item.ProductVariantID, err)
		}
	}
	return nil
}

// handleEmailsForPaymentStatusChange sends appropriate emails when payment status changes
func (uc *OrderUseCase) handleEmailsForPaymentStatusChange(order *entity.Order, previousStatus, newStatus entity.PaymentStatus) error {
	// Only send emails when payment status changes to authorized or paid
	shouldSendEmails := false

	switch {
	case previousStatus != entity.PaymentStatusAuthorized && newStatus == entity.PaymentStatusAuthorized:
		// Payment was just authorized - send order confirmation and notification emails
		shouldSendEmails = true
	case previousStatus != entity.PaymentStatusCaptured && newStatus == entity.PaymentStatusCaptured:
		// Payment was just captured/paid - send order confirmation and notification emails
		shouldSendEmails = true
	default:
		// No emails needed for other transitions
		return nil
	}

	if !shouldSendEmails {
		return nil
	}

	// Create user object for email sending
	var user *entity.User
	if order.IsGuestOrder || order.UserID == nil {
		// Guest order - create a temporary user object with customer details
		if order.CustomerDetails == nil {
			return fmt.Errorf("guest order missing customer details")
		}
		user = &entity.User{
			Email:     order.CustomerDetails.Email,
			FirstName: order.CustomerDetails.FullName, // Use FullName as FirstName for guest orders
		}
	} else {
		// Registered user - get from repository
		var err error
		user, err = uc.userRepo.GetByID(*order.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user %d: %w", *order.UserID, err)
		}
	}

	// Send order confirmation email to customer
	if err := uc.emailSvc.SendOrderConfirmation(order, user); err != nil {
		return fmt.Errorf("failed to send order confirmation email: %w", err)
	}

	// Send order notification email to admin
	if err := uc.emailSvc.SendOrderNotification(order, user); err != nil {
		return fmt.Errorf("failed to send order notification email: %w", err)
	}

	log.Printf("Sent order confirmation and notification emails for order %d (status: %s)", order.ID, newStatus)
	return nil
}
