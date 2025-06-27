package usecase

import (
	"testing"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/testutil/mock"
)

// Simple mock services for testing stock management
type mockPaymentService struct{}

func (m *mockPaymentService) GetAvailableProviders() []service.PaymentProvider { return nil }
func (m *mockPaymentService) GetAvailableProvidersForCurrency(currency string) []service.PaymentProvider {
	return nil
}
func (m *mockPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	return nil, nil
}
func (m *mockPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	return false, nil
}
func (m *mockPaymentService) CapturePayment(transactionID, currency string, amount int64, provider service.PaymentProviderType) (*service.PaymentResult, error) {
	return nil, nil
}
func (m *mockPaymentService) RefundPayment(transactionID, currency string, amount int64, provider service.PaymentProviderType) (*service.PaymentResult, error) {
	return nil, nil
}
func (m *mockPaymentService) CancelPayment(transactionID string, provider service.PaymentProviderType) (*service.PaymentResult, error) {
	return nil, nil
}
func (m *mockPaymentService) ForceApprovePayment(transactionID, phoneNumber string, provider service.PaymentProviderType) error {
	return nil
}

type mockEmailService struct{}

func (m *mockEmailService) SendEmail(data service.EmailData) error { return nil }
func (m *mockEmailService) SendOrderConfirmation(order *entity.Order, user *entity.User) error {
	return nil
}
func (m *mockEmailService) SendOrderNotification(order *entity.Order, user *entity.User) error {
	return nil
}

func TestOrderUseCase_HandleStockUpdatesForPaymentStatusChange(t *testing.T) {
	tests := []struct {
		name           string
		previousStatus entity.PaymentStatus
		newStatus      entity.PaymentStatus
		initialStock   int
		orderQuantity  int
		expectedStock  int
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "Stock decreased when payment authorized",
			previousStatus: entity.PaymentStatusPending,
			newStatus:      entity.PaymentStatusAuthorized,
			initialStock:   10,
			orderQuantity:  2,
			expectedStock:  8,
			expectError:    false,
		},
		{
			name:           "Stock increased when authorized payment cancelled",
			previousStatus: entity.PaymentStatusAuthorized,
			newStatus:      entity.PaymentStatusCancelled,
			initialStock:   8,
			orderQuantity:  2,
			expectedStock:  10,
			expectError:    false,
		},
		{
			name:           "Stock increased when authorized payment failed",
			previousStatus: entity.PaymentStatusAuthorized,
			newStatus:      entity.PaymentStatusFailed,
			initialStock:   8,
			orderQuantity:  2,
			expectedStock:  10,
			expectError:    false,
		},
		{
			name:           "Stock increased when captured payment refunded",
			previousStatus: entity.PaymentStatusCaptured,
			newStatus:      entity.PaymentStatusRefunded,
			initialStock:   8,
			orderQuantity:  2,
			expectedStock:  10,
			expectError:    false,
		},
		{
			name:           "No stock change for authorized to captured",
			previousStatus: entity.PaymentStatusAuthorized,
			newStatus:      entity.PaymentStatusCaptured,
			initialStock:   8,
			orderQuantity:  2,
			expectedStock:  8,
			expectError:    false,
		},
		{
			name:           "No stock change for pending to cancelled",
			previousStatus: entity.PaymentStatusPending,
			newStatus:      entity.PaymentStatusCancelled,
			initialStock:   10,
			orderQuantity:  2,
			expectedStock:  10,
			expectError:    false,
		},
		{
			name:           "No stock change for pending to failed",
			previousStatus: entity.PaymentStatusPending,
			newStatus:      entity.PaymentStatusFailed,
			initialStock:   10,
			orderQuantity:  2,
			expectedStock:  10,
			expectError:    false,
		},
		{
			name:           "Error when insufficient stock on authorization",
			previousStatus: entity.PaymentStatusPending,
			newStatus:      entity.PaymentStatusAuthorized,
			initialStock:   1,
			orderQuantity:  2,
			expectedStock:  1,
			expectError:    true,
			errorMessage:   "insufficient stock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test repositories
			orderRepo := mock.NewMockOrderRepository(false)
			productRepo := mock.NewMockProductRepository()
			productVariantRepo := mock.NewMockProductVariantRepository()
			userRepo := mock.NewMockUserRepository()
			paymentTxnRepo := mock.NewMockPaymentTransactionRepository()
			currencyRepo := mock.NewMockCurrencyRepository()

			// Setup payment service and email service mocks
			paymentSvc := &mockPaymentService{}
			emailSvc := &mockEmailService{}

			// Create the use case
			uc := NewOrderUseCase(
				orderRepo,
				productRepo,
				productVariantRepo,
				userRepo,
				paymentSvc,
				emailSvc,
				paymentTxnRepo,
				currencyRepo,
			)

			// Create a simple test order with pre-configured variant
			variant := &entity.ProductVariant{
				ProductID: 1,
				SKU:       "TEST-SKU",
				Stock:     tt.initialStock,
			}

			// Create the variant in the repository
			if err := productVariantRepo.Create(variant); err != nil {
				t.Fatalf("Failed to create variant: %v", err)
			}

			order := &entity.Order{
				ID: 1,
				Items: []entity.OrderItem{
					{
						ProductVariantID: variant.ID, // Use the ID assigned by the mock repository
						Quantity:         tt.orderQuantity,
						ProductName:      "Test Product",
						SKU:              "TEST-SKU",
					},
				},
				PaymentStatus: tt.previousStatus,
			}

			// Test the stock update logic
			stockErr := uc.handleStockUpdatesForPaymentStatusChange(order, tt.previousStatus, tt.newStatus)

			// Check error expectation
			if tt.expectError {
				if stockErr == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMessage != "" && !contains(stockErr.Error(), tt.errorMessage) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorMessage, stockErr)
				}
			} else if stockErr != nil {
				t.Errorf("Unexpected error: %v", stockErr)
			}

			// Check stock level
			updatedVariant, err := productVariantRepo.GetByID(variant.ID)
			if err != nil {
				t.Fatalf("Failed to get updated variant: %v", err)
			}

			if updatedVariant.Stock != tt.expectedStock {
				t.Errorf("Expected stock to be %d, got %d", tt.expectedStock, updatedVariant.Stock)
			}
		})
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		len(str) > len(substr) && (str[:len(substr)] == substr ||
			str[len(str)-len(substr):] == substr))
}
