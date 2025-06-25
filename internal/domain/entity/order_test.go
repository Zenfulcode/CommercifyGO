package entity

import (
	"testing"
	"time"
)

func TestOrderConstants(t *testing.T) {
	// Test OrderStatus constants
	if OrderStatusPending != "pending" {
		t.Errorf("Expected OrderStatusPending to be 'pending', got %s", OrderStatusPending)
	}
	if OrderStatusPaid != "paid" {
		t.Errorf("Expected OrderStatusPaid to be 'paid', got %s", OrderStatusPaid)
	}
	if OrderStatusShipped != "shipped" {
		t.Errorf("Expected OrderStatusShipped to be 'shipped', got %s", OrderStatusShipped)
	}
	if OrderStatusCancelled != "cancelled" {
		t.Errorf("Expected OrderStatusCancelled to be 'cancelled', got %s", OrderStatusCancelled)
	}
	if OrderStatusCompleted != "completed" {
		t.Errorf("Expected OrderStatusCompleted to be 'completed', got %s", OrderStatusCompleted)
	}

	// Test PaymentStatus constants
	if PaymentStatusPending != "pending" {
		t.Errorf("Expected PaymentStatusPending to be 'pending', got %s", PaymentStatusPending)
	}
	if PaymentStatusAuthorized != "authorized" {
		t.Errorf("Expected PaymentStatusAuthorized to be 'authorized', got %s", PaymentStatusAuthorized)
	}
	if PaymentStatusCaptured != "captured" {
		t.Errorf("Expected PaymentStatusCaptured to be 'captured', got %s", PaymentStatusCaptured)
	}
	if PaymentStatusRefunded != "refunded" {
		t.Errorf("Expected PaymentStatusRefunded to be 'refunded', got %s", PaymentStatusRefunded)
	}
	if PaymentStatusCancelled != "cancelled" {
		t.Errorf("Expected PaymentStatusCancelled to be 'cancelled', got %s", PaymentStatusCancelled)
	}
	if PaymentStatusFailed != "failed" {
		t.Errorf("Expected PaymentStatusFailed to be 'failed', got %s", PaymentStatusFailed)
	}
}

func TestNewOrder(t *testing.T) {
	// Test valid order creation
	items := []OrderItem{
		{
			ProductID: 1,
			Quantity:  2,
			Price:     1000, // $10.00
			Weight:    0.5,
		},
		{
			ProductID: 2,
			Quantity:  1,
			Price:     2000, // $20.00
			Weight:    1.0,
		},
	}

	shippingAddr := Address{
		Street:     "123 Main St",
		City:       "New York",
		State:      "NY",
		PostalCode: "10001",
		Country:    "USA",
	}

	billingAddr := Address{
		Street:     "456 Oak Ave",
		City:       "Los Angeles",
		State:      "CA",
		PostalCode: "90210",
		Country:    "USA",
	}

	customerDetails := CustomerDetails{
		Email:    "test@example.com",
		Phone:    "+1234567890",
		FullName: "John Doe",
	}

	order, err := NewOrder(1, items, "USD", shippingAddr, billingAddr, customerDetails)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify order properties
	if order.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", order.UserID)
	}
	if order.Currency != "USD" {
		t.Errorf("Expected Currency 'USD', got %s", order.Currency)
	}
	if order.TotalAmount != 4000 { // (2*1000) + (1*2000) = 4000
		t.Errorf("Expected TotalAmount 4000, got %d", order.TotalAmount)
	}
	if order.FinalAmount != 4000 {
		t.Errorf("Expected FinalAmount 4000, got %d", order.FinalAmount)
	}
	if order.TotalWeight != 2.0 { // (2*0.5) + (1*1.0) = 2.0
		t.Errorf("Expected TotalWeight 2.0, got %f", order.TotalWeight)
	}
	if order.Status != OrderStatusPending {
		t.Errorf("Expected Status %s, got %s", OrderStatusPending, order.Status)
	}
	if order.PaymentStatus != PaymentStatusPending {
		t.Errorf("Expected PaymentStatus %s, got %s", PaymentStatusPending, order.PaymentStatus)
	}
	if order.IsGuestOrder {
		t.Errorf("Expected IsGuestOrder false, got true")
	}
	if order.CustomerDetails.Email != "test@example.com" {
		t.Errorf("Expected customer email 'test@example.com', got %s", order.CustomerDetails.Email)
	}
	if len(order.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(order.Items))
	}

	// Verify order number format
	expectedPrefix := "ORD-" + time.Now().Format("20060102")
	if !contains(order.OrderNumber, expectedPrefix) {
		t.Errorf("Expected order number to contain %s, got %s", expectedPrefix, order.OrderNumber)
	}
}

