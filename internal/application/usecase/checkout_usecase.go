package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
)

// CheckoutInput defines the input for creating/adding to a checkout
type CheckoutInput struct {
	SKU         string
	Quantity    int
	Price       int64
	Weight      float64
	ProductName string
	VariantName string
	ProductID   uint // Internal use only - resolved from SKU
	VariantID   uint // Internal use only - resolved from SKU
}

// UpdateCheckoutItemInput defines the input for updating a checkout item
type UpdateCheckoutItemInput struct {
	SKU      string
	Quantity int
}

// RemoveItemInput defines the input for removing an item from a checkout
type RemoveItemInput struct {
	SKU string
}

// CheckoutUseCase implements checkout business logic
type CheckoutUseCase struct {
	checkoutRepo       repository.CheckoutRepository
	productRepo        repository.ProductRepository
	productVariantRepo repository.ProductVariantRepository
	shippingMethodRepo repository.ShippingMethodRepository
	shippingRateRepo   repository.ShippingRateRepository
	discountRepo       repository.DiscountRepository
	orderRepo          repository.OrderRepository
	currencyRepo       repository.CurrencyRepository
	paymentTxnRepo     repository.PaymentTransactionRepository
	paymentSvc         service.PaymentService
	shippingUsecase    *ShippingUseCase
}

type ProcessPaymentInput struct {
	PaymentProvider service.PaymentProviderType
	PaymentMethod   service.PaymentMethod
	CardDetails     *service.CardDetails `json:"card_details,omitempty"`
	PhoneNumber     string               `json:"phone_number,omitempty"`
}

