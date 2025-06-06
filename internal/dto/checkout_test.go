package dto

import (
	"testing"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

func TestCheckoutListResponse(t *testing.T) {
	checkouts := []CheckoutDTO{
		{
			ID:     1,
			UserID: 100,
			Status: "pending",
		},
		{
			ID:     2,
			UserID: 101,
			Status: "completed",
		},
	}

	response := CheckoutListResponse{
		ListResponseDTO: ListResponseDTO[CheckoutDTO]{
			Success: true,
			Data:    checkouts,
			Pagination: PaginationDTO{
				Page:     1,
				PageSize: 10,
				Total:    2,
			},
		},
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 checkouts in response, got %d", len(response.Data))
	}

	if response.Pagination.Total != 2 {
		t.Errorf("Expected total of 2, got %d", response.Pagination.Total)
	}

	if response.Data[0].ID != 1 {
		t.Errorf("Expected first checkout ID to be 1, got %d", response.Data[0].ID)
	}
}

func TestCheckoutDTO(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	checkout := CheckoutDTO{
		ID:               1,
		UserID:           100,
		SessionID:        "session-123",
		Items:            []CheckoutItemDTO{},
		Status:           "active",
		ShippingAddress:  AddressDTO{},
		BillingAddress:   AddressDTO{},
		ShippingMethodID: 1,
		PaymentProvider:  "stripe",
		TotalAmount:      99.99,
		ShippingCost:     9.99,
		TotalWeight:      1.5,
		CustomerDetails:  CustomerDetailsDTO{},
		Currency:         "USD",
		DiscountCode:     "SAVE10",
		DiscountAmount:   10.00,
		FinalAmount:      99.98,
		CreatedAt:        now,
		UpdatedAt:        now,
		LastActivityAt:   now,
		ExpiresAt:        expiresAt,
	}

	// Test basic fields
	if checkout.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", checkout.ID)
	}

	if checkout.UserID != 100 {
		t.Errorf("Expected UserID to be 100, got %d", checkout.UserID)
	}

	if checkout.SessionID != "session-123" {
		t.Errorf("Expected SessionID to be 'session-123', got %s", checkout.SessionID)
	}

	if checkout.Status != "active" {
		t.Errorf("Expected Status to be 'active', got %s", checkout.Status)
	}

	if checkout.TotalAmount != 99.99 {
		t.Errorf("Expected TotalAmount to be 99.99, got %f", checkout.TotalAmount)
	}

	if checkout.Currency != "USD" {
		t.Errorf("Expected Currency to be 'USD', got %s", checkout.Currency)
	}
}

func TestCheckoutItemDTO(t *testing.T) {
	now := time.Now()

	item := CheckoutItemDTO{
		ID:          1,
		ProductID:   10,
		VariantID:   20,
		ProductName: "Test Product",
		VariantName: "Blue / Large",
		ImageURL:    "/images/test.jpg",
		SKU:         "TEST-B-L",
		Price:       29.99,
		Quantity:    2,
		Weight:      0.5,
		Subtotal:    59.98,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test basic fields
	if item.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", item.ID)
	}

	if item.ProductID != 10 {
		t.Errorf("Expected ProductID to be 10, got %d", item.ProductID)
	}

	if item.VariantID != 20 {
		t.Errorf("Expected VariantID to be 20, got %d", item.VariantID)
	}

	if item.ProductName != "Test Product" {
		t.Errorf("Expected ProductName to be 'Test Product', got %s", item.ProductName)
	}

	if item.SKU != "TEST-B-L" {
		t.Errorf("Expected SKU to be 'TEST-B-L', got %s", item.SKU)
	}

	if item.Price != 29.99 {
		t.Errorf("Expected Price to be 29.99, got %f", item.Price)
	}

	if item.Quantity != 2 {
		t.Errorf("Expected Quantity to be 2, got %d", item.Quantity)
	}

	if item.Subtotal != 59.98 {
		t.Errorf("Expected Subtotal to be 59.98, got %f", item.Subtotal)
	}
}

func TestCustomerDetailsDTO(t *testing.T) {
	details := CustomerDetailsDTO{
		Email:    "test@example.com",
		Phone:    "+1234567890",
		FullName: "John Doe",
	}

	if details.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got %s", details.Email)
	}

	if details.Phone != "+1234567890" {
		t.Errorf("Expected Phone to be '+1234567890', got %s", details.Phone)
	}

	if details.FullName != "John Doe" {
		t.Errorf("Expected FullName to be 'John Doe', got %s", details.FullName)
	}
}

func TestAppliedDiscountDTO(t *testing.T) {
	discount := AppliedDiscountDTO{
		ID:     1,
		Code:   "SAVE10",
		Type:   "percentage",
		Method: "basket",
		Value:  10.0,
		Amount: 9.99,
	}

	if discount.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", discount.ID)
	}

	if discount.Code != "SAVE10" {
		t.Errorf("Expected Code to be 'SAVE10', got %s", discount.Code)
	}

	if discount.Type != "percentage" {
		t.Errorf("Expected Type to be 'percentage', got %s", discount.Type)
	}

	if discount.Value != 10.0 {
		t.Errorf("Expected Value to be 10.0, got %f", discount.Value)
	}

	if discount.Amount != 9.99 {
		t.Errorf("Expected Amount to be 9.99, got %f", discount.Amount)
	}
}