func TestNewOrderValidation(t *testing.T) {
	items := []OrderItem{
		{ProductID: 1, Quantity: 1, Price: 1000, Weight: 0.5},
	}
	addr := Address{Street: "123 Main St", City: "NYC", State: "NY", PostalCode: "10001", Country: "USA"}
	customer := CustomerDetails{Email: "test@example.com", Phone: "+1234567890", FullName: "John Doe"}

	// Test zero user ID
	_, err := NewOrder(0, items, "USD", addr, addr, customer)
	if err == nil {
		t.Error("Expected error for zero user ID")
	}

	// Test empty items
	_, err = NewOrder(1, []OrderItem{}, "USD", addr, addr, customer)
	if err == nil {
		t.Error("Expected error for empty items")
	}

	// Test empty currency
	_, err = NewOrder(1, items, "", addr, addr, customer)
	if err == nil {
		t.Error("Expected error for empty currency")
	}

	// Test zero quantity
	invalidItems := []OrderItem{
		{ProductID: 1, Quantity: 0, Price: 1000, Weight: 0.5},
	}
	_, err = NewOrder(1, invalidItems, "USD", addr, addr, customer)
	if err == nil {
		t.Error("Expected error for zero quantity")
	}

	// Test zero price
	invalidItems = []OrderItem{
		{ProductID: 1, Quantity: 1, Price: 0, Weight: 0.5},
	}
	_, err = NewOrder(1, invalidItems, "USD", addr, addr, customer)
	if err == nil {
		t.Error("Expected error for zero price")
	}
}

func TestNewGuestOrder(t *testing.T) {
	items := []OrderItem{
		{ProductID: 1, Quantity: 1, Price: 1500, Weight: 0.8},
	}

	shippingAddr := Address{
		Street:     "789 Guest St",
		City:       "Miami",
		State:      "FL",
		PostalCode: "33101",
		Country:    "USA",
	}

	billingAddr := Address{
		Street:     "789 Guest St",
		City:       "Miami",
		State:      "FL",
		PostalCode: "33101",
		Country:    "USA",
	}

	customerDetails := CustomerDetails{
		Email:    "guest@example.com",
		Phone:    "+1987654321",
		FullName: "Guest User",
	}

	order, err := NewGuestOrder(items, shippingAddr, billingAddr, customerDetails)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify guest order properties
	if order.UserID != 0 {
		t.Errorf("Expected UserID 0 for guest order, got %d", order.UserID)
	}
	if !order.IsGuestOrder {
		t.Errorf("Expected IsGuestOrder true, got false")
	}
	if order.TotalAmount != 1500 {
		t.Errorf("Expected TotalAmount 1500, got %d", order.TotalAmount)
	}
	if order.Status != OrderStatusPending {
		t.Errorf("Expected Status %s, got %s", OrderStatusPending, order.Status)
	}
	if order.PaymentStatus != PaymentStatusPending {
		t.Errorf("Expected PaymentStatus %s, got %s", PaymentStatusPending, order.PaymentStatus)
	}

	// Verify order number format for guest orders
	expectedPrefix := "GS-" + time.Now().Format("20060102")
	if !contains(order.OrderNumber, expectedPrefix) {
		t.Errorf("Expected guest order number to contain %s, got %s", expectedPrefix, order.OrderNumber)
	}
}

func TestUpdateStatus(t *testing.T) {
	order := createTestOrder(t)

	// Test valid transitions
	testCases := []struct {
		name       string
		fromStatus OrderStatus
		toStatus   OrderStatus
		shouldErr  bool
	}{
		{"Pending to Paid", OrderStatusPending, OrderStatusPaid, false},
		{"Pending to Cancelled", OrderStatusPending, OrderStatusCancelled, false},
		{"Paid to Shipped", OrderStatusPaid, OrderStatusShipped, false},
		{"Paid to Cancelled", OrderStatusPaid, OrderStatusCancelled, false},
		{"Shipped to Completed", OrderStatusShipped, OrderStatusCompleted, false},
		{"Shipped to Cancelled", OrderStatusShipped, OrderStatusCancelled, false},
		// Invalid transitions
		{"Pending to Shipped", OrderStatusPending, OrderStatusShipped, true},
		{"Pending to Completed", OrderStatusPending, OrderStatusCompleted, true},
		{"Paid to Completed", OrderStatusPaid, OrderStatusCompleted, true},
		{"Cancelled to Any", OrderStatusCancelled, OrderStatusPaid, true},
		{"Completed to Any", OrderStatusCompleted, OrderStatusPaid, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset order status
			order.Status = tc.fromStatus
			order.CompletedAt = nil

			err := order.UpdateStatus(tc.toStatus)

			if tc.shouldErr && err == nil {
				t.Errorf("Expected error for transition %s -> %s", tc.fromStatus, tc.toStatus)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("Unexpected error for transition %s -> %s: %v", tc.fromStatus, tc.toStatus, err)
			}

			if !tc.shouldErr {
				if order.Status != tc.toStatus {
					t.Errorf("Expected status %s, got %s", tc.toStatus, order.Status)
				}

				// Check if completed_at is set for terminal states
				if tc.toStatus == OrderStatusCancelled || tc.toStatus == OrderStatusCompleted {
					if order.CompletedAt == nil {
						t.Errorf("Expected CompletedAt to be set for status %s", tc.toStatus)
					}
				}
			}
		})
	}
}