func (uc *CheckoutUseCase) ProcessPayment(order *entity.Order, input ProcessPaymentInput) (*entity.Order, error) {
	// Validate order
	if order == nil {
		return nil, errors.New("order cannot be nil")
	}

	if order.ID == 0 {
		return nil, errors.New("order ID is required")
	}

	if order.Status != entity.OrderStatusPending {
		return nil, errors.New("order is not in a valid state for payment processing")
	}

	if order.CustomerDetails == nil || order.CustomerDetails.Email == "" {
		return nil, errors.New("customer details are required for payment processing")
	}

	if order.ShippingMethodID == 0 {
		return nil, errors.New("shipping method is required for payment processing")
	}

	// Check if order is already paid
	if order.Status == entity.OrderStatusPaid ||
		order.Status == entity.OrderStatusShipped ||
		order.Status == entity.OrderStatusDelivered {
		return nil, errors.New("order is already paid")
	}

	// Get default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to get default currency: %w", err)
	}

	// Validate payment provider supports the currency
	availableProviders := uc.GetAvailablePaymentProvidersForCurrency(defaultCurrency.Code)
	providerValid := false
	for _, p := range availableProviders {
		if p.Type == input.PaymentProvider && p.Enabled {
			providerValid = true
			break
		}
	}
	if !providerValid {
		return nil, fmt.Errorf("payment provider %s does not support currency %s", input.PaymentProvider, defaultCurrency.Code)
	}

	// Process payment
	paymentResult, err := uc.paymentSvc.ProcessPayment(service.PaymentRequest{
		OrderID:         order.ID,
		Amount:          order.FinalAmount, // Use final amount (after discounts)
		Currency:        defaultCurrency.Code,
		PaymentMethod:   input.PaymentMethod,
		PaymentProvider: input.PaymentProvider,
		CardDetails:     input.CardDetails,
		PhoneNumber:     input.PhoneNumber,
		CustomerEmail:   order.CustomerDetails.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	if paymentResult.RequiresAction && paymentResult.ActionURL != "" {
		// Update order with payment ID, provider, and status
		if err := order.SetPaymentID(paymentResult.TransactionID); err != nil {
			return nil, err
		}
		if err := order.SetPaymentProvider(string(paymentResult.Provider)); err != nil {
			return nil, err
		}
		if err := order.SetActionURL(paymentResult.ActionURL); err != nil {
			return nil, err
		}
		if err := order.UpdateStatus(entity.OrderStatusPendingAction); err != nil {
			return nil, err
		}

		// Update order in repository
		if err := uc.orderRepo.Update(order); err != nil {
			return nil, err
		}

		// Record the pending authorization transaction
		txn, err := entity.NewPaymentTransaction(
			order.ID,
			paymentResult.TransactionID,
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			order.FinalAmount,
			defaultCurrency.Code,
			string(paymentResult.Provider),
		)
		if err != nil {
			// Log the error but don't fail the payment process
			log.Printf("Failed to create payment transaction record: %v", err)
		} else {
			// Add metadata
			txn.AddMetadata("payment_method", string(order.PaymentMethod))
			txn.AddMetadata("requires_action", "true")
			txn.AddMetadata("action_url", paymentResult.ActionURL)

			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				// Log error but don't fail the payment process
				log.Printf("Failed to save payment transaction: %v\n", err)
			}
		}

		return order, nil
	}

	if !paymentResult.Success {
		// Record the failed transaction
		txn, err := entity.NewPaymentTransaction(
			order.ID,
			paymentResult.TransactionID,
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusFailed,
			order.FinalAmount,
			defaultCurrency.Code,
			string(paymentResult.Provider),
		)
		if err == nil {
			txn.AddMetadata("payment_method", string(order.PaymentMethod))
			txn.AddMetadata("error_message", paymentResult.ErrorMessage)

			if err := uc.paymentTxnRepo.Create(txn); err != nil {
				// Log error but don't fail the process
				log.Printf("Failed to save failed payment transaction: %v\n", err)
			}
		}

		return nil, errors.New(paymentResult.ErrorMessage)
	}

	// Update order with payment ID, provider, and status
	if err := order.SetPaymentID(paymentResult.TransactionID); err != nil {
		return nil, err
	}
	if err := order.SetPaymentProvider(string(paymentResult.Provider)); err != nil {
		return nil, err
	}
	if err := order.SetPaymentMethod(string(order.PaymentMethod)); err != nil {
		return nil, err
	}
	if err := order.UpdateStatus(entity.OrderStatusPaid); err != nil {
		return nil, err
	}

	// Update order in repository
	if err := uc.orderRepo.Update(order); err != nil {
		return nil, err
	}

	// Record the successful authorization transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		paymentResult.TransactionID,
		entity.TransactionTypeAuthorize,
		entity.TransactionStatusSuccessful,
		order.FinalAmount,
		defaultCurrency.Code,
		string(paymentResult.Provider),
	)
	if err != nil {
		// Log the error but don't fail the payment process
		log.Printf("Failed to create payment transaction record: %v", err)
	} else {
		if err := uc.paymentTxnRepo.Create(txn); err != nil {
			// Log error but don't fail the payment process
			log.Printf("Failed to save payment transaction: %v\n", err)
		}
	}

	return order, nil
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (uc *CheckoutUseCase) GetAvailablePaymentProviders() []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProviders()
}

// GetAvailablePaymentProvidersForCurrency returns a list of available payment providers that support the given currency
func (uc *CheckoutUseCase) GetAvailablePaymentProvidersForCurrency(currency string) []service.PaymentProvider {
	return uc.paymentSvc.GetAvailableProvidersForCurrency(currency)
}

// NewCheckoutUseCase creates a new checkout use case
func NewCheckoutUseCase(
	checkoutRepo repository.CheckoutRepository,
	productRepo repository.ProductRepository,
	productVariantRepo repository.ProductVariantRepository,
	shippingMethodRepo repository.ShippingMethodRepository,
	shippingRateRepo repository.ShippingRateRepository,
	discountRepo repository.DiscountRepository,
	orderRepo repository.OrderRepository,
	currencyRepo repository.CurrencyRepository,
	paymentTxnRepo repository.PaymentTransactionRepository,
	paymentSvc service.PaymentService,
	shippingUsecase *ShippingUseCase,

) *CheckoutUseCase {
	return &CheckoutUseCase{
		checkoutRepo:       checkoutRepo,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		shippingMethodRepo: shippingMethodRepo,
		shippingRateRepo:   shippingRateRepo,
		discountRepo:       discountRepo,
		orderRepo:          orderRepo,
		paymentTxnRepo:     paymentTxnRepo,
		currencyRepo:       currencyRepo,
		paymentSvc:         paymentSvc,
		shippingUsecase:    shippingUsecase,
	}
}