func TestAddToCheckoutRequest(t *testing.T) {
	request := AddToCheckoutRequest{
		SKU:      "TEST-B-L",
		Quantity: 2,
	}

	if request.SKU != "TEST-B-L" {
		t.Errorf("Expected SKU to be 'TEST-B-L', got %s", request.SKU)
	}

	if request.Quantity != 2 {
		t.Errorf("Expected Quantity to be 2, got %d", request.Quantity)
	}
}

func TestUpdateCheckoutItemRequest(t *testing.T) {
	request := UpdateCheckoutItemRequest{
		Quantity: 3,
	}

	if request.Quantity != 3 {
		t.Errorf("Expected Quantity to be 3, got %d", request.Quantity)
	}
}

func TestSetShippingAddressRequest(t *testing.T) {
	request := SetShippingAddressRequest{
		AddressLine1: "123 Main St",
		AddressLine2: "Apt 4B",
		City:         "New York",
		State:        "NY",
		PostalCode:   "10001",
		Country:      "USA",
	}

	if request.AddressLine1 != "123 Main St" {
		t.Errorf("Expected AddressLine1 to be '123 Main St', got %s", request.AddressLine1)
	}

	if request.City != "New York" {
		t.Errorf("Expected City to be 'New York', got %s", request.City)
	}

	if request.Country != "USA" {
		t.Errorf("Expected Country to be 'USA', got %s", request.Country)
	}
}

func TestSetCustomerDetailsRequest(t *testing.T) {
	request := SetCustomerDetailsRequest{
		Email:    "customer@example.com",
		Phone:    "+1234567890",
		FullName: "Jane Smith",
	}

	if request.Email != "customer@example.com" {
		t.Errorf("Expected Email to be 'customer@example.com', got %s", request.Email)
	}

	if request.Phone != "+1234567890" {
		t.Errorf("Expected Phone to be '+1234567890', got %s", request.Phone)
	}

	if request.FullName != "Jane Smith" {
		t.Errorf("Expected FullName to be 'Jane Smith', got %s", request.FullName)
	}
}

func TestApplyDiscountRequest(t *testing.T) {
	request := ApplyDiscountRequest{
		DiscountCode: "WELCOME10",
	}

	if request.DiscountCode != "WELCOME10" {
		t.Errorf("Expected DiscountCode to be 'WELCOME10', got %s", request.DiscountCode)
	}
}

func TestSetCurrencyRequest(t *testing.T) {
	request := SetCurrencyRequest{
		Currency: "EUR",
	}

	if request.Currency != "EUR" {
		t.Errorf("Expected Currency to be 'EUR', got %s", request.Currency)
	}
}

func TestCompleteCheckoutRequest(t *testing.T) {
	cardDetails := &CardDetailsDTO{
		CardNumber:     "4111111111111111",
		ExpiryMonth:    12,
		ExpiryYear:     2025,
		CVV:            "123",
		CardholderName: "John Doe",
	}

	request := CompleteCheckoutRequest{
		PaymentProvider: "stripe",
		PaymentData: PaymentData{
			CardDetails: cardDetails,
		},
	}

	if request.PaymentProvider != "stripe" {
		t.Errorf("Expected PaymentProvider to be 'stripe', got %s", request.PaymentProvider)
	}

	if request.PaymentData.CardDetails == nil {
		t.Error("Expected CardDetails to not be nil")
	}

	if request.PaymentData.CardDetails.CardNumber != "4111111111111111" {
		t.Errorf("Expected CardNumber to be '4111111111111111', got %s", request.PaymentData.CardDetails.CardNumber)
	}
}

func TestCardDetailsDTO(t *testing.T) {
	card := CardDetailsDTO{
		CardNumber:     "4111111111111111",
		ExpiryMonth:    12,
		ExpiryYear:     2025,
		CVV:            "123",
		CardholderName: "John Doe",
		Token:          "tok_123456",
	}

	if card.CardNumber != "4111111111111111" {
		t.Errorf("Expected CardNumber to be '4111111111111111', got %s", card.CardNumber)
	}

	if card.ExpiryMonth != 12 {
		t.Errorf("Expected ExpiryMonth to be 12, got %d", card.ExpiryMonth)
	}

	if card.ExpiryYear != 2025 {
		t.Errorf("Expected ExpiryYear to be 2025, got %d", card.ExpiryYear)
	}

	if card.CVV != "123" {
		t.Errorf("Expected CVV to be '123', got %s", card.CVV)
	}

	if card.CardholderName != "John Doe" {
		t.Errorf("Expected CardholderName to be 'John Doe', got %s", card.CardholderName)
	}

	if card.Token != "tok_123456" {
		t.Errorf("Expected Token to be 'tok_123456', got %s", card.Token)
	}
}