func TestUpdatePaymentStatus(t *testing.T) {
	testCases := []struct {
		name                string
		fromPaymentStatus   PaymentStatus
		toPaymentStatus     PaymentStatus
		initialOrderStatus  OrderStatus
		expectedOrderStatus OrderStatus
		shouldErr           bool
		shouldSetCompleted  bool
	}{
		{
			name:                "Pending to Authorized",
			fromPaymentStatus:   PaymentStatusPending,
			toPaymentStatus:     PaymentStatusAuthorized,
			initialOrderStatus:  OrderStatusPending,
			expectedOrderStatus: OrderStatusPaid,
			shouldErr:           false,
		},
		{
			name:                "Pending to Failed",
			fromPaymentStatus:   PaymentStatusPending,
			toPaymentStatus:     PaymentStatusFailed,
			initialOrderStatus:  OrderStatusPending,
			expectedOrderStatus: OrderStatusCancelled,
			shouldErr:           false,
			shouldSetCompleted:  true,
		},
		{
			name:                "Authorized to Captured (Shipped Order)",
			fromPaymentStatus:   PaymentStatusAuthorized,
			toPaymentStatus:     PaymentStatusCaptured,
			initialOrderStatus:  OrderStatusShipped,
			expectedOrderStatus: OrderStatusCompleted,
			shouldErr:           false,
			shouldSetCompleted:  true,
		},
		{
			name:                "Authorized to Cancelled",
			fromPaymentStatus:   PaymentStatusAuthorized,
			toPaymentStatus:     PaymentStatusCancelled,
			initialOrderStatus:  OrderStatusPaid,
			expectedOrderStatus: OrderStatusCancelled,
			shouldErr:           false,
			shouldSetCompleted:  true,
		},
		{
			name:                "Captured to Refunded",
			fromPaymentStatus:   PaymentStatusCaptured,
			toPaymentStatus:     PaymentStatusRefunded,
			initialOrderStatus:  OrderStatusCompleted,
			expectedOrderStatus: OrderStatusCompleted, // Order status doesn't change on refund
			shouldErr:           false,
		},
		// Invalid transitions
		{
			name:              "Pending to Captured (invalid)",
			fromPaymentStatus: PaymentStatusPending,
			toPaymentStatus:   PaymentStatusCaptured,
			shouldErr:         true,
		},
		{
			name:              "Failed to any (invalid)",
			fromPaymentStatus: PaymentStatusFailed,
			toPaymentStatus:   PaymentStatusAuthorized,
			shouldErr:         true,
		},
		{
			name:              "Refunded to any (invalid)",
			fromPaymentStatus: PaymentStatusRefunded,
			toPaymentStatus:   PaymentStatusCaptured,
			shouldErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := createTestOrder(t)
			order.PaymentStatus = tc.fromPaymentStatus
			order.Status = tc.initialOrderStatus
			order.CompletedAt = nil

			err := order.UpdatePaymentStatus(tc.toPaymentStatus)

			if tc.shouldErr && err == nil {
				t.Errorf("Expected error for payment transition %s -> %s", tc.fromPaymentStatus, tc.toPaymentStatus)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("Unexpected error for payment transition %s -> %s: %v", tc.fromPaymentStatus, tc.toPaymentStatus, err)
			}

			if !tc.shouldErr {
				if order.PaymentStatus != tc.toPaymentStatus {
					t.Errorf("Expected payment status %s, got %s", tc.toPaymentStatus, order.PaymentStatus)
				}

				if tc.expectedOrderStatus != "" && order.Status != tc.expectedOrderStatus {
					t.Errorf("Expected order status %s, got %s", tc.expectedOrderStatus, order.Status)
				}

				if tc.shouldSetCompleted && order.CompletedAt == nil {
					t.Errorf("Expected CompletedAt to be set")
				}
			}
		})
	}
}