// GetOrCreateCheckout retrieves or creates a checkout for a user
func (uc *CheckoutUseCase) GetOrCreateCheckout(sessionId string) (*entity.Checkout, error) {
	// If not found, create a new one
	checkout, err := entity.NewCheckout(sessionId)
	if err != nil {
		return nil, err
	}

	// Set default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err == nil && defaultCurrency != nil {
		checkout.Currency = defaultCurrency.Code
	}

	// Save to repository
	err = uc.checkoutRepo.Create(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetShippingAddress sets the shipping address for the user's checkout
func (uc *CheckoutUseCase) SetShippingAddress(userID uint, address entity.Address) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set shipping address
	checkout.SetShippingAddress(address)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetBillingAddress sets the billing address for the user's checkout
func (uc *CheckoutUseCase) SetBillingAddress(userID uint, address entity.Address) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set billing address
	checkout.SetBillingAddress(address)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetCustomerDetails sets the customer details for the user's checkout
func (uc *CheckoutUseCase) SetCustomerDetails(userID uint, details entity.CustomerDetails) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set customer details
	checkout.SetCustomerDetails(details)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// SetShippingMethod sets the shipping method for the user's checkout
func (uc *CheckoutUseCase) SetShippingMethod(checkout *entity.Checkout, methodID uint) (*entity.Checkout, error) {
	// Validate inputs
	if checkout == nil {
		return nil, errors.New("checkout cannot be nil")
	}

	if methodID == 0 {
		return nil, errors.New("shipping method ID is required")
	}

	// Check if checkout is active
	if checkout.Status != entity.CheckoutStatusActive {
		return nil, errors.New("cannot modify a non-active checkout")
	}

	// Validate shipping address is set
	if checkout.ShippingAddr.Street == "" || checkout.ShippingAddr.Country == "" {
		return nil, errors.New("shipping address is required to calculate shipping options")
	}

	// Verify shipping method exists
	method, err := uc.shippingMethodRepo.GetByID(methodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping method: %w", err)
	}

	// Check if shipping method is active
	if !method.Active {
		return nil, errors.New("shipping method is not available")
	}

	// Calculate shipping options
	options, err := uc.shippingUsecase.CalculateShippingOptions(checkout.ShippingAddr, checkout.TotalAmount, checkout.TotalWeight)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate shipping options: %w", err)
	}

	// Find the selected shipping option
	var selectedOption *entity.ShippingOption
	for _, option := range options.Options {
		if option.ShippingMethodID == methodID {
			selectedOption = option
			break
		}
	}

	if selectedOption == nil {
		return nil, fmt.Errorf("shipping method %d is not available for the current checkout", methodID)
	}

	// Set shipping method and cost
	checkout.SetShippingMethod(selectedOption)

	// Update checkout in repository
	if err := uc.checkoutRepo.Update(checkout); err != nil {
		return nil, fmt.Errorf("failed to update checkout: %w", err)
	}

	return checkout, nil
}

// SetPaymentProvider sets the payment provider for the user's checkout
func (uc *CheckoutUseCase) SetPaymentProvider(userID uint, provider string) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set payment provider
	checkout.SetPaymentProvider(provider)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ApplyDiscountCode applies a discount code to the user's checkout
