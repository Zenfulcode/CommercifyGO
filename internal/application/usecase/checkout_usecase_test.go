package usecase

import (
	"strings"
	"testing"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/testutil/mock"
)

func TestCheckout_Currency_Validation(t *testing.T) {
	t.Run("NewCheckout should reject empty currency", func(t *testing.T) {
		_, err := entity.NewCheckout("test-session", "")

		if err == nil {
			t.Error("Expected error when creating checkout with empty currency")
		}

		expectedMsg := "currency cannot be empty"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("NewCheckout should accept valid currency", func(t *testing.T) {
		checkout, err := entity.NewCheckout("test-session", "EUR")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if checkout.Currency != "EUR" {
			t.Errorf("Expected currency EUR, got %s", checkout.Currency)
		}
	})
}

func TestCheckout_Currency_DefaultHandling(t *testing.T) {
	// Create mock currency repository
	currencyRepo := mock.NewMockCurrencyRepository()

	// Setup test currencies
	usd, _ := entity.NewCurrency("USD", "US Dollar", "$", 1.0, true, true)
	eur, _ := entity.NewCurrency("EUR", "Euro", "€", 0.92, true, false)
	dkk, _ := entity.NewCurrency("DKK", "Danish Krone", "kr", 6.54, true, false)

	currencyRepo.Create(usd)
	currencyRepo.Create(eur)
	currencyRepo.Create(dkk)

	t.Run("Should use USD as default currency initially", func(t *testing.T) {
		defaultCurrency, err := currencyRepo.GetDefault()
		if err != nil {
			t.Fatalf("Expected no error getting default currency, got %v", err)
		}

		if defaultCurrency.Code != "USD" {
			t.Errorf("Expected default currency to be USD, got %s", defaultCurrency.Code)
		}

		// Create checkout with this default currency
		checkout, err := entity.NewCheckout("test-session-1", defaultCurrency.Code)
		if err != nil {
			t.Fatalf("Expected no error creating checkout, got %v", err)
		}

		if checkout.Currency != "USD" {
			t.Errorf("Expected checkout currency to be USD, got %s", checkout.Currency)
		}
	})

	t.Run("Should use EUR as default after changing default", func(t *testing.T) {
		// Change default currency to EUR
		err := currencyRepo.SetDefault("EUR")
		if err != nil {
			t.Fatalf("Failed to set EUR as default: %v", err)
		}

		defaultCurrency, err := currencyRepo.GetDefault()
		if err != nil {
			t.Fatalf("Expected no error getting default currency, got %v", err)
		}

		if defaultCurrency.Code != "EUR" {
			t.Errorf("Expected default currency to be EUR, got %s", defaultCurrency.Code)
		}

		// Create checkout with this default currency
		checkout, err := entity.NewCheckout("test-session-2", defaultCurrency.Code)
		if err != nil {
			t.Fatalf("Expected no error creating checkout, got %v", err)
		}

		if checkout.Currency != "EUR" {
			t.Errorf("Expected checkout currency to be EUR, got %s", checkout.Currency)
		}
	})

	t.Run("Should use DKK as default after changing default", func(t *testing.T) {
		// Change default currency to DKK
		err := currencyRepo.SetDefault("DKK")
		if err != nil {
			t.Fatalf("Failed to set DKK as default: %v", err)
		}

		defaultCurrency, err := currencyRepo.GetDefault()
		if err != nil {
			t.Fatalf("Expected no error getting default currency, got %v", err)
		}

		if defaultCurrency.Code != "DKK" {
			t.Errorf("Expected default currency to be DKK, got %s", defaultCurrency.Code)
		}

		// Create checkout with this default currency
		checkout, err := entity.NewCheckout("test-session-3", defaultCurrency.Code)
		if err != nil {
			t.Fatalf("Expected no error creating checkout, got %v", err)
		}

		if checkout.Currency != "DKK" {
			t.Errorf("Expected checkout currency to be DKK, got %s", checkout.Currency)
		}
	})
}

func TestCheckout_AddItem_CurrencyConversion(t *testing.T) {
	// Setup mock repositories
	checkoutRepo := mock.NewMockCheckoutRepository()
	currencyRepo := mock.NewMockCurrencyRepository()
	productRepo := mock.NewMockProductRepository()
	productVariantRepo := mock.NewMockProductVariantRepository()
	discountRepo := mock.NewMockDiscountRepository()
	orderRepo := mock.NewMockOrderRepository(false)
	paymentTransactionRepo := mock.NewMockPaymentTransactionRepository()

	// Setup test currencies
	usd, _ := entity.NewCurrency("USD", "US Dollar", "$", 1.0, true, true)
	eur, _ := entity.NewCurrency("EUR", "Euro", "€", 0.85, true, false)         // 1 USD = 0.85 EUR
	dkk, _ := entity.NewCurrency("DKK", "Danish Krone", "kr", 6.8, true, false) // 1 USD = 6.8 DKK

	currencyRepo.Create(usd)
	currencyRepo.Create(eur)
	currencyRepo.Create(dkk)

	// Create checkout usecase with minimal dependencies
	usecase := NewCheckoutUseCase(
		checkoutRepo,
		productRepo,
		productVariantRepo,
		nil, // shippingMethodRepo - not needed for currency tests
		nil, // shippingRateRepo - not needed for currency tests
		discountRepo,
		orderRepo,
		currencyRepo,
		paymentTransactionRepo,
		nil, // paymentSvc - not needed for these tests
		nil, // shippingUsecase - not needed for currency tests
	)

	t.Run("Should convert variant price to checkout currency", func(t *testing.T) {
		// Create a checkout in EUR
		checkoutEUR, _ := entity.NewCheckout("test-session", "EUR")
		checkoutEUR.ID = 1
		checkoutRepo.Create(checkoutEUR)

		// Create a product
		product, _ := entity.NewProduct("Test Product", "A test product", "USD", 1, nil)
		product.ID = 1
		product.Active = true
		productRepo.Create(product)

		// Create a variant with USD price (100 USD = 10000 cents)
		variant, _ := entity.NewProductVariant(1, "TEST-SKU", 100.0, "USD", 10, nil, nil, true)
		variant.ID = 1
		productVariantRepo.Create(variant)

		// Add item to checkout
		input := CheckoutInput{
			SKU:      "TEST-SKU",
			Quantity: 2,
		}

		updatedCheckout, err := usecase.AddItemToCheckout(1, input)
		if err != nil {
			t.Fatalf("Expected no error adding item to checkout, got %v", err)
		}

		// Check that the item was added with converted price
		if len(updatedCheckout.Items) != 1 {
			t.Fatalf("Expected 1 item in checkout, got %d", len(updatedCheckout.Items))
		}

		item := updatedCheckout.Items[0]

		// Expected price: 100 USD * 0.85 = 85 EUR = 8500 cents
		expectedPriceInCents := int64(8500)
		if item.Price != expectedPriceInCents {
			t.Errorf("Expected item price to be %d EUR cents (converted from USD), got %d", expectedPriceInCents, item.Price)
		}

		// Check that checkout currency is still EUR
		if updatedCheckout.Currency != "EUR" {
			t.Errorf("Expected checkout currency to remain EUR, got %s", updatedCheckout.Currency)
		}
	})

	t.Run("Should use variant price directly if currencies match", func(t *testing.T) {
		// Create a checkout in USD
		checkoutUSD, _ := entity.NewCheckout("test-session-2", "USD")
		checkoutUSD.ID = 2
		checkoutRepo.Create(checkoutUSD)

		// Create a product priced in USD
		product2, _ := entity.NewProduct("Test Product 2", "Another test product", "USD", 1, nil)
		product2.ID = 2
		product2.Active = true
		productRepo.Create(product2)

		// Create a variant with USD price (50 USD = 5000 cents)
		variant2, _ := entity.NewProductVariant(2, "TEST-SKU-2", 50.0, "USD", 5, nil, nil, true)
		variant2.ID = 2
		productVariantRepo.Create(variant2)

		// Add item to checkout
		input := CheckoutInput{
			SKU:      "TEST-SKU-2",
			Quantity: 1,
		}

		updatedCheckout, err := usecase.AddItemToCheckout(2, input)
		if err != nil {
			t.Fatalf("Expected no error adding item to checkout, got %v", err)
		}

		// Check that the item was added with original price (no conversion needed)
		if len(updatedCheckout.Items) != 1 {
			t.Fatalf("Expected 1 item in checkout, got %d", len(updatedCheckout.Items))
		}

		item := updatedCheckout.Items[0]

		// Expected price: 50 USD = 5000 cents (no conversion)
		expectedPriceInCents := int64(5000)
		if item.Price != expectedPriceInCents {
			t.Errorf("Expected item price to be %d USD cents (no conversion needed), got %d", expectedPriceInCents, item.Price)
		}
	})

	t.Run("Should use variant's multi-currency price if available", func(t *testing.T) {
		// Create a checkout in DKK
		checkoutDKK, _ := entity.NewCheckout("test-session-3", "DKK")
		checkoutDKK.ID = 3
		checkoutRepo.Create(checkoutDKK)

		// Create a product
		product3, _ := entity.NewProduct("Test Product 3", "Product with multiple currency prices", "USD", 1, nil)
		product3.ID = 3
		product3.Active = true
		productRepo.Create(product3)

		// Create a variant with USD base price but also DKK price
		variant3, _ := entity.NewProductVariant(3, "TEST-SKU-3", 75.0, "USD", 8, nil, nil, true)
		variant3.ID = 3

		// Add a specific DKK price to the variant (500 DKK = 50000 cents)
		variant3.Prices = []entity.ProductVariantPrice{
			{
				VariantID:    3,
				CurrencyCode: "DKK",
				Price:        50000, // 500 DKK in cents
			},
		}
		productVariantRepo.Create(variant3)

		// Verify the variant has the DKK price set
		dkkPrice, hasDkkPrice := variant3.GetPriceInCurrency("DKK")
		if !hasDkkPrice {
			t.Fatalf("Expected variant to have DKK price, but GetPriceInCurrency returned false")
		}
		if dkkPrice != 50000 {
			t.Fatalf("Expected DKK price to be 50000 cents, got %d", dkkPrice)
		}

		// Add item to checkout
		input := CheckoutInput{
			SKU:      "TEST-SKU-3",
			Quantity: 1,
		}

		updatedCheckout, err := usecase.AddItemToCheckout(3, input)
		if err != nil {
			t.Fatalf("Expected no error adding item to checkout, got %v", err)
		}

		// Check that the item was added with the DKK-specific price
		if len(updatedCheckout.Items) != 1 {
			t.Fatalf("Expected 1 item in checkout, got %d", len(updatedCheckout.Items))
		}

		item := updatedCheckout.Items[0]

		// Expected price: 500 DKK = 50000 cents (from variant's specific DKK price)
		expectedPriceInCents := int64(50000)
		if item.Price != expectedPriceInCents {
			t.Errorf("Expected item price to be %d DKK cents (from variant's specific price), got %d", expectedPriceInCents, item.Price)
		}
	})
}

func TestCheckout_Currency_ParameterHandling(t *testing.T) {
	// Setup mock repositories
	checkoutRepo := mock.NewMockCheckoutRepository()
	currencyRepo := mock.NewMockCurrencyRepository()
	productRepo := mock.NewMockProductRepository()
	productVariantRepo := mock.NewMockProductVariantRepository()
	discountRepo := mock.NewMockDiscountRepository()
	orderRepo := mock.NewMockOrderRepository(false)
	paymentTransactionRepo := mock.NewMockPaymentTransactionRepository()

	// Setup test currencies
	usd, _ := entity.NewCurrency("USD", "US Dollar", "$", 1.0, true, true)
	eur, _ := entity.NewCurrency("EUR", "Euro", "€", 0.85, true, false)
	dkk, _ := entity.NewCurrency("DKK", "Danish Krone", "kr", 6.8, true, false)

	currencyRepo.Create(usd)
	currencyRepo.Create(eur)
	currencyRepo.Create(dkk)

	// Create checkout usecase
	usecase := NewCheckoutUseCase(
		checkoutRepo,
		productRepo,
		productVariantRepo,
		nil, // shippingMethodRepo - not needed for currency tests
		nil, // shippingRateRepo - not needed for currency tests
		discountRepo,
		orderRepo,
		currencyRepo,
		paymentTransactionRepo,
		nil, // paymentSvc - not needed for these tests
		nil, // shippingUsecase - not needed for currency tests
	)

	t.Run("Should create checkout with specified currency", func(t *testing.T) {
		checkout, err := usecase.GetOrCreateCheckoutBySessionIDWithCurrency("test-session-1", "EUR")
		if err != nil {
			t.Fatalf("Expected no error creating checkout with EUR, got %v", err)
		}

		if checkout.Currency != "EUR" {
			t.Errorf("Expected checkout currency to be EUR, got %s", checkout.Currency)
		}
	})

	t.Run("Should change currency of existing checkout", func(t *testing.T) {
		// First create a checkout with USD
		checkout1, _ := usecase.GetOrCreateCheckoutBySessionIDWithCurrency("test-session-2", "USD")
		if checkout1.Currency != "USD" {
			t.Errorf("Expected initial checkout currency to be USD, got %s", checkout1.Currency)
		}

		// Then request the same session with EUR - should convert
		checkout2, err := usecase.GetOrCreateCheckoutBySessionIDWithCurrency("test-session-2", "EUR")
		if err != nil {
			t.Fatalf("Expected no error changing checkout currency to EUR, got %v", err)
		}

		if checkout2.Currency != "EUR" {
			t.Errorf("Expected checkout currency to be changed to EUR, got %s", checkout2.Currency)
		}

		// Should be the same checkout object with updated currency
		if checkout2.ID != checkout1.ID {
			t.Errorf("Expected same checkout ID %d, got %d", checkout1.ID, checkout2.ID)
		}
	})

	t.Run("Should use default currency when no currency specified", func(t *testing.T) {
		checkout, err := usecase.GetOrCreateCheckoutBySessionIDWithCurrency("test-session-3", "")
		if err != nil {
			t.Fatalf("Expected no error creating checkout with default currency, got %v", err)
		}

		// Should use the default currency (USD)
		if checkout.Currency != "USD" {
			t.Errorf("Expected checkout currency to be USD (default), got %s", checkout.Currency)
		}
	})

	t.Run("Should return error for invalid currency", func(t *testing.T) {
		_, err := usecase.GetOrCreateCheckoutBySessionIDWithCurrency("test-session-4", "INVALID")
		if err == nil {
			t.Error("Expected error for invalid currency, got nil")
		}

		expectedErrSubstring := "invalid currency INVALID"
		if !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected error to contain '%s', got '%s'", expectedErrSubstring, err.Error())
		}
	})

	t.Run("Should return error for disabled currency", func(t *testing.T) {
		// Create a disabled currency
		gbp, _ := entity.NewCurrency("GBP", "British Pound", "£", 0.8, false, false) // disabled
		currencyRepo.Create(gbp)

		_, err := usecase.GetOrCreateCheckoutBySessionIDWithCurrency("test-session-5", "GBP")
		if err == nil {
			t.Error("Expected error for disabled currency, got nil")
		}

		expectedErrSubstring := "currency GBP is not enabled"
		if !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected error to contain '%s', got '%s'", expectedErrSubstring, err.Error())
		}
	})

	t.Run("Should maintain backward compatibility", func(t *testing.T) {
		// Test the original method without currency parameter
		checkout, err := usecase.GetOrCreateCheckoutBySessionID("test-session-6")
		if err != nil {
			t.Fatalf("Expected no error with original method, got %v", err)
		}

		// Should use default currency
		if checkout.Currency != "USD" {
			t.Errorf("Expected checkout currency to be USD (default), got %s", checkout.Currency)
		}
	})
}

