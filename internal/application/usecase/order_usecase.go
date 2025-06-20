package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
)

// OrderUseCase implements order-related use cases
type OrderUseCase struct {
	orderRepo      repository.OrderRepository
	productRepo    repository.ProductRepository
	userRepo       repository.UserRepository
	paymentSvc     service.PaymentService
	emailSvc       service.EmailService
	paymentTxnRepo repository.PaymentTransactionRepository
	currencyRepo   repository.CurrencyRepository
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	paymentSvc service.PaymentService,
	emailSvc service.EmailService,
	paymentTxnRepo repository.PaymentTransactionRepository,
	currencyRepo repository.CurrencyRepository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		userRepo:       userRepo,
		paymentSvc:     paymentSvc,
		emailSvc:       emailSvc,
		paymentTxnRepo: paymentTxnRepo,
		currencyRepo:   currencyRepo,
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

// GetUserOrders retrieves orders for a user
func (uc *OrderUseCase) GetUserOrders(userID uint, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.GetByUser(userID, offset, limit)
}

func (uc *OrderUseCase) ListOrdersByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	return uc.orderRepo.ListByStatus(status, offset, limit)
}

func (uc *OrderUseCase) FailOrder(order *entity.Order) error {
	// Update the order status to failed
	if err := order.UpdateStatus(entity.OrderStatusFailed); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
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

	// Check if the order is already captured
	if order.Status == entity.OrderStatusCaptured {
		return errors.New("payment already captured")
	}
	// Check if the order is in a state that allows capture
	if order.Status != entity.OrderStatusPaid {
		return errors.New("payment capture not allowed in current order status")
	}

	// Check if the amount is valid
	if amount <= 0 {
		return errors.New("capture amount must be greater than zero")
	}

	// Check if amount is greater than the order amount
	if amount > order.FinalAmount {
		return errors.New("capture amount cannot exceed the original payment amount")
	}

	providerType := service.PaymentProviderType(order.PaymentProvider)

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	// Call payment service to capture payment
	err = uc.paymentSvc.CapturePayment(transactionID, amount, providerType)
	if err != nil {
		// Record failed capture attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			entity.TransactionTypeCapture,
			entity.TransactionStatusFailed,
			amount,
			defaultCurrency.Code,
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

	if err := order.UpdateStatus(entity.OrderStatusCaptured); err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// Save the updated order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save order status: %v", err)
	}

	// Record successful capture transaction
	// Track if this is a full or partial capture
	isFullCapture := amount >= order.FinalAmount

	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		entity.TransactionTypeCapture,
		entity.TransactionStatusSuccessful,
		amount,
		defaultCurrency.Code,
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

	// Check if the order is already canceled
	if order.Status == entity.OrderStatusCancelled {
		return errors.New("payment already canceled")
	}
	// Check if the order is in a state that allows cancellation
	if order.Status != entity.OrderStatusPendingAction {
		return errors.New("payment cancellation not allowed in current order status")
	}
	// Check if the transaction ID is valid
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	providerType := service.PaymentProviderType(order.PaymentProvider)

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	err = uc.paymentSvc.CancelPayment(transactionID, providerType)
	if err != nil {
		// Record failed cancellation attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			entity.TransactionTypeCancel,
			entity.TransactionStatusFailed,
			0, // No amount for cancellation
			defaultCurrency.Code,
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

	// Update the order status to cancelled after successful payment cancellation
	if err := order.UpdateStatus(entity.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// Save the updated order in the repository
	if err := uc.orderRepo.Update(order); err != nil {
		return fmt.Errorf("failed to save order status: %v", err)
	}

	// Record successful cancellation transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		transactionID,
		entity.TransactionTypeCancel,
		entity.TransactionStatusSuccessful,
		0, // No amount for cancellation
		defaultCurrency.Code,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("previous_status", string(entity.OrderStatusPendingAction))

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

	// Check if the order is already refunded
	if order.Status == entity.OrderStatusRefunded {
		return errors.New("payment already refunded")
	}
	// Check if the order is in a state that allows refund
	if order.Status != entity.OrderStatusPaid && order.Status != entity.OrderStatusCaptured {
		return errors.New("payment refund not allowed in current order status")
	}
	// Check if the amount is valid
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Check if the refund amount exceeds the original amount
	if amount > order.FinalAmount {
		return errors.New("refund amount cannot exceed the original payment amount")
	}

	providerType := service.PaymentProviderType(order.PaymentProvider)

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return fmt.Errorf("failed to get default currency: %w", err)
	}

	// Get total refunded amount so far (if any)
	var totalRefundedSoFar int64 = 0
	totalRefundedSoFar, _ = uc.paymentTxnRepo.SumAmountByOrderIDAndType(order.ID, entity.TransactionTypeRefund)

	// Check if we're trying to refund more than the original amount when combining with previous refunds
	if totalRefundedSoFar+amount > order.FinalAmount {
		return errors.New("total refund amount would exceed the original payment amount")
	}

	err = uc.paymentSvc.RefundPayment(transactionID, amount, providerType)
	if err != nil {
		// Record failed refund attempt
		txn, txErr := entity.NewPaymentTransaction(
			order.ID,
			transactionID,
			entity.TransactionTypeRefund,
			entity.TransactionStatusFailed,
			amount,
			defaultCurrency.Code,
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

	// Calculate if this is a full refund
	isFullRefund := false
	if amount >= order.FinalAmount || (totalRefundedSoFar+amount) >= order.FinalAmount {
		isFullRefund = true
	}

	// Only update the order status to refunded if it's a full refund
	if isFullRefund {
		if err := order.UpdateStatus(entity.OrderStatusRefunded); err != nil {
			return fmt.Errorf("failed to update order status: %v", err)
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
		entity.TransactionTypeRefund,
		entity.TransactionStatusSuccessful,
		amount,
		defaultCurrency.Code,
		string(providerType),
	)
	if err == nil {
		txn.AddMetadata("full_refund", fmt.Sprintf("%t", isFullRefund))
		txn.AddMetadata("previous_status", string(order.Status))

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
	return paymentSvc.ForceApprovePayment(paymentID, phoneNumber, service.PaymentProviderMobilePay)
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
