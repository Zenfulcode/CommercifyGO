package gorm

import (
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

type CheckoutRepository struct {
	db *gorm.DB
}

// ConvertGuestCheckoutToUserCheckout implements repository.CheckoutRepository.
func (c *CheckoutRepository) ConvertGuestCheckoutToUserCheckout(sessionID string, userID uint) (*entity.Checkout, error) {
	var checkout entity.Checkout

	// First, find the guest checkout
	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Where("session_id = ? AND user_id IS NULL OR user_id = 0", sessionID).
		First(&checkout).Error
	if err != nil {
		return nil, fmt.Errorf("guest checkout not found: %w", err)
	}

	// Update the checkout to assign it to the user
	checkout.UserID = &userID
	checkout.LastActivityAt = time.Now()

	err = c.db.Save(&checkout).Error
	if err != nil {
		return nil, fmt.Errorf("failed to convert guest checkout to user checkout: %w", err)
	}

	return &checkout, nil
}

// Create implements repository.CheckoutRepository.
func (c *CheckoutRepository) Create(checkout *entity.Checkout) error {
	return c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").Create(checkout).Error
}

// Delete implements repository.CheckoutRepository.
func (c *CheckoutRepository) Delete(checkoutID uint) error {
	return c.db.Delete(&entity.Checkout{}, checkoutID).Error
}

// GetActiveCheckoutsByUserID implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetActiveCheckoutsByUserID(userID uint) ([]*entity.Checkout, error) {
	var checkouts []*entity.Checkout

	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Where("user_id = ? AND status = ?", userID, entity.CheckoutStatusActive).
		Order("created_at DESC").
		Find(&checkouts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active checkouts by user ID: %w", err)
	}

	return checkouts, nil
}

// GetByID implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetByID(checkoutID uint) (*entity.Checkout, error) {
	var checkout entity.Checkout
	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").Preload("ConvertedOrder").
		First(&checkout, checkoutID).Error
	if err != nil {
		return nil, err
	}
	return &checkout, nil
}

// GetBySessionID implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetBySessionID(sessionID string) (*entity.Checkout, error) {
	var checkout entity.Checkout
	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Where("session_id = ? AND status = ?", sessionID, entity.CheckoutStatusActive).
		First(&checkout).Error
	if err != nil {
		return nil, err
	}
	return &checkout, nil
}

// GetByUserID implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetByUserID(userID uint) (*entity.Checkout, error) {
	var checkout entity.Checkout
	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Where("user_id = ? AND status = ?", userID, entity.CheckoutStatusActive).
		First(&checkout).Error
	if err != nil {
		return nil, err
	}
	return &checkout, nil
}

// GetCheckoutsByStatus implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetCheckoutsByStatus(status entity.CheckoutStatus, offset int, limit int) ([]*entity.Checkout, error) {
	var checkouts []*entity.Checkout

	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Where("status = ?", status).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&checkouts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkouts by status: %w", err)
	}

	return checkouts, nil
}

// GetCheckoutsToAbandon implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetCheckoutsToAbandon() ([]*entity.Checkout, error) {
	var checkouts []*entity.Checkout
	abandonThreshold := time.Now().Add(-15 * time.Minute)

	// Find active checkouts with customer/shipping info that haven't been active for 15 minutes
	// Check if there's any customer details or shipping address data (JSON fields are not empty/null)
	err := c.db.Preload("Items").
		Where("status = ? AND last_activity_at < ? AND (customer_email != '' OR customer_phone != '' OR customer_full_name != '' OR shipping_address IS NOT NULL)",
			entity.CheckoutStatusActive, abandonThreshold).
		Find(&checkouts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkouts to abandon: %w", err)
	}

	return checkouts, nil
}

// GetCheckoutsToDelete implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetCheckoutsToDelete() ([]*entity.Checkout, error) {
	var checkouts []*entity.Checkout
	now := time.Now()

	// Delete empty checkouts after 24 hours OR abandoned checkouts after 7 days OR all expired checkouts
	emptyThreshold := now.Add(-24 * time.Hour)
	abandonedThreshold := now.Add(-7 * 24 * time.Hour)

	err := c.db.Where(
		"(customer_email = '' AND customer_phone = '' AND customer_full_name = '' AND shipping_address IS NULL AND last_activity_at < ?) OR "+
			"(status = ? AND updated_at < ?) OR "+
			"status = ?",
		emptyThreshold, entity.CheckoutStatusAbandoned, abandonedThreshold, entity.CheckoutStatusExpired).
		Find(&checkouts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkouts to delete: %w", err)
	}

	return checkouts, nil
}

// GetCompletedCheckoutsByUserID implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetCompletedCheckoutsByUserID(userID uint, offset int, limit int) ([]*entity.Checkout, error) {
	var checkouts []*entity.Checkout

	err := c.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").Preload("ConvertedOrder").
		Where("user_id = ? AND status = ?", userID, entity.CheckoutStatusCompleted).
		Offset(offset).Limit(limit).
		Order("completed_at DESC").
		Find(&checkouts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch completed checkouts by user ID: %w", err)
	}

	return checkouts, nil
}

// GetExpiredCheckouts implements repository.CheckoutRepository.
func (c *CheckoutRepository) GetExpiredCheckouts() ([]*entity.Checkout, error) {
	var checkouts []*entity.Checkout
	now := time.Now()

	err := c.db.Preload("Items").
		Where("status = ? AND expires_at < ?", entity.CheckoutStatusActive, now).
		Find(&checkouts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expired checkouts: %w", err)
	}

	return checkouts, nil
}

// HasActiveCheckoutsWithProduct implements repository.CheckoutRepository.
func (c *CheckoutRepository) HasActiveCheckoutsWithProduct(productID uint) (bool, error) {
	var count int64

	err := c.db.Model(&entity.Checkout{}).
		Joins("JOIN checkout_items ON checkouts.id = checkout_items.checkout_id").
		Where("checkouts.status = ? AND checkout_items.product_id = ?", entity.CheckoutStatusActive, productID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check active checkouts with product: %w", err)
	}

	return count > 0, nil
}

// Update implements repository.CheckoutRepository.
func (c *CheckoutRepository) Update(checkout *entity.Checkout) error {
	return c.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(checkout).Error
}

// NewCheckoutRepository creates a new GORM-based CheckoutRepository
func NewCheckoutRepository(db *gorm.DB) repository.CheckoutRepository {
	return &CheckoutRepository{db: db}
}