func (uc *CheckoutUseCase) ApplyDiscountCode(checkout *entity.Checkout, code string) (*entity.Checkout, error) {
	// Get discount
	discount, err := uc.discountRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Check if discount is valid
	if !discount.IsValid() {
		return nil, errors.New("discount is not valid")
	}

	// Apply discount
	checkout.ApplyDiscount(discount)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// RemoveDiscountCode removes a discount code from the user's checkout
func (uc *CheckoutUseCase) RemoveDiscountCode(checkout *entity.Checkout) (*entity.Checkout, error) {
	// Remove discount
	checkout.ApplyDiscount(nil)

	// Update checkout in repository
	err := uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ExpireOldCheckouts marks expired checkouts as expired
func (uc *CheckoutUseCase) ExpireOldCheckouts() error {
	// Get expired checkouts
	expiredCheckouts, err := uc.checkoutRepo.GetExpiredCheckouts()
	if err != nil {
		return err
	}

	// Mark each as expired
	for _, checkout := range expiredCheckouts {
		checkout.MarkAsExpired()
		err = uc.checkoutRepo.Update(checkout)
		if err != nil {
			// Continue despite errors
			continue
		}
	}

	return nil
}

// CreateOrderFromCheckout creates an order from a checkout
func (uc *CheckoutUseCase) CreateOrderFromCheckout(checkoutID uint) (*entity.Order, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Validate checkout
	if len(checkout.Items) == 0 {
		return nil, errors.New("checkout has no items")
	}

	if checkout.ShippingAddr.Street == "" || checkout.ShippingAddr.Country == "" {
		return nil, errors.New("shipping address is required")
	}

	if checkout.BillingAddr.Street == "" || checkout.BillingAddr.Country == "" {
		return nil, errors.New("billing address is required")
	}

	if checkout.CustomerDetails.Email == "" || checkout.CustomerDetails.FullName == "" {
		return nil, errors.New("customer details are required")
	}

	if checkout.ShippingMethodID == 0 {
		return nil, errors.New("shipping method is required")
	}

	// Convert checkout to order
	order := checkout.ToOrder()

	// Create order in repository
	err = uc.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	// Mark checkout as completed
	checkout.MarkAsCompleted(order.ID)
	err = uc.checkoutRepo.Update(checkout)
	// TODO: Handle error but do not return it, as we want to proceed with order creation even if updating the checkout fails
	if err != nil {
		fmt.Printf("Failed to update checkout after order creation: %v\n", err)
	}

	// Increment discount usage if a discount was applied
	if checkout.AppliedDiscount != nil {
		discount, err := uc.discountRepo.GetByID(checkout.AppliedDiscount.DiscountID)
		if err == nil {
			discount.IncrementUsage()
			uc.discountRepo.Update(discount)
		}
	}

	return order, nil
}

// ExtendCheckoutExpiry extends the expiry time of a checkout
func (uc *CheckoutUseCase) ExtendCheckoutExpiry(checkoutID uint, duration time.Duration) (*entity.Checkout, error) {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Extend expiry
	checkout.ExtendExpiry(duration)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// GetCheckoutByID retrieves a checkout by ID
func (uc *CheckoutUseCase) GetCheckoutByID(checkoutID uint) (*entity.Checkout, error) {
	return uc.checkoutRepo.GetByID(checkoutID)
}

// AbandonCheckout marks a checkout as abandoned
func (uc *CheckoutUseCase) AbandonCheckout(checkoutID uint) error {
	// Get checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return err
	}

	// Mark as abandoned
	checkout.MarkAsAbandoned()

	// Update checkout in repository
	return uc.checkoutRepo.Update(checkout)
}

// GetCheckoutsByStatus retrieves checkouts by status with pagination
func (uc *CheckoutUseCase) GetCheckoutsByStatus(status entity.CheckoutStatus, offset, limit int) ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetCheckoutsByStatus(status, offset, limit)
}

// GetAllCheckouts retrieves all checkouts with pagination
func (uc *CheckoutUseCase) GetAllCheckouts(offset, limit int) ([]*entity.Checkout, error) {
	// If no specific status is requested, get checkouts regardless of status
	return uc.checkoutRepo.GetCheckoutsByStatus("", offset, limit)
}

// DeleteCheckout deletes a checkout by ID
func (uc *CheckoutUseCase) DeleteCheckout(checkoutID uint) error {
	return uc.checkoutRepo.Delete(checkoutID)
}

// GetExpiredCheckouts retrieves all expired checkouts
func (uc *CheckoutUseCase) GetExpiredCheckouts() ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetExpiredCheckouts()
}

// GetAbandonedCheckouts retrieves all abandoned checkouts
func (uc *CheckoutUseCase) GetAbandonedCheckouts(offset, limit int) ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetCheckoutsByStatus(entity.CheckoutStatusAbandoned, offset, limit)
}