func TestOrderSetters(t *testing.T) {
	order := createTestOrder(t)

	// Test SetPaymentID
	err := order.SetPaymentID("payment_12345")
	if err != nil {
		t.Errorf("Unexpected error setting payment ID: %v", err)
	}
	if order.PaymentID != "payment_12345" {
		t.Errorf("Expected payment ID 'payment_12345', got %s", order.PaymentID)
	}

	// Test SetPaymentID with empty value
	err = order.SetPaymentID("")
	if err == nil {
		t.Error("Expected error for empty payment ID")
	}

	// Test SetPaymentProvider
	err = order.SetPaymentProvider("stripe")
	if err != nil {
		t.Errorf("Unexpected error setting payment provider: %v", err)
	}
	if order.PaymentProvider != "stripe" {
		t.Errorf("Expected payment provider 'stripe', got %s", order.PaymentProvider)
	}

	// Test SetPaymentProvider with empty value
	err = order.SetPaymentProvider("")
	if err == nil {
		t.Error("Expected error for empty payment provider")
	}

	// Test SetPaymentMethod
	err = order.SetPaymentMethod("card")
	if err != nil {
		t.Errorf("Unexpected error setting payment method: %v", err)
	}
	if order.PaymentMethod != "card" {
		t.Errorf("Expected payment method 'card', got %s", order.PaymentMethod)
	}

	// Test SetTrackingCode
	err = order.SetTrackingCode("TRACK123456")
	if err != nil {
		t.Errorf("Unexpected error setting tracking code: %v", err)
	}
	if order.TrackingCode != "TRACK123456" {
		t.Errorf("Expected tracking code 'TRACK123456', got %s", order.TrackingCode)
	}

	// Test SetActionURL
	err = order.SetActionURL("https://payment.example.com/checkout")
	if err != nil {
		t.Errorf("Unexpected error setting action URL: %v", err)
	}
	if order.ActionURL != "https://payment.example.com/checkout" {
		t.Errorf("Expected action URL 'https://payment.example.com/checkout', got %s", order.ActionURL)
	}
}

func TestSetOrderNumber(t *testing.T) {
	order := createTestOrder(t)
	orderID := uint(12345)

	order.SetOrderNumber(orderID)

	expectedOrderNumber := "ORD-" + order.CreatedAt.Format("20060102") + "-012345"
	if order.OrderNumber != expectedOrderNumber {
		t.Errorf("Expected order number %s, got %s", expectedOrderNumber, order.OrderNumber)
	}
}

func TestSetShippingMethod(t *testing.T) {
	order := createTestOrder(t)
	originalFinalAmount := order.FinalAmount

	shippingOption := &ShippingOption{
		ShippingMethodID:      1,
		Name:                  "Express Shipping",
		Cost:                  500, // $5.00
		EstimatedDeliveryDays: 2,
	}

	err := order.SetShippingMethod(shippingOption)
	if err != nil {
		t.Errorf("Unexpected error setting shipping method: %v", err)
	}

	if order.ShippingMethodID != 1 {
		t.Errorf("Expected shipping method ID 1, got %d", order.ShippingMethodID)
	}
	if order.ShippingCost != 500 {
		t.Errorf("Expected shipping cost 500, got %d", order.ShippingCost)
	}
	if order.FinalAmount != originalFinalAmount+500 {
		t.Errorf("Expected final amount %d, got %d", originalFinalAmount+500, order.FinalAmount)
	}
	if order.ShippingOption == nil || order.ShippingOption.Name != "Express Shipping" {
		t.Errorf("Expected shipping option to be set correctly")
	}

	// Test with nil shipping option
	err = order.SetShippingMethod(nil)
	if err == nil {
		t.Error("Expected error for nil shipping option")
	}
}

func TestCalculateTotalWeight(t *testing.T) {
	order := createTestOrder(t)

	// Modify items for testing
	order.Items = []OrderItem{
		{ProductID: 1, Quantity: 2, Price: 1000, Weight: 0.5}, // 2 * 0.5 = 1.0
		{ProductID: 2, Quantity: 3, Price: 1500, Weight: 1.2}, // 3 * 1.2 = 3.6
	}

	totalWeight := order.CalculateTotalWeight()
	expectedWeight := 4.6 // 1.0 + 3.6

	if totalWeight != expectedWeight {
		t.Errorf("Expected total weight %.2f, got %.2f", expectedWeight, totalWeight)
	}
	if order.TotalWeight != expectedWeight {
		t.Errorf("Expected order total weight %.2f, got %.2f", expectedWeight, order.TotalWeight)
	}
}

