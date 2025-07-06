package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// CheckoutRepository defines the interface for checkout data access
type CheckoutRepository interface {
	// Create creates a new checkout
	Create(checkout *entity.Checkout) error

	// GetByID retrieves a checkout by ID
	GetByID(checkoutID uint) (*entity.Checkout, error)

	// GetByUserID retrieves an active checkout by user ID
	GetByUserID(userID uint) (*entity.Checkout, error)

	// GetBySessionID retrieves an active checkout by session ID
	GetBySessionID(sessionID string) (*entity.Checkout, error)

	// GetAbandonedBySessionID retrieves an abandoned checkout by session ID
	GetAbandonedBySessionID(sessionID string) (*entity.Checkout, error)

	// Update updates a checkout
	Update(checkout *entity.Checkout) error

	// Delete deletes a checkout
	Delete(checkoutID uint) error

	// ConvertGuestCheckoutToUserCheckout converts a guest checkout to a user checkout
	ConvertGuestCheckoutToUserCheckout(sessionID string, userID uint) (*entity.Checkout, error)

	// GetExpiredCheckouts retrieves all checkouts that have expired
	GetExpiredCheckouts() ([]*entity.Checkout, error)

	// GetCheckoutsToAbandon retrieves active checkouts with customer/shipping info that should be marked as abandoned
	GetCheckoutsToAbandon() ([]*entity.Checkout, error)

	// GetCheckoutsToDelete retrieves checkouts that should be deleted (empty checkouts > 24h or abandoned > 7 days)
	GetCheckoutsToDelete() ([]*entity.Checkout, error)

	// GetCheckoutsByStatus retrieves checkouts by status
	GetCheckoutsByStatus(status entity.CheckoutStatus, offset, limit int) ([]*entity.Checkout, error)

	// GetActiveCheckoutsByUserID retrieves all active checkouts for a user
	GetActiveCheckoutsByUserID(userID uint) ([]*entity.Checkout, error)

	// GetCompletedCheckoutsByUserID retrieves all completed checkouts for a user
	GetCompletedCheckoutsByUserID(userID uint, offset, limit int) ([]*entity.Checkout, error)

	// HasActiveCheckoutsWithProduct checks if a product has any active checkouts
	HasActiveCheckoutsWithProduct(productID uint) (bool, error)

	// GetAllExpiredCheckoutsForDeletion retrieves all expired checkouts for force deletion
	GetAllExpiredCheckoutsForDeletion() ([]*entity.Checkout, error)
}
