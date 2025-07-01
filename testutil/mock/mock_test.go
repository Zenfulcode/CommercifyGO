package mock_test

import (
	"fmt"
	"testing"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/testutil/mock"
)

func TestPaymentTransactionRepository(t *testing.T) {
	repo := mock.NewPaymentTransactionRepository()

	// Test Create
	txn, err := entity.NewPaymentTransaction(
		1,                                  // orderID
		"test-123",                         // transactionID
		entity.TransactionTypeAuthorize,    // type
		entity.TransactionStatusSuccessful, // status
		10000,                              // amount (100.00)
		"USD",                              // currency
		"stripe",                           // provider
	)
	if err != nil {
		t.Fatalf("Failed to create payment transaction: %v", err)
	}

	err = repo.Create(txn)
	if err != nil {
		t.Fatalf("Failed to create payment transaction in repo: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(txn.ID)
	if err != nil {
		t.Fatalf("Failed to get payment transaction by ID: %v", err)
	}

	if retrieved.TransactionID != "test-123" {
		t.Errorf("Expected transaction ID 'test-123', got '%s'", retrieved.TransactionID)
	}

	// Test GetByTransactionID
	retrieved2, err := repo.GetByTransactionID("test-123")
	if err != nil {
		t.Fatalf("Failed to get payment transaction by transaction ID: %v", err)
	}

	if retrieved2.OrderID != 1 {
		t.Errorf("Expected order ID 1, got %d", retrieved2.OrderID)
	}

	// Test GetByOrderID
	orderTxns, err := repo.GetByOrderID(1)
	if err != nil {
		t.Fatalf("Failed to get payment transactions by order ID: %v", err)
	}

	if len(orderTxns) != 1 {
		t.Errorf("Expected 1 transaction for order, got %d", len(orderTxns))
	}

	// Test CountSuccessfulByOrderIDAndType
	count, err := repo.CountSuccessfulByOrderIDAndType(1, entity.TransactionTypeAuthorize)
	if err != nil {
		t.Fatalf("Failed to count successful transactions: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 successful authorization, got %d", count)
	}

	// Test SumAmountByOrderIDAndType
	sum, err := repo.SumAmountByOrderIDAndType(1, entity.TransactionTypeAuthorize)
	if err != nil {
		t.Fatalf("Failed to sum transaction amounts: %v", err)
	}

	if sum != 10000 {
		t.Errorf("Expected sum of 10000, got %d", sum)
	}
}

func TestPaymentProviderRepository(t *testing.T) {
	repo := mock.NewPaymentProviderRepository()

	// Create a test provider
	provider := &entity.PaymentProvider{
		Type:                common.PaymentProviderStripe,
		Name:                "Stripe",
		Description:         "Stripe payment processor",
		Methods:             []common.PaymentMethod{common.PaymentMethodCreditCard},
		Enabled:             true,
		SupportedCurrencies: []string{"USD", "EUR"},
	}

	// Test Create
	err := repo.Create(provider)
	if err != nil {
		t.Fatalf("Failed to create payment provider: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(provider.ID)
	if err != nil {
		t.Fatalf("Failed to get payment provider by ID: %v", err)
	}

	if retrieved.Name != "Stripe" {
		t.Errorf("Expected provider name 'Stripe', got '%s'", retrieved.Name)
	}

	// Test GetByType
	retrieved2, err := repo.GetByType(common.PaymentProviderStripe)
	if err != nil {
		t.Fatalf("Failed to get payment provider by type: %v", err)
	}

	if retrieved2.Description != "Stripe payment processor" {
		t.Errorf("Expected description 'Stripe payment processor', got '%s'", retrieved2.Description)
	}

	// Test GetEnabled
	enabled, err := repo.GetEnabled()
	if err != nil {
		t.Fatalf("Failed to get enabled payment providers: %v", err)
	}

	if len(enabled) != 1 {
		t.Errorf("Expected 1 enabled provider, got %d", len(enabled))
	}

	// Test GetEnabledByMethod
	methodProviders, err := repo.GetEnabledByMethod(common.PaymentMethodCreditCard)
	if err != nil {
		t.Fatalf("Failed to get providers by method: %v", err)
	}

	if len(methodProviders) != 1 {
		t.Errorf("Expected 1 provider for credit card method, got %d", len(methodProviders))
	}

	// Test GetEnabledByCurrency
	currencyProviders, err := repo.GetEnabledByCurrency("USD")
	if err != nil {
		t.Fatalf("Failed to get providers by currency: %v", err)
	}

	if len(currencyProviders) != 1 {
		t.Errorf("Expected 1 provider for USD currency, got %d", len(currencyProviders))
	}
}

func TestOrderRepository(t *testing.T) {
	repo := mock.NewOrderRepository()

	// Create a test order
	order := &entity.Order{
		OrderNumber:       "ORD-123",
		Currency:          "USD",
		UserID:            &[]uint{1}[0],
		TotalAmount:       10000,
		FinalAmount:       10000,
		Status:            entity.OrderStatusPending,
		PaymentStatus:     entity.PaymentStatusPending,
		CheckoutSessionID: "cs_test_123",
		PaymentID:         "pi_test_123",
	}

	// Test Create
	err := repo.Create(order)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(order.ID)
	if err != nil {
		t.Fatalf("Failed to get order by ID: %v", err)
	}

	if retrieved.OrderNumber != "ORD-123" {
		t.Errorf("Expected order number 'ORD-123', got '%s'", retrieved.OrderNumber)
	}

	// Test GetByCheckoutSessionID
	retrieved2, err := repo.GetByCheckoutSessionID("cs_test_123")
	if err != nil {
		t.Fatalf("Failed to get order by checkout session ID: %v", err)
	}

	if retrieved2.Currency != "USD" {
		t.Errorf("Expected currency 'USD', got '%s'", retrieved2.Currency)
	}

	// Test GetByPaymentID
	retrieved3, err := repo.GetByPaymentID("pi_test_123")
	if err != nil {
		t.Fatalf("Failed to get order by payment ID: %v", err)
	}

	if retrieved3.TotalAmount != 10000 {
		t.Errorf("Expected total amount 10000, got %d", retrieved3.TotalAmount)
	}

	// Test GetByUser
	userOrders, err := repo.GetByUser(1, 0, 10)
	if err != nil {
		t.Fatalf("Failed to get orders by user: %v", err)
	}

	if len(userOrders) != 1 {
		t.Errorf("Expected 1 order for user, got %d", len(userOrders))
	}

	// Test ListByStatus
	statusOrders, err := repo.ListByStatus(entity.OrderStatusPending, 0, 10)
	if err != nil {
		t.Fatalf("Failed to get orders by status: %v", err)
	}

	if len(statusOrders) != 1 {
		t.Errorf("Expected 1 pending order, got %d", len(statusOrders))
	}
}

func TestLogger(t *testing.T) {
	logger := mock.NewLogger()

	// Test logging different levels
	logger.Debug("Debug message with %s", "args")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Test getting logs
	logs := logger.(*mock.Logger).GetLogs()
	if len(logs) != 4 {
		t.Errorf("Expected 4 log entries, got %d", len(logs))
	}

	// Test getting logs by level
	errorLogs := logger.(*mock.Logger).GetLogsByLevel("ERROR")
	if len(errorLogs) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(errorLogs))
	}

	// Test checking for specific messages
	if !logger.(*mock.Logger).HasLogWithMessage("Info message") {
		t.Error("Expected to find 'Info message' in logs")
	}

	// Test checking for specific levels
	if !logger.(*mock.Logger).HasLogWithLevel("WARN") {
		t.Error("Expected to find WARN level in logs")
	}

	// Test log count
	if logger.(*mock.Logger).LogCount() != 4 {
		t.Errorf("Expected log count of 4, got %d", logger.(*mock.Logger).LogCount())
	}

	// Test clearing logs
	logger.(*mock.Logger).Clear()
	if logger.(*mock.Logger).LogCount() != 0 {
		t.Errorf("Expected log count of 0 after clear, got %d", logger.(*mock.Logger).LogCount())
	}
}

func TestProductRepository(t *testing.T) {
	repo := mock.NewProductRepository()

	// Create a test variant
	variant := &entity.ProductVariant{
		SKU:        "TEST-SKU-001",
		Stock:      10,
		Price:      9999, // $99.99
		IsDefault:  true,
		Weight:     1.5,
		Attributes: entity.VariantAttributes{"color": "red", "size": "large"},
		Images:     common.StringSlice{"variant1.jpg"},
	}

	// Create a test product
	product := &entity.Product{
		Name:        "Test Product",
		Description: "A test product description",
		Currency:    "USD",
		CategoryID:  1,
		Images:      common.StringSlice{"product1.jpg", "product2.jpg"},
		Active:      true,
		Variants:    []*entity.ProductVariant{variant},
	}

	// Test Create
	err := repo.Create(product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(product.ID)
	if err != nil {
		t.Fatalf("Failed to get product by ID: %v", err)
	}

	if retrieved.Name != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%s'", retrieved.Name)
	}

	if len(retrieved.Variants) != 1 {
		t.Errorf("Expected 1 variant, got %d", len(retrieved.Variants))
	}

	// Test GetByIDAndCurrency
	retrieved2, err := repo.GetByIDAndCurrency(product.ID, "USD")
	if err != nil {
		t.Fatalf("Failed to get product by ID and currency: %v", err)
	}

	if retrieved2.Currency != "USD" {
		t.Errorf("Expected currency 'USD', got '%s'", retrieved2.Currency)
	}

	// Test wrong currency
	_, err = repo.GetByIDAndCurrency(product.ID, "EUR")
	if err == nil {
		t.Error("Expected error for wrong currency, got none")
	}

	// Test GetBySKU
	retrieved3, err := repo.GetBySKU("TEST-SKU-001")
	if err != nil {
		t.Fatalf("Failed to get product by SKU: %v", err)
	}

	if retrieved3.Name != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%s'", retrieved3.Name)
	}

	// Test List with no filters
	products, err := repo.List("", "", 0, 0, 10, 0, 0, true)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if len(products) != 1 {
		t.Errorf("Expected 1 product in list, got %d", len(products))
	}

	// Test List with name query
	products2, err := repo.List("Test", "", 0, 0, 10, 0, 0, true)
	if err != nil {
		t.Fatalf("Failed to list products with query: %v", err)
	}

	if len(products2) != 1 {
		t.Errorf("Expected 1 product matching 'Test', got %d", len(products2))
	}

	// Test List with no matching query
	products3, err := repo.List("NonExistent", "", 0, 0, 10, 0, 0, true)
	if err != nil {
		t.Fatalf("Failed to list products with non-matching query: %v", err)
	}

	if len(products3) != 0 {
		t.Errorf("Expected 0 products matching 'NonExistent', got %d", len(products3))
	}

	// Test List with price range
	products4, err := repo.List("", "", 0, 0, 10, 5000, 15000, true)
	if err != nil {
		t.Fatalf("Failed to list products with price range: %v", err)
	}

	if len(products4) != 1 {
		t.Errorf("Expected 1 product in price range, got %d", len(products4))
	}

	// Test List with category filter
	products5, err := repo.List("", "", 1, 0, 10, 0, 0, true)
	if err != nil {
		t.Fatalf("Failed to list products with category filter: %v", err)
	}

	if len(products5) != 1 {
		t.Errorf("Expected 1 product in category 1, got %d", len(products5))
	}

	// Test Count
	count, err := repo.Count("", "", 0, 0, 0, true)
	if err != nil {
		t.Fatalf("Failed to count products: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count of 1, got %d", count)
	}

	// Test Update
	product.Description = "Updated description"
	err = repo.Update(product)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	updated, err := repo.GetByID(product.ID)
	if err != nil {
		t.Fatalf("Failed to get updated product: %v", err)
	}

	if updated.Description != "Updated description" {
		t.Errorf("Expected updated description, got '%s'", updated.Description)
	}

	// Test Delete
	err = repo.Delete(product.ID)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(product.ID)
	if err == nil {
		t.Error("Expected error after deletion, got none")
	}

	// Verify SKU mapping was also removed
	_, err = repo.GetBySKU("TEST-SKU-001")
	if err == nil {
		t.Error("Expected error for deleted SKU, got none")
	}
}

func TestProductRepositoryEdgeCases(t *testing.T) {
	repo := mock.NewProductRepository()

	// Test duplicate SKU creation
	variant1 := &entity.ProductVariant{
		SKU:       "DUPLICATE-SKU",
		Stock:     5,
		Price:     1000,
		IsDefault: true,
	}

	variant2 := &entity.ProductVariant{
		SKU:       "DUPLICATE-SKU",
		Stock:     3,
		Price:     2000,
		IsDefault: false,
	}

	product1 := &entity.Product{
		Name:       "Product 1",
		Currency:   "USD",
		CategoryID: 1,
		Active:     true,
		Variants:   []*entity.ProductVariant{variant1},
	}

	product2 := &entity.Product{
		Name:       "Product 2",
		Currency:   "USD",
		CategoryID: 1,
		Active:     true,
		Variants:   []*entity.ProductVariant{variant2},
	}

	// First product should create successfully
	err := repo.Create(product1)
	if err != nil {
		t.Fatalf("Failed to create first product: %v", err)
	}

	// Second product should fail due to duplicate SKU
	err = repo.Create(product2)
	if err == nil {
		t.Error("Expected error for duplicate SKU, got none")
	}

	// Test error scenarios
	repo.(*mock.ProductRepository).SetGetByIDError(fmt.Errorf("database error"))
	_, err = repo.GetByID(1)
	if err == nil {
		t.Error("Expected configured error, got none")
	}

	// Reset errors
	repo.(*mock.ProductRepository).Reset()

	// Verify reset worked
	products := repo.(*mock.ProductRepository).GetAllProducts()
	if len(products) != 0 {
		t.Errorf("Expected 0 products after reset, got %d", len(products))
	}
}