// GetCheckoutsByUserID retrieves all checkouts for a user with pagination
func (uc *CheckoutUseCase) GetCheckoutsByUserID(userID uint, offset, limit int) ([]*entity.Checkout, error) {
	return uc.checkoutRepo.GetCompletedCheckoutsByUserID(userID, offset, limit)
}

// GetCheckoutBySessionID retrieves a checkout by session ID
func (uc *CheckoutUseCase) GetCheckoutBySessionID(sessionID string) (*entity.Checkout, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}
	return uc.checkoutRepo.GetBySessionID(sessionID)
}

// UpdateCheckout updates a checkout in the repository
func (uc *CheckoutUseCase) UpdateCheckout(checkout *entity.Checkout) (*entity.Checkout, error) {
	if checkout == nil {
		return nil, errors.New("checkout cannot be nil")
	}

	// Make sure the checkout is active
	if checkout.Status != entity.CheckoutStatusActive {
		return nil, errors.New("cannot update a non-active checkout")
	}

	// Update timestamps
	now := time.Now()
	checkout.UpdatedAt = now
	checkout.LastActivityAt = now

	// Save to repository
	err := uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// GetOrCreateCheckoutBySessionID retrieves or creates a checkout using a session ID
func (uc *CheckoutUseCase) GetOrCreateCheckoutBySessionID(sessionID string) (*entity.Checkout, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}

	// Try to get an existing active checkout
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err == nil {
		// If found, return it
		return checkout, nil
	}

	// If not found, create a new one
	checkout, err = entity.NewCheckout(sessionID)
	if err != nil {
		return nil, err
	}

	// Set default currency
	defaultCurrency, err := uc.currencyRepo.GetDefault()
	if err == nil && defaultCurrency != nil {
		checkout.Currency = defaultCurrency.Code
	}

	// Save to repository
	err = uc.checkoutRepo.Create(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// UpdateOrder updates an order in the repository
func (uc *CheckoutUseCase) UpdateOrder(order *entity.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	return uc.orderRepo.Update(order)
}

// AddItemToCheckout adds an item to a checkout by ID using SKU
func (uc *CheckoutUseCase) AddItemToCheckout(checkoutID uint, input CheckoutInput) (*entity.Checkout, error) {
	// Get the checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Check if checkout is active
	if checkout.Status != entity.CheckoutStatusActive {
		return nil, errors.New("cannot modify a non-active checkout")
	}

	// Validate SKU is provided
	if input.SKU == "" {
		return nil, errors.New("SKU is required")
	}

	// Find the product variant by SKU (all products now have variants)
	variant, err := uc.productVariantRepo.GetBySKU(input.SKU)
	if err != nil {
		return nil, fmt.Errorf("product variant not found with SKU '%s'", input.SKU)
	}

	// Get the parent product
	product, err := uc.productRepo.GetByID(variant.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product for variant: %w", err)
	}

	// Check if product is active
	if !product.Active {
		return nil, errors.New("product is not available")
	}

	// Extract variant name from attributes
	variantName := ""
	for _, attr := range variant.Attributes {
		if variantName == "" {
			variantName = attr.Value
		} else {
			variantName += " / " + attr.Value
		}
	}

	// Populate input with variant details
	input.ProductID = variant.ProductID
	input.VariantID = variant.ID
	input.ProductName = product.Name
	input.VariantName = variantName
	input.Price = variant.Price
	input.Weight = product.Weight

	// Add the item to the checkout
	err = checkout.AddItem(input.ProductID, input.VariantID, input.Quantity, input.Price, input.Weight, input.ProductName, input.VariantName, input.SKU)
	if err != nil {
		return nil, err
	}

	// Save the updated checkout
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// UpdateCheckoutItemBySKU updates an item in a checkout by SKU
func (uc *CheckoutUseCase) UpdateCheckoutItemBySKU(checkoutID uint, input UpdateCheckoutItemInput) (*entity.Checkout, error) {
	// Get the checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Check if checkout is active
	if checkout.Status != entity.CheckoutStatusActive {
		return nil, errors.New("cannot modify a non-active checkout")
	}

	// Validate SKU is provided
	if input.SKU == "" {
		return nil, errors.New("SKU is required")
	}

	// Validate quantity
	if input.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	// Find the product variant by SKU (all products now have variants)
	variant, err := uc.productVariantRepo.GetBySKU(input.SKU)
	if err != nil {
		return nil, fmt.Errorf("product variant not found with SKU '%s'", input.SKU)
	}

	// Get the parent product
	product, err := uc.productRepo.GetByID(variant.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product for variant: %w", err)
	}

	// Check if product is active
	if !product.Active {
		return nil, errors.New("product is not available")
	}

	productID := variant.ProductID
	variantID := variant.ID

	// Update the item in the checkout
	err = checkout.UpdateItem(productID, variantID, input.Quantity)
	if err != nil {
		return nil, err
	}

	// Save the updated checkout
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// RemoveItemBySKU removes an item from a checkout by SKU
func (uc *CheckoutUseCase) RemoveItemBySKU(checkoutID uint, input RemoveItemInput) (*entity.Checkout, error) {
	// Get the checkout
	checkout, err := uc.checkoutRepo.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	// Check if checkout is active
	if checkout.Status != entity.CheckoutStatusActive {
		return nil, errors.New("cannot modify a non-active checkout")
	}

	// Validate SKU is provided
	if input.SKU == "" {
		return nil, errors.New("SKU is required")
	}

	// Find the product variant by SKU (all products now have variants)
	variant, err := uc.productVariantRepo.GetBySKU(input.SKU)
	if err != nil {
		return nil, fmt.Errorf("product variant not found with SKU '%s'", input.SKU)
	}

	// Get the parent product
	product, err := uc.productRepo.GetByID(variant.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product for variant: %w", err)
	}

	// Check if product is active
	if !product.Active {
		return nil, errors.New("product is not available")
	}

	productID := variant.ProductID
	variantID := variant.ID

	// Remove the item from the checkout
	err = checkout.RemoveItem(productID, variantID)
	if err != nil {
		return nil, err
	}

	// Save the updated checkout
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// ChangeCurrency changes the currency of a checkout and converts all prices
func (uc *CheckoutUseCase) ChangeCurrency(checkout *entity.Checkout, newCurrencyCode string) (*entity.Checkout, error) {
	// Validate that the new currency exists and is enabled
	toCurrency, err := uc.currencyRepo.GetByCode(newCurrencyCode)
	if err != nil {
		return nil, fmt.Errorf("currency %s not found: %w", newCurrencyCode, err)
	}

	if !toCurrency.IsEnabled {
		return nil, fmt.Errorf("currency %s is not enabled", newCurrencyCode)
	}

	// Get the current currency
	fromCurrency, err := uc.currencyRepo.GetByCode(checkout.Currency)
	if err != nil {
		return nil, fmt.Errorf("current currency %s not found: %w", checkout.Currency, err)
	}

	// Convert the checkout currency and all prices
	checkout.SetCurrency(newCurrencyCode, fromCurrency, toCurrency)

	// Update checkout in repository
	err = uc.checkoutRepo.Update(checkout)
	if err != nil {
		return nil, fmt.Errorf("failed to update checkout: %w", err)
	}

	return checkout, nil
}

// ChangeCurrencyBySessionID changes the currency of a checkout by session ID
func (uc *CheckoutUseCase) ChangeCurrencyBySessionID(sessionID string, newCurrencyCode string) (*entity.Checkout, error) {
	// Get checkout by session ID
	checkout, err := uc.checkoutRepo.GetBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("checkout not found for session %s: %w", sessionID, err)
	}

	return uc.ChangeCurrency(checkout, newCurrencyCode)
}

// ChangeCurrencyByUserID changes the currency of a checkout by user ID
func (uc *CheckoutUseCase) ChangeCurrencyByUserID(userID uint, newCurrencyCode string) (*entity.Checkout, error) {
	// Get checkout by user ID
	checkout, err := uc.checkoutRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("checkout not found for user %d: %w", userID, err)
	}

	return uc.ChangeCurrency(checkout, newCurrencyCode)
}