func TestCheckout_MultiCurrencyPricing_Integration(t *testing.T) {
	// Setup mock repositories
	checkoutRepo := mock.NewMockCheckoutRepository()
	currencyRepo := mock.NewMockCurrencyRepository()
	productRepo := mock.NewMockProductRepository()
	productVariantRepo := mock.NewMockProductVariantRepository()
	discountRepo := mock.NewMockDiscountRepository()
	orderRepo := mock.NewMockOrderRepository(false)
	paymentTransactionRepo := mock.NewMockPaymentTransactionRepository()

	// Setup test currencies
	usd, _ := entity.NewCurrency("USD", "US Dollar", "$", 1.0, true, true)
	eur, _ := entity.NewCurrency("EUR", "Euro", "€", 0.85, true, false)
	dkk, _ := entity.NewCurrency("DKK", "Danish Krone", "kr", 6.8, true, false)

	currencyRepo.Create(usd)
	currencyRepo.Create(eur)
	currencyRepo.Create(dkk)

	// Create checkout usecase
	usecase := NewCheckoutUseCase(
		checkoutRepo,
		productRepo,
		productVariantRepo,
		nil, nil, discountRepo, orderRepo, currencyRepo, paymentTransactionRepo, nil, nil,
	)

	t.Run("Should use exact DKK price to avoid conversion precision issues", func(t *testing.T) {
		// Create a checkout in DKK
		checkoutDKK, _ := entity.NewCheckout("test-session-dkk", "DKK")
		checkoutDKK.ID = 1
		checkoutRepo.Create(checkoutDKK)

		// Create a product
		product, _ := entity.NewProduct("Test Product", "Product with exact DKK pricing", "USD", 1, nil)
		product.ID = 1
		product.Active = true
		productRepo.Create(product)

		// Create a variant with USD base price (25.00 USD) and exact DKK price (250.00 DKK)
		variant, _ := entity.NewProductVariant(1, "TEST-SKU-EXACT", 25.0, "USD", 10, nil, nil, true)
		variant.ID = 1

		// Set exact DKK price (250.00 DKK = 25000 cents) to avoid conversion issues
		variant.Prices = []entity.ProductVariantPrice{
			{
				VariantID:    1,
				CurrencyCode: "DKK",
				Price:        25000, // Exactly 250.00 DKK
			},
		}
		productVariantRepo.Create(variant)

		// Add item to checkout
		input := CheckoutInput{
			SKU:      "TEST-SKU-EXACT",
			Quantity: 1,
		}

		updatedCheckout, err := usecase.AddItemToCheckout(1, input)
		if err != nil {
			t.Fatalf("Expected no error adding item to checkout, got %v", err)
		}

		// Check that the item was added with the exact DKK price
		if len(updatedCheckout.Items) != 1 {
			t.Fatalf("Expected 1 item in checkout, got %d", len(updatedCheckout.Items))
		}

		item := updatedCheckout.Items[0]

		// Expected price: exactly 250.00 DKK = 25000 cents (no conversion, no precision loss)
		expectedPriceInCents := int64(25000)
		if item.Price != expectedPriceInCents {
			t.Errorf("Expected item price to be exactly %d DKK cents (250.00 DKK), got %d", expectedPriceInCents, item.Price)
		}

		// Convert to float to verify the exact amount
		itemPriceInDKK := float64(item.Price) / 100
		if itemPriceInDKK != 250.00 {
			t.Errorf("Expected exact price 250.00 DKK, got %.2f DKK", itemPriceInDKK)
		}
	})

	t.Run("Should fallback to conversion when specific currency price not available", func(t *testing.T) {
		// Create a checkout in EUR
		checkoutEUR, _ := entity.NewCheckout("test-session-eur", "EUR")
		checkoutEUR.ID = 2
		checkoutRepo.Create(checkoutEUR)

		// Create a product
		product2, _ := entity.NewProduct("Test Product 2", "Product without EUR pricing", "USD", 1, nil)
		product2.ID = 2
		product2.Active = true
		productRepo.Create(product2)

		// Create a variant with only USD base price (no EUR price set)
		variant2, _ := entity.NewProductVariant(2, "TEST-SKU-FALLBACK", 100.0, "USD", 5, nil, nil, true)
		variant2.ID = 2
		productVariantRepo.Create(variant2)

		// Add item to checkout
		input := CheckoutInput{
			SKU:      "TEST-SKU-FALLBACK",
			Quantity: 1,
		}

		updatedCheckout, err := usecase.AddItemToCheckout(2, input)
		if err != nil {
			t.Fatalf("Expected no error adding item to checkout, got %v", err)
		}

		// Check that the item was added with converted price
		if len(updatedCheckout.Items) != 1 {
			t.Fatalf("Expected 1 item in checkout, got %d", len(updatedCheckout.Items))
		}

		item := updatedCheckout.Items[0]

		// Expected price: 100 USD * 0.85 = 85 EUR = 8500 cents (converted)
		expectedPriceInCents := int64(8500)
		if item.Price != expectedPriceInCents {
			t.Errorf("Expected item price to be %d EUR cents (converted from USD), got %d", expectedPriceInCents, item.Price)
		}
	})

	t.Run("Should handle multiple items with different currency configurations", func(t *testing.T) {
		// Create a checkout in DKK
		checkoutDKK, _ := entity.NewCheckout("test-session-mixed", "DKK")
		checkoutDKK.ID = 3
		checkoutRepo.Create(checkoutDKK)

		// Create products
		product3, _ := entity.NewProduct("Product A", "Product with DKK pricing", "USD", 1, nil)
		product3.ID = 3
		product3.Active = true
		productRepo.Create(product3)

		product4, _ := entity.NewProduct("Product B", "Product without DKK pricing", "USD", 1, nil)
		product4.ID = 4
		product4.Active = true
		productRepo.Create(product4)

		// Variant A: Has specific DKK price
		variantA, _ := entity.NewProductVariant(3, "SKU-A", 50.0, "USD", 10, nil, nil, true)
		variantA.ID = 3
		variantA.Prices = []entity.ProductVariantPrice{
			{
				VariantID:    3,
				CurrencyCode: "DKK",
				Price:        400000, // 4000.00 DKK
			},
		}
		productVariantRepo.Create(variantA)

		// Variant B: No DKK price, will be converted
		variantB, _ := entity.NewProductVariant(4, "SKU-B", 75.0, "USD", 8, nil, nil, true)
		variantB.ID = 4
		productVariantRepo.Create(variantB)

		// Add both items to checkout
		inputA := CheckoutInput{SKU: "SKU-A", Quantity: 1}
		_, err := usecase.AddItemToCheckout(3, inputA)
		if err != nil {
			t.Fatalf("Expected no error adding item A, got %v", err)
		}

		inputB := CheckoutInput{SKU: "SKU-B", Quantity: 1}
		updatedCheckout, err := usecase.AddItemToCheckout(3, inputB)
		if err != nil {
			t.Fatalf("Expected no error adding item B, got %v", err)
		}

		// Check both items
		if len(updatedCheckout.Items) != 2 {
			t.Fatalf("Expected 2 items in checkout, got %d", len(updatedCheckout.Items))
		}

		// Item A should use exact DKK price
		itemA := updatedCheckout.Items[0]
		if itemA.Price != 400000 {
			t.Errorf("Expected item A price to be 400000 DKK cents (exact), got %d", itemA.Price)
		}

		// Item B should use converted price: 75 USD * 6.8 = 510 DKK = 51000 cents
		itemB := updatedCheckout.Items[1]
		expectedPriceBInCents := int64(51000)
		if itemB.Price != expectedPriceBInCents {
			t.Errorf("Expected item B price to be %d DKK cents (converted), got %d", expectedPriceBInCents, itemB.Price)
		}
	})
}