func TestConvertToCheckoutDTO_MinimalCheckout(t *testing.T) {
	now := time.Now()

	// Create a minimal checkout entity with only required fields
	checkout := &entity.Checkout{
		ID:              1,
		Status:          entity.CheckoutStatusActive,
		Currency:        "USD",
		TotalAmount:     0,
		ShippingCost:    0,
		FinalAmount:     0,
		CreatedAt:       now,
		UpdatedAt:       now,
		LastActivityAt:  now,
		ExpiresAt:       now.Add(24 * time.Hour),
		Items:           []entity.CheckoutItem{},
		ShippingAddr:    entity.Address{},
		BillingAddr:     entity.Address{},
		CustomerDetails: entity.CustomerDetails{},
	}

	dto := ConvertToCheckoutDTO(checkout)

	// Test that conversion doesn't fail with minimal data
	if dto.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", dto.ID)
	}

	if dto.Status != "active" {
		t.Errorf("Expected Status to be 'active', got %s", dto.Status)
	}

	if dto.Currency != "USD" {
		t.Errorf("Expected Currency to be 'USD', got %s", dto.Currency)
	}

	if len(dto.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(dto.Items))
	}

	if dto.ShippingOption != nil {
		t.Error("Expected shipping method to be nil")
	}

	if dto.AppliedDiscount != nil {
		t.Error("Expected applied discount to be nil")
	}

	if dto.CompletedAt != nil {
		t.Error("Expected CompletedAt to be nil")
	}

	if dto.ConvertedOrderID != 0 {
		t.Errorf("Expected ConvertedOrderID to be 0, got %d", dto.ConvertedOrderID)
	}
}

func TestConvertToCheckoutDTO_MultipleItems(t *testing.T) {
	now := time.Now()

	checkout := &entity.Checkout{
		ID:             1,
		Status:         entity.CheckoutStatusActive,
		Currency:       "USD",
		TotalAmount:    7998, // 79.98 in cents
		CreatedAt:      now,
		UpdatedAt:      now,
		LastActivityAt: now,
		ExpiresAt:      now.Add(24 * time.Hour),
		Items: []entity.CheckoutItem{
			{
				ID:               1,
				ProductID:        10,
				ProductVariantID: 20,
				ProductName:      "Product 1",
				VariantName:      "Red / Small",
				SKU:              "PROD1-R-S",
				Price:            1999, // 19.99 in cents
				Quantity:         2,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			{
				ID:               2,
				ProductID:        11,
				ProductVariantID: 21,
				ProductName:      "Product 2",
				VariantName:      "Blue / Large",
				SKU:              "PROD2-B-L",
				Price:            2000, // 20.00 in cents
				Quantity:         2,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
		},
		ShippingAddr:    entity.Address{},
		BillingAddr:     entity.Address{},
		CustomerDetails: entity.CustomerDetails{},
	}

	dto := ConvertToCheckoutDTO(checkout)

	// Test multiple items conversion
	if len(dto.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(dto.Items))
	}

	// Test first item
	item1 := dto.Items[0]
	if item1.SKU != "PROD1-R-S" {
		t.Errorf("Expected first item SKU to be 'PROD1-R-S', got %s", item1.SKU)
	}

	if item1.Price != 19.99 {
		t.Errorf("Expected first item price to be 19.99, got %f", item1.Price)
	}

	if item1.Subtotal != 39.98 {
		t.Errorf("Expected first item subtotal to be 39.98, got %f", item1.Subtotal)
	}

	// Test second item
	item2 := dto.Items[1]
	if item2.SKU != "PROD2-B-L" {
		t.Errorf("Expected second item SKU to be 'PROD2-B-L', got %s", item2.SKU)
	}

	if item2.Price != 20.00 {
		t.Errorf("Expected second item price to be 20.00, got %f", item2.Price)
	}

	if item2.Subtotal != 40.00 {
		t.Errorf("Expected second item subtotal to be 40.00, got %f", item2.Subtotal)
	}
}

func TestCheckoutCompleteResponse(t *testing.T) {
	now := time.Now()

	order := OrderDTO{
		ID:          1,
		Status:      "confirmed",
		TotalAmount: 99.99,
		Currency:    "USD",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := CheckoutCompleteResponse{
		Order:          order,
		ActionRequired: true,
		ActionURL:      "https://payment.example.com/confirm",
	}

	if response.Order.ID != 1 {
		t.Errorf("Expected Order ID to be 1, got %d", response.Order.ID)
	}

	if !response.ActionRequired {
		t.Error("Expected ActionRequired to be true")
	}

	if response.ActionURL != "https://payment.example.com/confirm" {
		t.Errorf("Expected ActionURL to be 'https://payment.example.com/confirm', got %s", response.ActionURL)
	}
}
