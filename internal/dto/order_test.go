package dto

import (
	"testing"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/service"
)

func TestOrderDTO(t *testing.T) {
	now := time.Now()
	items := []OrderItemDTO{
		{
			ID:          1,
			OrderID:     1,
			ProductID:   1,
			VariantID:   1,
			SKU:         "PROD-001",
			ProductName: "Test Product",
			VariantName: "Red/Large",
			Quantity:    2,
			UnitPrice:   29.99,
			TotalPrice:  59.98,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	shippingAddress := AddressDTO{
		AddressLine1: "123 Shipping St",
		City:         "New York",
		State:        "NY",
		PostalCode:   "10001",
		Country:      "US",
	}

	billingAddress := AddressDTO{
		AddressLine1: "456 Billing Ave",
		City:         "Boston",
		State:        "MA",
		PostalCode:   "02101",
		Country:      "US",
	}

	paymentDetails := PaymentDetails{
		PaymentID: "pay_123",
		Provider:  PaymentProviderStripe,
		Method:    PaymentMethodCard,
		Status:    "completed",
		Captured:  true,
		Refunded:  false,
	}

	shippingDetails := ShippingOptionDTO{
		ShippingRateID:        2,
		ShippingMethodID:      1,
		Name:                  "Standard Shipping",
		Description:           "Delivery in 5-7 business days",
		Cost:                  9.99,
		EstimatedDeliveryDays: 5,
		FreeShipping:          false,
	}

	customer := CustomerDetailsDTO{
		Email:    "customer@example.com",
		Phone:    "+1234567890",
		FullName: "John Doe",
	}

	discountDetails := AppliedDiscountDTO{
		Code:   "SAVE10",
		Amount: 10.00,
	}

	order := OrderDTO{
		ID:              1,
		UserID:          1,
		OrderNumber:     "ORD-001",
		Items:           items,
		Status:          OrderStatusPaid,
		PaymentStatus:   PaymentStatusCaptured,
		TotalAmount:     69.97,
		FinalAmount:     59.97,
		Currency:        "USD",
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		PaymentDetails:  paymentDetails,
		ShippingDetails: shippingDetails,
		DiscountDetails: discountDetails,
		Customer:        customer,
		CheckoutID:      "checkout_123",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if order.ID != 1 {
		t.Errorf("Expected ID 1, got %d", order.ID)
	}
	if order.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", order.UserID)
	}
	if order.OrderNumber != "ORD-001" {
		t.Errorf("Expected OrderNumber 'ORD-001', got %s", order.OrderNumber)
	}
	if order.Status != OrderStatusPaid {
		t.Errorf("Expected Status %s, got %s", OrderStatusPaid, order.Status)
	}
	if order.TotalAmount != 69.97 {
		t.Errorf("Expected TotalAmount 69.97, got %f", order.TotalAmount)
	}
	if order.FinalAmount != 59.97 {
		t.Errorf("Expected FinalAmount 59.97, got %f", order.FinalAmount)
	}
	if order.Currency != "USD" {
		t.Errorf("Expected Currency 'USD', got %s", order.Currency)
	}
	if order.CheckoutID != "checkout_123" {
		t.Errorf("Expected CheckoutID 'checkout_123', got %s", order.CheckoutID)
	}
	if len(order.Items) != 1 {
		t.Errorf("Expected Items length 1, got %d", len(order.Items))
	}
	if order.Items[0].ProductName != "Test Product" {
		t.Errorf("Expected Items[0].ProductName 'Test Product', got %s", order.Items[0].ProductName)
	}
}

func TestOrderItemDTO(t *testing.T) {
	now := time.Now()
	item := OrderItemDTO{
		ID:          1,
		OrderID:     1,
		ProductID:   1,
		VariantID:   2,
		SKU:         "PROD-001-VAR",
		ProductName: "Test Product",
		VariantName: "Blue/Medium",
		Quantity:    3,
		UnitPrice:   25.00,
		TotalPrice:  75.00,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if item.ID != 1 {
		t.Errorf("Expected ID 1, got %d", item.ID)
	}
	if item.OrderID != 1 {
		t.Errorf("Expected OrderID 1, got %d", item.OrderID)
	}
	if item.ProductID != 1 {
		t.Errorf("Expected ProductID 1, got %d", item.ProductID)
	}
	if item.VariantID != 2 {
		t.Errorf("Expected VariantID 2, got %d", item.VariantID)
	}
	if item.SKU != "PROD-001-VAR" {
		t.Errorf("Expected SKU 'PROD-001-VAR', got %s", item.SKU)
	}
	if item.ProductName != "Test Product" {
		t.Errorf("Expected ProductName 'Test Product', got %s", item.ProductName)
	}
	if item.VariantName != "Blue/Medium" {
		t.Errorf("Expected VariantName 'Blue/Medium', got %s", item.VariantName)
	}
	if item.Quantity != 3 {
		t.Errorf("Expected Quantity 3, got %d", item.Quantity)
	}
	if item.UnitPrice != 25.00 {
		t.Errorf("Expected UnitPrice 25.00, got %f", item.UnitPrice)
	}
	if item.TotalPrice != 75.00 {
		t.Errorf("Expected TotalPrice 75.00, got %f", item.TotalPrice)
	}
}

func TestPaymentDetails(t *testing.T) {
	details := PaymentDetails{
		PaymentID: "pay_456",
		Provider:  PaymentProviderMobilePay,
		Method:    PaymentMethodWallet,
		Status:    "pending",
		Captured:  false,
		Refunded:  false,
	}

	if details.PaymentID != "pay_456" {
		t.Errorf("Expected PaymentID 'pay_456', got %s", details.PaymentID)
	}
	if details.Provider != PaymentProviderMobilePay {
		t.Errorf("Expected Provider %s, got %s", PaymentProviderMobilePay, details.Provider)
	}
	if details.Method != PaymentMethodWallet {
		t.Errorf("Expected Method %s, got %s", PaymentMethodWallet, details.Method)
	}
	if details.Status != "pending" {
		t.Errorf("Expected Status 'pending', got %s", details.Status)
	}
	if details.Captured {
		t.Errorf("Expected Captured false, got %t", details.Captured)
	}
	if details.Refunded {
		t.Errorf("Expected Refunded false, got %t", details.Refunded)
	}
}

func TestShippingDetails(t *testing.T) {
	details := ShippingOptionDTO{
		ShippingMethodID: 2,
		Name:             "Express Shipping",
		Cost:             19.99,
	}

	if details.ShippingMethodID != 2 {
		t.Errorf("Expected MethodID 2, got %d", details.ShippingMethodID)
	}
	if details.Name != "Express Shipping" {
		t.Errorf("Expected Method 'Express Shipping', got %s", details.Name)
	}
	if details.Cost != 19.99 {
		t.Errorf("Expected Cost 19.99, got %f", details.Cost)
	}
}

func TestCustomerDetails(t *testing.T) {
	customer := CustomerDetailsDTO{
		Email:    "test@example.com",
		Phone:    "+1-555-123-4567",
		FullName: "Jane Smith",
	}

	if customer.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got %s", customer.Email)
	}
	if customer.Phone != "+1-555-123-4567" {
		t.Errorf("Expected Phone '+1-555-123-4567', got %s", customer.Phone)
	}
	if customer.FullName != "Jane Smith" {
		t.Errorf("Expected FullName 'Jane Smith', got %s", customer.FullName)
	}
}

func TestDiscountDetails(t *testing.T) {
	discount := AppliedDiscountDTO{
		Code:   "WINTER20",
		Amount: 15.50,
	}

	if discount.Code != "WINTER20" {
		t.Errorf("Expected Code 'WINTER20', got %s", discount.Code)
	}
	if discount.Amount != 15.50 {
		t.Errorf("Expected Amount 15.50, got %f", discount.Amount)
	}
}

func TestCreateOrderRequest(t *testing.T) {
	shippingAddress := AddressDTO{
		AddressLine1: "789 Test St",
		City:         "Chicago",
		State:        "IL",
		PostalCode:   "60601",
		Country:      "US",
	}

	billingAddress := AddressDTO{
		AddressLine1: "321 Billing Rd",
		City:         "Miami",
		State:        "FL",
		PostalCode:   "33101",
		Country:      "US",
	}

	request := CreateOrderRequest{
		FirstName:        "Alice",
		LastName:         "Johnson",
		Email:            "alice@example.com",
		PhoneNumber:      "+1-555-987-6543",
		ShippingAddress:  shippingAddress,
		BillingAddress:   billingAddress,
		ShippingMethodID: 3,
	}

	if request.FirstName != "Alice" {
		t.Errorf("Expected FirstName 'Alice', got %s", request.FirstName)
	}
	if request.LastName != "Johnson" {
		t.Errorf("Expected LastName 'Johnson', got %s", request.LastName)
	}
	if request.Email != "alice@example.com" {
		t.Errorf("Expected Email 'alice@example.com', got %s", request.Email)
	}
	if request.PhoneNumber != "+1-555-987-6543" {
		t.Errorf("Expected PhoneNumber '+1-555-987-6543', got %s", request.PhoneNumber)
	}
	if request.ShippingMethodID != 3 {
		t.Errorf("Expected ShippingMethodID 3, got %d", request.ShippingMethodID)
	}
	if request.ShippingAddress.City != "Chicago" {
		t.Errorf("Expected ShippingAddress.City 'Chicago', got %s", request.ShippingAddress.City)
	}
	if request.BillingAddress.City != "Miami" {
		t.Errorf("Expected BillingAddress.City 'Miami', got %s", request.BillingAddress.City)
	}
}

func TestCreateOrderItemRequest(t *testing.T) {
	request := CreateOrderItemRequest{
		ProductID: 5,
		VariantID: 3,
		Quantity:  4,
	}

	if request.ProductID != 5 {
		t.Errorf("Expected ProductID 5, got %d", request.ProductID)
	}
	if request.VariantID != 3 {
		t.Errorf("Expected VariantID 3, got %d", request.VariantID)
	}
	if request.Quantity != 4 {
		t.Errorf("Expected Quantity 4, got %d", request.Quantity)
	}
}

func TestUpdateOrderRequest(t *testing.T) {
	estimatedDelivery := time.Now().Add(24 * time.Hour)

	request := UpdateOrderRequest{
		Status:            "shipped",
		PaymentStatus:     "captured",
		TrackingNumber:    "TRACK123456",
		EstimatedDelivery: &estimatedDelivery,
	}

	if request.Status != "shipped" {
		t.Errorf("Expected Status 'shipped', got %s", request.Status)
	}
	if request.PaymentStatus != "captured" {
		t.Errorf("Expected PaymentStatus 'captured', got %s", request.PaymentStatus)
	}
	if request.TrackingNumber != "TRACK123456" {
		t.Errorf("Expected TrackingNumber 'TRACK123456', got %s", request.TrackingNumber)
	}
	if request.EstimatedDelivery == nil {
		t.Error("Expected EstimatedDelivery not nil")
	}
}

func TestOrderSearchRequest(t *testing.T) {
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now()

	request := OrderSearchRequest{
		UserID:        1,
		Status:        OrderStatusPaid,
		PaymentStatus: string(PaymentStatusCaptured),
		StartDate:     &startDate,
		EndDate:       &endDate,
		PaginationDTO: PaginationDTO{
			Page:     1,
			PageSize: 20,
			Total:    0,
		},
	}

	if request.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", request.UserID)
	}
	if request.Status != OrderStatusPaid {
		t.Errorf("Expected Status %s, got %s", OrderStatusPaid, request.Status)
	}
	if request.PaymentStatus != string(PaymentStatusCaptured) {
		t.Errorf("Expected PaymentStatus '%s', got %s", PaymentStatusCaptured, request.PaymentStatus)
	}
	if request.StartDate == nil {
		t.Error("Expected StartDate not nil")
	}
	if request.EndDate == nil {
		t.Error("Expected EndDate not nil")
	}
	if request.Page != 1 {
		t.Errorf("Expected Page 1, got %d", request.Page)
	}
}

func TestProcessPaymentRequest(t *testing.T) {
	cardDetails := &service.CardDetails{
		CardNumber:     "4111111111111111",
		ExpiryMonth:    12,
		ExpiryYear:     2025,
		CVV:            "123",
		CardholderName: "John Doe",
	}

	request := ProcessPaymentRequest{
		PaymentMethod:   PaymentMethodCard,
		PaymentProvider: PaymentProviderStripe,
		CardDetails:     cardDetails,
		PhoneNumber:     "+1-555-123-4567",
	}

	if request.PaymentMethod != PaymentMethodCard {
		t.Errorf("Expected PaymentMethod %s, got %s", PaymentMethodCard, request.PaymentMethod)
	}
	if request.PaymentProvider != PaymentProviderStripe {
		t.Errorf("Expected PaymentProvider %s, got %s", PaymentProviderStripe, request.PaymentProvider)
	}
	if request.CardDetails.CardNumber != "4111111111111111" {
		t.Errorf("Expected CardDetails.CardNumber '4111111111111111', got %s", request.CardDetails.CardNumber)
	}
	if request.PhoneNumber != "+1-555-123-4567" {
		t.Errorf("Expected PhoneNumber '+1-555-123-4567', got %s", request.PhoneNumber)
	}
}

func TestOrderStatusConstants(t *testing.T) {
	if OrderStatusPending != "pending" {
		t.Errorf("Expected OrderStatusPending 'pending', got %s", OrderStatusPending)
	}
	if OrderStatusPaid != "paid" {
		t.Errorf("Expected OrderStatusPaid 'paid', got %s", OrderStatusPaid)
	}
	if OrderStatusShipped != "shipped" {
		t.Errorf("Expected OrderStatusShipped 'shipped', got %s", OrderStatusShipped)
	}
	if OrderStatusCancelled != "cancelled" {
		t.Errorf("Expected OrderStatusCancelled 'cancelled', got %s", OrderStatusCancelled)
	}
	if OrderStatusCompleted != "completed" {
		t.Errorf("Expected OrderStatusCompleted 'completed', got %s", OrderStatusCompleted)
	}
}

func TestPaymentStatusConstants(t *testing.T) {
	if PaymentStatusPending != "pending" {
		t.Errorf("Expected PaymentStatusPending 'pending', got %s", PaymentStatusPending)
	}
	if PaymentStatusAuthorized != "authorized" {
		t.Errorf("Expected PaymentStatusAuthorized 'authorized', got %s", PaymentStatusAuthorized)
	}
	if PaymentStatusCaptured != "captured" {
		t.Errorf("Expected PaymentStatusCaptured 'captured', got %s", PaymentStatusCaptured)
	}
	if PaymentStatusRefunded != "refunded" {
		t.Errorf("Expected PaymentStatusRefunded 'refunded', got %s", PaymentStatusRefunded)
	}
	if PaymentStatusCancelled != "cancelled" {
		t.Errorf("Expected PaymentStatusCancelled 'cancelled', got %s", PaymentStatusCancelled)
	}
	if PaymentStatusFailed != "failed" {
		t.Errorf("Expected PaymentStatusFailed 'failed', got %s", PaymentStatusFailed)
	}
}

func TestPaymentMethodConstants(t *testing.T) {
	if PaymentMethodCard != "credit_card" {
		t.Errorf("Expected PaymentMethodCard 'credit_card', got %s", PaymentMethodCard)
	}
	if PaymentMethodWallet != "wallet" {
		t.Errorf("Expected PaymentMethodWallet 'wallet', got %s", PaymentMethodWallet)
	}
}

func TestPaymentProviderConstants(t *testing.T) {
	if PaymentProviderStripe != "stripe" {
		t.Errorf("Expected PaymentProviderStripe 'stripe', got %s", PaymentProviderStripe)
	}
	if PaymentProviderMobilePay != "mobilepay" {
		t.Errorf("Expected PaymentProviderMobilePay 'mobilepay', got %s", PaymentProviderMobilePay)
	}
}

func TestOrderListResponse(t *testing.T) {
	orders := []OrderSummaryDTO{
		{
			ID:            1,
			OrderNumber:   "ORD-001",
			Status:        OrderStatusPaid,
			PaymentStatus: PaymentStatusCaptured,
			TotalAmount:   99.99,
			Currency:      "USD",
		},
		{
			ID:            2,
			OrderNumber:   "ORD-002",
			Status:        OrderStatusShipped,
			PaymentStatus: PaymentStatusCaptured,
			TotalAmount:   149.99,
			Currency:      "EUR",
		},
	}

	pagination := PaginationDTO{
		Page:     1,
		PageSize: 10,
		Total:    2,
	}

	response := ListResponseDTO[OrderSummaryDTO]{
		Success: true,
		Data:    orders,
		Pagination: PaginationDTO{
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Total:    pagination.Total,
		},
	}

	if !response.Success {
		t.Errorf("Expected Success true, got %t", response.Success)
	}
	if len(response.Data) != 2 {
		t.Errorf("Expected Data length 2, got %d", len(response.Data))
	}
	if response.Data[0].OrderNumber != "ORD-001" {
		t.Errorf("Expected Data[0].OrderNumber 'ORD-001', got %s", response.Data[0].OrderNumber)
	}
	if response.Data[1].Status != OrderStatusShipped {
		t.Errorf("Expected Data[1].Status %s, got %s", OrderStatusShipped, response.Data[1].Status)
	}
	if response.Pagination.Total != 2 {
		t.Errorf("Expected Pagination.Total 2, got %d", response.Pagination.Total)
	}
}