func TestIsCaptured(t *testing.T) {
	order := createTestOrder(t)

	// Test when not captured
	order.PaymentStatus = PaymentStatusPending
	if order.IsCaptured() {
		t.Error("Expected IsCaptured to be false for pending payment")
	}

	// Test when captured
	order.PaymentStatus = PaymentStatusCaptured
	if !order.IsCaptured() {
		t.Error("Expected IsCaptured to be true for captured payment")
	}
}

func TestIsRefunded(t *testing.T) {
	order := createTestOrder(t)

	// Test when not refunded
	order.PaymentStatus = PaymentStatusCaptured
	if order.IsRefunded() {
		t.Error("Expected IsRefunded to be false for captured payment")
	}

	// Test when refunded
	order.PaymentStatus = PaymentStatusRefunded
	if !order.IsRefunded() {
		t.Error("Expected IsRefunded to be true for refunded payment")
	}
}

func TestApplyDiscount(t *testing.T) {
	order := createTestOrder(t)
	order.TotalAmount = 10000 // $100.00
	order.FinalAmount = 10000
	order.ShippingCost = 500 // $5.00

	// Create a test discount
	discount := &Discount{
		ID:           1,
		Code:         "SAVE10",
		Type:         DiscountTypeBasket,
		Method:       DiscountMethodPercentage,
		Value:        10.0, // 10% off
		Active:       true,
		StartDate:    time.Now().Add(-24 * time.Hour),
		EndDate:      time.Now().Add(24 * time.Hour),
		UsageLimit:   100,
		CurrentUsage: 5,
	}

	err := order.ApplyDiscount(discount)
	if err != nil {
		t.Errorf("Unexpected error applying discount: %v", err)
	}

	expectedDiscountAmount := int64(1000) // 10% of $100.00
	if order.DiscountAmount != expectedDiscountAmount {
		t.Errorf("Expected discount amount %d, got %d", expectedDiscountAmount, order.DiscountAmount)
	}

	expectedFinalAmount := order.TotalAmount + order.ShippingCost - expectedDiscountAmount
	if order.FinalAmount != expectedFinalAmount {
		t.Errorf("Expected final amount %d, got %d", expectedFinalAmount, order.FinalAmount)
	}

	if order.AppliedDiscount == nil {
		t.Error("Expected applied discount to be set")
	} else {
		if order.AppliedDiscount.DiscountID != discount.ID {
			t.Errorf("Expected applied discount ID %d, got %d", discount.ID, order.AppliedDiscount.DiscountID)
		}
		if order.AppliedDiscount.DiscountCode != discount.Code {
			t.Errorf("Expected applied discount code %s, got %s", discount.Code, order.AppliedDiscount.DiscountCode)
		}
	}

	// Test applying nil discount
	err = order.ApplyDiscount(nil)
	if err == nil {
		t.Error("Expected error for nil discount")
	}
}

func TestRemoveDiscount(t *testing.T) {
	order := createTestOrder(t)
	order.TotalAmount = 10000
	order.ShippingCost = 500
	order.DiscountAmount = 1000
	order.FinalAmount = 9500 // 10000 + 500 - 1000
	order.AppliedDiscount = &AppliedDiscount{
		DiscountID:     1,
		DiscountCode:   "SAVE10",
		DiscountAmount: 1000,
	}

	order.RemoveDiscount()

	if order.DiscountAmount != 0 {
		t.Errorf("Expected discount amount 0, got %d", order.DiscountAmount)
	}
	if order.FinalAmount != 10500 { // 10000 + 500
		t.Errorf("Expected final amount 10500, got %d", order.FinalAmount)
	}
	if order.AppliedDiscount != nil {
		t.Error("Expected applied discount to be nil")
	}
}

// Helper functions

func createTestOrder(t *testing.T) *Order {
	items := []OrderItem{
		{ProductID: 1, Quantity: 1, Price: 1000, Weight: 0.5},
	}

	addr := Address{
		Street:     "123 Test St",
		City:       "Test City",
		State:      "TS",
		PostalCode: "12345",
		Country:    "USA",
	}

	customer := CustomerDetails{
		Email:    "test@example.com",
		Phone:    "+1234567890",
		FullName: "Test User",
	}

	order, err := NewOrder(1, items, "USD", addr, addr, customer)
	if err != nil {
		t.Fatalf("Failed to create test order: %v", err)
	}

	return order
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
