package dto

import (
	"reflect"
	"testing"
	"time"
)

func TestCheckoutDTO(t *testing.T) {
	tests := []struct {
		name     string
		dto      CheckoutDTO
		expected CheckoutDTO
	}{
		{
			name: "full checkout DTO",
			dto: CheckoutDTO{
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
						Subtotal:    99.98,
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
				CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				LastActivityAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				ExpiresAt:      time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
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
						Subtotal:    99.98,
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
				CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				LastActivityAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				ExpiresAt:      time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "empty checkout DTO",
			dto:  CheckoutDTO{},
			expected: CheckoutDTO{
				Items: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.dto, tt.expected) {
				t.Errorf("CheckoutDTO mismatch.\nGot: %+v\nWant: %+v", tt.dto, tt.expected)
			}
		})
	}
}

func TestCheckoutItemDTO(t *testing.T) {
	tests := []struct {
		name     string
		dto      CheckoutItemDTO
		expected CheckoutItemDTO
	}{
		{
			name: "full checkout item DTO",
			dto: CheckoutItemDTO{
				ID:          1,
				ProductID:   10,
				VariantID:   20,
				ProductName: "Test Product",
				VariantName: "Size M",
				SKU:         "TEST-M",
				Price:       49.99,
				Quantity:    2,
				Weight:      0.75,
				Subtotal:    99.98,
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expected: CheckoutItemDTO{
				ID:          1,
				ProductID:   10,
				VariantID:   20,
				ProductName: "Test Product",
				VariantName: "Size M",
				SKU:         "TEST-M",
				Price:       49.99,
				Quantity:    2,
				Weight:      0.75,
				Subtotal:    99.98,
				CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:     "empty checkout item DTO",
			dto:      CheckoutItemDTO{},
			expected: CheckoutItemDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.dto, tt.expected) {
				t.Errorf("CheckoutItemDTO mismatch.\nGot: %+v\nWant: %+v", tt.dto, tt.expected)
			}
		})
	}
}

func TestCustomerDetailsDTO(t *testing.T) {
	tests := []struct {
		name     string
		dto      CustomerDetailsDTO
		expected CustomerDetailsDTO
	}{
		{
			name: "full customer details DTO",
			dto: CustomerDetailsDTO{
				Email:    "test@example.com",
				Phone:    "+1234567890",
				FullName: "John Doe",
			},
			expected: CustomerDetailsDTO{
				Email:    "test@example.com",
				Phone:    "+1234567890",
				FullName: "John Doe",
			},
		},
		{
			name:     "empty customer details DTO",
			dto:      CustomerDetailsDTO{},
			expected: CustomerDetailsDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.dto, tt.expected) {
				t.Errorf("CustomerDetailsDTO mismatch.\nGot: %+v\nWant: %+v", tt.dto, tt.expected)
			}
		})
	}
}

func TestAppliedDiscountDTO(t *testing.T) {
	tests := []struct {
		name     string
		dto      AppliedDiscountDTO
		expected AppliedDiscountDTO
	}{
		{
			name: "full applied discount DTO",
			dto: AppliedDiscountDTO{
				ID:     1,
				Code:   "SAVE10",
				Type:   "percentage",
				Method: "fixed",
				Value:  10.0,
				Amount: 5.99,
			},
			expected: AppliedDiscountDTO{
				ID:     1,
				Code:   "SAVE10",
				Type:   "percentage",
				Method: "fixed",
				Value:  10.0,
				Amount: 5.99,
			},
		},
		{
			name:     "empty applied discount DTO",
			dto:      AppliedDiscountDTO{},
			expected: AppliedDiscountDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.dto, tt.expected) {
				t.Errorf("AppliedDiscountDTO mismatch.\nGot: %+v\nWant: %+v", tt.dto, tt.expected)
			}
		})
	}
}

func TestAddToCheckoutRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  AddToCheckoutRequest
		expected AddToCheckoutRequest
	}{
		{
			name: "full add to checkout request",
			request: AddToCheckoutRequest{
				SKU:      "TEST-M",
				Quantity: 2,
			},
			expected: AddToCheckoutRequest{
				SKU:      "TEST-M",
				Quantity: 2,
			},
		},
		{
			name: "without variant ID",
			request: AddToCheckoutRequest{
				SKU:      "TEST-M",
				Quantity: 1,
			},
			expected: AddToCheckoutRequest{
				SKU:      "TEST-M",
				Quantity: 1,
			},
		},
		{
			name:     "empty add to checkout request",
			request:  AddToCheckoutRequest{},
			expected: AddToCheckoutRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("AddToCheckoutRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestUpdateCheckoutItemRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateCheckoutItemRequest
		expected UpdateCheckoutItemRequest
	}{
		{
			name: "full update checkout item request",
			request: UpdateCheckoutItemRequest{
				Quantity: 3,
				SKU:      "TEST-M",
			},
			expected: UpdateCheckoutItemRequest{
				Quantity: 3,
				SKU:      "TEST-M",
			},
		},
		{
			name: "quantity only",
			request: UpdateCheckoutItemRequest{
				Quantity: 5,
			},
			expected: UpdateCheckoutItemRequest{
				Quantity: 5,
			},
		},
		{
			name:     "empty update checkout item request",
			request:  UpdateCheckoutItemRequest{},
			expected: UpdateCheckoutItemRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("UpdateCheckoutItemRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestSetShippingAddressRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  SetShippingAddressRequest
		expected SetShippingAddressRequest
	}{
		{
			name: "full shipping address request",
			request: SetShippingAddressRequest{
				AddressLine1: "123 Main St",
				AddressLine2: "Apt 4B",
				City:         "New York",
				State:        "NY",
				PostalCode:   "10001",
				Country:      "US",
			},
			expected: SetShippingAddressRequest{
				AddressLine1: "123 Main St",
				AddressLine2: "Apt 4B",
				City:         "New York",
				State:        "NY",
				PostalCode:   "10001",
				Country:      "US",
			},
		},
		{
			name: "without address line 2",
			request: SetShippingAddressRequest{
				AddressLine1: "456 Oak Ave",
				City:         "Los Angeles",
				State:        "CA",
				PostalCode:   "90210",
				Country:      "US",
			},
			expected: SetShippingAddressRequest{
				AddressLine1: "456 Oak Ave",
				City:         "Los Angeles",
				State:        "CA",
				PostalCode:   "90210",
				Country:      "US",
			},
		},
		{
			name:     "empty shipping address request",
			request:  SetShippingAddressRequest{},
			expected: SetShippingAddressRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("SetShippingAddressRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestSetBillingAddressRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  SetBillingAddressRequest
		expected SetBillingAddressRequest
	}{
		{
			name: "full billing address request",
			request: SetBillingAddressRequest{
				AddressLine1: "789 Pine St",
				AddressLine2: "Suite 100",
				City:         "Chicago",
				State:        "IL",
				PostalCode:   "60601",
				Country:      "US",
			},
			expected: SetBillingAddressRequest{
				AddressLine1: "789 Pine St",
				AddressLine2: "Suite 100",
				City:         "Chicago",
				State:        "IL",
				PostalCode:   "60601",
				Country:      "US",
			},
		},
		{
			name:     "empty billing address request",
			request:  SetBillingAddressRequest{},
			expected: SetBillingAddressRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("SetBillingAddressRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestSetCustomerDetailsRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  SetCustomerDetailsRequest
		expected SetCustomerDetailsRequest
	}{
		{
			name: "full customer details request",
			request: SetCustomerDetailsRequest{
				Email:    "test@example.com",
				Phone:    "+1234567890",
				FullName: "Jane Smith",
			},
			expected: SetCustomerDetailsRequest{
				Email:    "test@example.com",
				Phone:    "+1234567890",
				FullName: "Jane Smith",
			},
		},
		{
			name:     "empty customer details request",
			request:  SetCustomerDetailsRequest{},
			expected: SetCustomerDetailsRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("SetCustomerDetailsRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestSetShippingMethodRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  SetShippingMethodRequest
		expected SetShippingMethodRequest
	}{
		{
			name: "valid shipping method request",
			request: SetShippingMethodRequest{
				ShippingMethodID: 5,
			},
			expected: SetShippingMethodRequest{
				ShippingMethodID: 5,
			},
		},
		{
			name:     "empty shipping method request",
			request:  SetShippingMethodRequest{},
			expected: SetShippingMethodRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("SetShippingMethodRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestSetCurrencyRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  SetCurrencyRequest
		expected SetCurrencyRequest
	}{
		{
			name: "valid currency request",
			request: SetCurrencyRequest{
				Currency: "EUR",
			},
			expected: SetCurrencyRequest{
				Currency: "EUR",
			},
		},
		{
			name:     "empty currency request",
			request:  SetCurrencyRequest{},
			expected: SetCurrencyRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("SetCurrencyRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestApplyDiscountRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  ApplyDiscountRequest
		expected ApplyDiscountRequest
	}{
		{
			name: "valid discount request",
			request: ApplyDiscountRequest{
				DiscountCode: "SAVE20",
			},
			expected: ApplyDiscountRequest{
				DiscountCode: "SAVE20",
			},
		},
		{
			name:     "empty discount request",
			request:  ApplyDiscountRequest{},
			expected: ApplyDiscountRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("ApplyDiscountRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestCheckoutSearchRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  CheckoutSearchRequest
		expected CheckoutSearchRequest
	}{
		{
			name: "full checkout search request",
			request: CheckoutSearchRequest{
				UserID: 100,
				Status: "pending",
				PaginationDTO: PaginationDTO{
					Page:     2,
					PageSize: 20,
				},
			},
			expected: CheckoutSearchRequest{
				UserID: 100,
				Status: "pending",
				PaginationDTO: PaginationDTO{
					Page:     2,
					PageSize: 20,
				},
			},
		},
		{
			name: "with user ID only",
			request: CheckoutSearchRequest{
				UserID: 50,
			},
			expected: CheckoutSearchRequest{
				UserID: 50,
			},
		},
		{
			name:     "empty checkout search request",
			request:  CheckoutSearchRequest{},
			expected: CheckoutSearchRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("CheckoutSearchRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestCheckoutCompleteResponse(t *testing.T) {
	tests := []struct {
		name     string
		response CheckoutCompleteResponse
		expected CheckoutCompleteResponse
	}{
		{
			name: "complete response with action required",
			response: CheckoutCompleteResponse{
				Order: OrderDTO{
					ID:     100,
					UserID: 50,
					Status: "pending",
				},
				ActionRequired: true,
				ActionURL:      "https://payment.example.com/confirm",
			},
			expected: CheckoutCompleteResponse{
				Order: OrderDTO{
					ID:     100,
					UserID: 50,
					Status: "pending",
				},
				ActionRequired: true,
				ActionURL:      "https://payment.example.com/confirm",
			},
		},
		{
			name: "complete response without action",
			response: CheckoutCompleteResponse{
				Order: OrderDTO{
					ID:     101,
					UserID: 51,
					Status: "completed",
				},
				ActionRequired: false,
			},
			expected: CheckoutCompleteResponse{
				Order: OrderDTO{
					ID:     101,
					UserID: 51,
					Status: "completed",
				},
				ActionRequired: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.response, tt.expected) {
				t.Errorf("CheckoutCompleteResponse mismatch.\nGot: %+v\nWant: %+v", tt.response, tt.expected)
			}
		})
	}
}

func TestCompleteCheckoutRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  CompleteCheckoutRequest
		expected CompleteCheckoutRequest
	}{
		{
			name: "complete checkout with card details",
			request: CompleteCheckoutRequest{
				PaymentProvider: "stripe",
				PaymentData: PaymentData{
					CardDetails: &CardDetailsDTO{
						CardNumber:     "4111111111111111",
						ExpiryMonth:    12,
						ExpiryYear:     2025,
						CVV:            "123",
						CardholderName: "John Doe",
					},
				},
			},
			expected: CompleteCheckoutRequest{
				PaymentProvider: "stripe",
				PaymentData: PaymentData{
					CardDetails: &CardDetailsDTO{
						CardNumber:     "4111111111111111",
						ExpiryMonth:    12,
						ExpiryYear:     2025,
						CVV:            "123",
						CardholderName: "John Doe",
					},
				},
			},
		},
		{
			name: "complete checkout with phone number",
			request: CompleteCheckoutRequest{
				PaymentProvider: "mpesa",
				PaymentData: PaymentData{
					PhoneNumber: "+254700000000",
				},
			},
			expected: CompleteCheckoutRequest{
				PaymentProvider: "mpesa",
				PaymentData: PaymentData{
					PhoneNumber: "+254700000000",
				},
			},
		},
		{
			name:     "empty complete checkout request",
			request:  CompleteCheckoutRequest{},
			expected: CompleteCheckoutRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.request, tt.expected) {
				t.Errorf("CompleteCheckoutRequest mismatch.\nGot: %+v\nWant: %+v", tt.request, tt.expected)
			}
		})
	}
}

func TestPaymentData(t *testing.T) {
	tests := []struct {
		name     string
		data     PaymentData
		expected PaymentData
	}{
		{
			name: "payment data with card details",
			data: PaymentData{
				CardDetails: &CardDetailsDTO{
					CardNumber:     "4000000000000002",
					ExpiryMonth:    6,
					ExpiryYear:     2026,
					CVV:            "456",
					CardholderName: "Jane Smith",
					Token:          "tok_123456789",
				},
			},
			expected: PaymentData{
				CardDetails: &CardDetailsDTO{
					CardNumber:     "4000000000000002",
					ExpiryMonth:    6,
					ExpiryYear:     2026,
					CVV:            "456",
					CardholderName: "Jane Smith",
					Token:          "tok_123456789",
				},
			},
		},
		{
			name: "payment data with phone number",
			data: PaymentData{
				PhoneNumber: "+254711111111",
			},
			expected: PaymentData{
				PhoneNumber: "+254711111111",
			},
		},
		{
			name:     "empty payment data",
			data:     PaymentData{},
			expected: PaymentData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.data, tt.expected) {
				t.Errorf("PaymentData mismatch.\nGot: %+v\nWant: %+v", tt.data, tt.expected)
			}
		})
	}
}

func TestCardDetailsDTO(t *testing.T) {
	tests := []struct {
		name     string
		dto      CardDetailsDTO
		expected CardDetailsDTO
	}{
		{
			name: "full card details DTO",
			dto: CardDetailsDTO{
				CardNumber:     "5555555555554444",
				ExpiryMonth:    3,
				ExpiryYear:     2027,
				CVV:            "789",
				CardholderName: "Alice Johnson",
				Token:          "tok_987654321",
			},
			expected: CardDetailsDTO{
				CardNumber:     "5555555555554444",
				ExpiryMonth:    3,
				ExpiryYear:     2027,
				CVV:            "789",
				CardholderName: "Alice Johnson",
				Token:          "tok_987654321",
			},
		},
		{
			name: "card details without token",
			dto: CardDetailsDTO{
				CardNumber:     "4242424242424242",
				ExpiryMonth:    8,
				ExpiryYear:     2028,
				CVV:            "321",
				CardholderName: "Bob Wilson",
			},
			expected: CardDetailsDTO{
				CardNumber:     "4242424242424242",
				ExpiryMonth:    8,
				ExpiryYear:     2028,
				CVV:            "321",
				CardholderName: "Bob Wilson",
			},
		},
		{
			name:     "empty card details DTO",
			dto:      CardDetailsDTO{},
			expected: CardDetailsDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.dto, tt.expected) {
				t.Errorf("CardDetailsDTO mismatch.\nGot: %+v\nWant: %+v", tt.dto, tt.expected)
			}
		})
	}
}

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
