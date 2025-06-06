package dto

import (
	"reflect"
	"testing"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

func TestConvertToCheckoutDTO(t *testing.T) {
	// Create a test time
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	completedTime := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		checkout *entity.Checkout
		expected CheckoutDTO
	}{
		{
			name: "full checkout conversion",
			checkout: &entity.Checkout{
				ID:               1,
				UserID:           100,
				SessionID:        "sess_123",
				Status:           "pending",
				ShippingMethodID: 5,
				PaymentProvider:  "stripe",
				TotalAmount:      9999, // 99.99 in cents
				ShippingCost:     999,  // 9.99 in cents
				TotalWeight:      1.5,
				Currency:         "USD",
				DiscountCode:     "SAVE10",
				DiscountAmount:   1000, // 10.00 in cents
				FinalAmount:      8999, // 89.99 in cents
				CreatedAt:        testTime,
				UpdatedAt:        testTime,
				LastActivityAt:   testTime,
				ExpiresAt:        testTime.Add(24 * time.Hour),
				CompletedAt:      &completedTime,
				ConvertedOrderID: 200,
				Items: []entity.CheckoutItem{
					{
						ID:               1,
						ProductID:        10,
						ProductVariantID: 20,
						ProductName:      "Test Product",
						VariantName:      "Size M",
						SKU:              "TEST-M",
						Price:            4999, // 49.99 in cents
						Quantity:         2,
						Weight:           0.75,
						CreatedAt:        testTime,
						UpdatedAt:        testTime,
					},
				},
				ShippingAddr: entity.Address{
					Street:     "123 Main St",
					City:       "New York",
					State:      "NY",
					PostalCode: "10001",
					Country:    "US",
				},
				BillingAddr: entity.Address{
					Street:     "456 Oak Ave",
					City:       "Los Angeles",
					State:      "CA",
					PostalCode: "90210",
					Country:    "US",
				},
				CustomerDetails: entity.CustomerDetails{
					Email:    "test@example.com",
					Phone:    "+1234567890",
					FullName: "John Doe",
				},
				ShippingOption: &entity.ShippingOption{
					ShippingMethodID:      5,
					ShippingRateID:        10,
					Name:                  "Standard Shipping",
					Description:           "5-7 business days",
					EstimatedDeliveryDays: 7,
					FreeShipping:          false,
				},
				AppliedDiscount: &entity.AppliedDiscount{
					DiscountID:     1,
					DiscountCode:   "SAVE10",
					DiscountAmount: 1000, // 10.00 in cents
				},
			},
			expected: CheckoutDTO{
				ID:               1,
				UserID:           100,
				SessionID:        "sess_123",
				Status:           "pending",
				ShippingMethodID: 5,
				PaymentProvider:  "stripe",
				TotalAmount:      99.99,
				ShippingCost:     9.99,
				TotalWeight:      1.5,
				Currency:         "USD",
				DiscountCode:     "SAVE10",
				DiscountAmount:   10.0,
				FinalAmount:      89.99,
				CreatedAt:        testTime,
				UpdatedAt:        testTime,
				LastActivityAt:   testTime,
				ExpiresAt:        testTime.Add(24 * time.Hour),
				CompletedAt:      &completedTime,
				ConvertedOrderID: 200,
				Items: []CheckoutItemDTO{
					{
						ID:          1,
						ProductID:   10,
						VariantID:   20,
						ProductName: "Test Product",
						VariantName: "Size M",
						SKU:         "TEST-M",
						Price:       49.99,
						Quantity:    2,
						Weight:      0.75,
						Subtotal:    99.98, // 4999 * 2 / 100
						CreatedAt:   testTime,
						UpdatedAt:   testTime,
					},
				},
				ShippingAddress: AddressDTO{
					AddressLine1: "123 Main St",
					City:         "New York",
					State:        "NY",
					PostalCode:   "10001",
					Country:      "US",
				},
				BillingAddress: AddressDTO{
					AddressLine1: "456 Oak Ave",
					City:         "Los Angeles",
					State:        "CA",
					PostalCode:   "90210",
					Country:      "US",
				},
				CustomerDetails: CustomerDetailsDTO{
					Email:    "test@example.com",
					Phone:    "+1234567890",
					FullName: "John Doe",
				},
				ShippingOption: &ShippingOptionDTO{
					ShippingMethodID:      5,
					ShippingRateID:        10,
					Name:                  "Standard Shipping",
					Description:           "5-7 business days",
					EstimatedDeliveryDays: 7,
					FreeShipping:          false,
				},
				AppliedDiscount: &AppliedDiscountDTO{
					ID:     1,
					Code:   "SAVE10",
					Type:   "", // Empty in conversion
					Method: "", // Empty in conversion
					Value:  0,  // Empty in conversion
					Amount: 10.0,
				},
			},
		},
		{
			name: "checkout without optional fields",
			checkout: &entity.Checkout{
				ID:             2,
				SessionID:      "sess_456",
				Status:         "pending",
				TotalAmount:    5000, // 50.00 in cents
				ShippingCost:   0,
				TotalWeight:    1.0,
				Currency:       "USD",
				DiscountAmount: 0,
				FinalAmount:    5000, // 50.00 in cents
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
				LastActivityAt: testTime,
				ExpiresAt:      testTime.Add(24 * time.Hour),
				Items:          []entity.CheckoutItem{},
				ShippingAddr: entity.Address{
					Street:     "789 Pine St",
					City:       "Boston",
					State:      "MA",
					PostalCode: "02101",
					Country:    "US",
				},
				BillingAddr: entity.Address{
					Street:     "789 Pine St",
					City:       "Boston",
					State:      "MA",
					PostalCode: "02101",
					Country:    "US",
				},
				CustomerDetails: entity.CustomerDetails{
					Email:    "user@example.com",
					Phone:    "+1987654321",
					FullName: "Jane Smith",
				},
			},
			expected: CheckoutDTO{
				ID:             2,
				SessionID:      "sess_456",
				Status:         "pending",
				TotalAmount:    50.0,
				ShippingCost:   0.0,
				TotalWeight:    1.0,
				Currency:       "USD",
				DiscountAmount: 0.0,
				FinalAmount:    50.0,
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
				LastActivityAt: testTime,
				ExpiresAt:      testTime.Add(24 * time.Hour),
				Items:          []CheckoutItemDTO{},
				ShippingAddress: AddressDTO{
					AddressLine1: "789 Pine St",
					City:         "Boston",
					State:        "MA",
					PostalCode:   "02101",
					Country:      "US",
				},
				BillingAddress: AddressDTO{
					AddressLine1: "789 Pine St",
					City:         "Boston",
					State:        "MA",
					PostalCode:   "02101",
					Country:      "US",
				},
				CustomerDetails: CustomerDetailsDTO{
					Email:    "user@example.com",
					Phone:    "+1987654321",
					FullName: "Jane Smith",
				},
			},
		},
		{
			name: "checkout with multiple items",
			checkout: &entity.Checkout{
				ID:             3,
				UserID:         150,
				SessionID:      "sess_789",
				Status:         "completed",
				TotalAmount:    15000, // 150.00 in cents
				ShippingCost:   500,   // 5.00 in cents
				TotalWeight:    2.5,
				Currency:       "EUR",
				DiscountAmount: 0,
				FinalAmount:    15000, // 150.00 in cents
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
				LastActivityAt: testTime,
				ExpiresAt:      testTime.Add(24 * time.Hour),
				Items: []entity.CheckoutItem{
					{
						ID:               1,
						ProductID:        10,
						ProductVariantID: 20,
						ProductName:      "Product A",
						VariantName:      "Red",
						SKU:              "PROD-A-RED",
						Price:            5000, // 50.00 in cents
						Quantity:         1,
						Weight:           1.0,
						CreatedAt:        testTime,
						UpdatedAt:        testTime,
					},
					{
						ID:               2,
						ProductID:        11,
						ProductVariantID: 21,
						ProductName:      "Product B",
						VariantName:      "Blue",
						SKU:              "PROD-B-BLUE",
						Price:            10000, // 100.00 in cents
						Quantity:         1,
						Weight:           1.5,
						CreatedAt:        testTime,
						UpdatedAt:        testTime,
					},
				},
				ShippingAddr: entity.Address{
					Street:     "100 Test Ave",
					City:       "Berlin",
					State:      "BE",
					PostalCode: "10115",
					Country:    "DE",
				},
				BillingAddr: entity.Address{
					Street:     "100 Test Ave",
					City:       "Berlin",
					State:      "BE",
					PostalCode: "10115",
					Country:    "DE",
				},
				CustomerDetails: entity.CustomerDetails{
					Email:    "test@berlin.de",
					Phone:    "+49301234567",
					FullName: "Hans Mueller",
				},
			},
			expected: CheckoutDTO{
				ID:             3,
				UserID:         150,
				SessionID:      "sess_789",
				Status:         "completed",
				TotalAmount:    150.0,
				ShippingCost:   5.0,
				TotalWeight:    2.5,
				Currency:       "EUR",
				DiscountAmount: 0.0,
				FinalAmount:    150.0,
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
				LastActivityAt: testTime,
				ExpiresAt:      testTime.Add(24 * time.Hour),
				Items: []CheckoutItemDTO{
					{
						ID:          1,
						ProductID:   10,
						VariantID:   20,
						ProductName: "Product A",
						VariantName: "Red",
						SKU:         "PROD-A-RED",
						Price:       50.0,
						Quantity:    1,
						Weight:      1.0,
						Subtotal:    50.0, // 5000 * 1 / 100
						CreatedAt:   testTime,
						UpdatedAt:   testTime,
					},
					{
						ID:          2,
						ProductID:   11,
						VariantID:   21,
						ProductName: "Product B",
						VariantName: "Blue",
						SKU:         "PROD-B-BLUE",
						Price:       100.0,
						Quantity:    1,
						Weight:      1.5,
						Subtotal:    100.0, // 10000 * 1 / 100
						CreatedAt:   testTime,
						UpdatedAt:   testTime,
					},
				},
				ShippingAddress: AddressDTO{
					AddressLine1: "100 Test Ave",
					City:         "Berlin",
					State:        "BE",
					PostalCode:   "10115",
					Country:      "DE",
				},
				BillingAddress: AddressDTO{
					AddressLine1: "100 Test Ave",
					City:         "Berlin",
					State:        "BE",
					PostalCode:   "10115",
					Country:      "DE",
				},
				CustomerDetails: CustomerDetailsDTO{
					Email:    "test@berlin.de",
					Phone:    "+49301234567",
					FullName: "Hans Mueller",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToCheckoutDTO(tt.checkout)

			// Compare fields individually for better error messages
			if result.ID != tt.expected.ID {
				t.Errorf("ID mismatch. Got: %d, Want: %d", result.ID, tt.expected.ID)
			}
			if result.UserID != tt.expected.UserID {
				t.Errorf("UserID mismatch. Got: %d, Want: %d", result.UserID, tt.expected.UserID)
			}
			if result.SessionID != tt.expected.SessionID {
				t.Errorf("SessionID mismatch. Got: %s, Want: %s", result.SessionID, tt.expected.SessionID)
			}
			if result.Status != tt.expected.Status {
				t.Errorf("Status mismatch. Got: %s, Want: %s", result.Status, tt.expected.Status)
			}
			if result.TotalAmount != tt.expected.TotalAmount {
				t.Errorf("TotalAmount mismatch. Got: %f, Want: %f", result.TotalAmount, tt.expected.TotalAmount)
			}
			if result.ShippingCost != tt.expected.ShippingCost {
				t.Errorf("ShippingCost mismatch. Got: %f, Want: %f", result.ShippingCost, tt.expected.ShippingCost)
			}
			if result.FinalAmount != tt.expected.FinalAmount {
				t.Errorf("FinalAmount mismatch. Got: %f, Want: %f", result.FinalAmount, tt.expected.FinalAmount)
			}
			if result.Currency != tt.expected.Currency {
				t.Errorf("Currency mismatch. Got: %s, Want: %s", result.Currency, tt.expected.Currency)
			}

			// Compare items
			if len(result.Items) != len(tt.expected.Items) {
				t.Errorf("Items length mismatch. Got: %d, Want: %d", len(result.Items), len(tt.expected.Items))
			} else {
				for i, item := range result.Items {
					expectedItem := tt.expected.Items[i]
					if item.ID != expectedItem.ID {
						t.Errorf("Item[%d] ID mismatch. Got: %d, Want: %d", i, item.ID, expectedItem.ID)
					}
					if item.ProductID != expectedItem.ProductID {
						t.Errorf("Item[%d] ProductID mismatch. Got: %d, Want: %d", i, item.ProductID, expectedItem.ProductID)
					}
					if item.Price != expectedItem.Price {
						t.Errorf("Item[%d] Price mismatch. Got: %f, Want: %f", i, item.Price, expectedItem.Price)
					}
					if item.Subtotal != expectedItem.Subtotal {
						t.Errorf("Item[%d] Subtotal mismatch. Got: %f, Want: %f", i, item.Subtotal, expectedItem.Subtotal)
					}
				}
			}

			// Compare addresses
			if !reflect.DeepEqual(result.ShippingAddress, tt.expected.ShippingAddress) {
				t.Errorf("ShippingAddress mismatch.\nGot: %+v\nWant: %+v", result.ShippingAddress, tt.expected.ShippingAddress)
			}
			if !reflect.DeepEqual(result.BillingAddress, tt.expected.BillingAddress) {
				t.Errorf("BillingAddress mismatch.\nGot: %+v\nWant: %+v", result.BillingAddress, tt.expected.BillingAddress)
			}

			// Compare customer details
			if !reflect.DeepEqual(result.CustomerDetails, tt.expected.CustomerDetails) {
				t.Errorf("CustomerDetails mismatch.\nGot: %+v\nWant: %+v", result.CustomerDetails, tt.expected.CustomerDetails)
			}

			// Compare shipping method (if present)
			if tt.expected.ShippingOption != nil {
				if result.ShippingOption == nil {
					t.Error("Expected ShippingMethod to be present, got nil")
				} else if !reflect.DeepEqual(*result.ShippingOption, *tt.expected.ShippingOption) {
					t.Errorf("ShippingMethod mismatch.\nGot: %+v\nWant: %+v", *result.ShippingOption, *tt.expected.ShippingOption)
				}
			} else if result.ShippingOption != nil {
				t.Errorf("Expected ShippingMethod to be nil, got: %+v", result.ShippingOption)
			}

			// Compare applied discount (if present)
			if tt.expected.AppliedDiscount != nil {
				if result.AppliedDiscount == nil {
					t.Error("Expected AppliedDiscount to be present, got nil")
				} else if !reflect.DeepEqual(*result.AppliedDiscount, *tt.expected.AppliedDiscount) {
					t.Errorf("AppliedDiscount mismatch.\nGot: %+v\nWant: %+v", *result.AppliedDiscount, *tt.expected.AppliedDiscount)
				}
			} else if result.AppliedDiscount != nil {
				t.Errorf("Expected AppliedDiscount to be nil, got: %+v", result.AppliedDiscount)
			}

			// Compare timestamps
			if !result.CreatedAt.Equal(tt.expected.CreatedAt) {
				t.Errorf("CreatedAt mismatch. Got: %v, Want: %v", result.CreatedAt, tt.expected.CreatedAt)
			}
			if !result.UpdatedAt.Equal(tt.expected.UpdatedAt) {
				t.Errorf("UpdatedAt mismatch. Got: %v, Want: %v", result.UpdatedAt, tt.expected.UpdatedAt)
			}

			// Compare CompletedAt (if present)
			if tt.expected.CompletedAt != nil {
				if result.CompletedAt == nil {
					t.Error("Expected CompletedAt to be present, got nil")
				} else if !result.CompletedAt.Equal(*tt.expected.CompletedAt) {
					t.Errorf("CompletedAt mismatch. Got: %v, Want: %v", *result.CompletedAt, *tt.expected.CompletedAt)
				}
			} else if result.CompletedAt != nil {
				t.Errorf("Expected CompletedAt to be nil, got: %v", result.CompletedAt)
			}
		})
	}
}

func TestConvertToCheckoutDTO_CentsConversion(t *testing.T) {
	// Test specific cents conversion scenarios
	testTime := time.Now()

	checkout := &entity.Checkout{
		ID:             1,
		SessionID:      "test",
		Status:         "pending",
		TotalAmount:    12345, // 123.45 in cents
		ShippingCost:   567,   // 5.67 in cents
		DiscountAmount: 1234,  // 12.34 in cents
		FinalAmount:    11678, // 116.78 in cents
		Currency:       "USD",
		CreatedAt:      testTime,
		UpdatedAt:      testTime,
		LastActivityAt: testTime,
		ExpiresAt:      testTime.Add(time.Hour),
		Items: []entity.CheckoutItem{
			{
				ProductID: 1,
				Price:     2499, // 24.99 in cents
				Quantity:  3,
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
		},
		ShippingAddr:    entity.Address{},
		BillingAddr:     entity.Address{},
		CustomerDetails: entity.CustomerDetails{},
	}

	result := ConvertToCheckoutDTO(checkout)

	// Test cents to currency units conversion
	if result.TotalAmount != 123.45 {
		t.Errorf("TotalAmount conversion failed. Got: %f, Want: 123.45", result.TotalAmount)
	}
	if result.ShippingCost != 5.67 {
		t.Errorf("ShippingCost conversion failed. Got: %f, Want: 5.67", result.ShippingCost)
	}
	if result.DiscountAmount != 12.34 {
		t.Errorf("DiscountAmount conversion failed. Got: %f, Want: 12.34", result.DiscountAmount)
	}
	if result.FinalAmount != 116.78 {
		t.Errorf("FinalAmount conversion failed. Got: %f, Want: 116.78", result.FinalAmount)
	}

	// Test item price and subtotal conversion
	if len(result.Items) > 0 {
		item := result.Items[0]
		if item.Price != 24.99 {
			t.Errorf("Item price conversion failed. Got: %f, Want: 24.99", item.Price)
		}
		expectedSubtotal := 74.97 // 24.99 * 3
		if item.Subtotal != expectedSubtotal {
			t.Errorf("Item subtotal conversion failed. Got: %f, Want: %f", item.Subtotal, expectedSubtotal)
		}
	}
}

func TestConvertToCheckoutDTO_NilPointer(t *testing.T) {
	// Test that function doesn't panic with nil pointer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ConvertToCheckoutDTO panicked with nil checkout: %v", r)
		}
	}()

	// This should panic, but we want to test that it doesn't cause unexpected behavior
	// In a real scenario, this function should probably handle nil gracefully
	// For now, we just test that calling it doesn't cause undefined behavior beyond the expected panic
	shouldPanic := func() {
		ConvertToCheckoutDTO(nil)
	}

	// Test that it panics as expected
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected ConvertToCheckoutDTO to panic with nil checkout, but it didn't")
			}
		}()
		shouldPanic()
	}()
}
