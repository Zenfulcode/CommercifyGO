package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// MockCheckoutRepository is a mock implementation of the CheckoutRepository interface
type MockCheckoutRepository struct {
	mutex            sync.Mutex
	checkouts        map[uint]*entity.Checkout
	userCheckouts    map[uint]*entity.Checkout
	sessionCheckouts map[string]*entity.Checkout
	nextID           uint
}

// NewMockCheckoutRepository creates a new mock checkout repository
func NewMockCheckoutRepository() repository.CheckoutRepository {
	return &MockCheckoutRepository{
		checkouts:        make(map[uint]*entity.Checkout),
		userCheckouts:    make(map[uint]*entity.Checkout),
		sessionCheckouts: make(map[string]*entity.Checkout),
		nextID:           1,
	}
}

// Create adds a checkout to the repository
func (r *MockCheckoutRepository) Create(checkout *entity.Checkout) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	checkout.ID = r.nextID
	r.nextID++

	// Add to checkouts map
	r.checkouts[checkout.ID] = checkout

	// Store checkout based on whether it's a user checkout or guest checkout
	if checkout.SessionID != "" {
		r.sessionCheckouts[checkout.SessionID] = checkout
	}

	if checkout.UserID > 0 {
		r.userCheckouts[checkout.UserID] = checkout
	}

	return nil
}

// GetByID retrieves a checkout by ID
func (r *MockCheckoutRepository) GetByID(checkoutID uint) (*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if checkout, found := r.checkouts[checkoutID]; found {
		return checkout, nil
	}

	return nil, errors.New("checkout not found")
}

// GetByUserID retrieves an active checkout by user ID
func (r *MockCheckoutRepository) GetByUserID(userID uint) (*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if checkout, found := r.userCheckouts[userID]; found && checkout.Status == entity.CheckoutStatusActive {
		return checkout, nil
	}

	return nil, errors.New("active checkout not found for user")
}

// GetBySessionID retrieves an active checkout by session ID
func (r *MockCheckoutRepository) GetBySessionID(sessionID string) (*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if checkout, found := r.sessionCheckouts[sessionID]; found && checkout.Status == entity.CheckoutStatusActive {
		return checkout, nil
	}

	return nil, errors.New("active checkout not found for session")
}

// Update updates a checkout
func (r *MockCheckoutRepository) Update(checkout *entity.Checkout) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.checkouts[checkout.ID]; !found {
		return errors.New("checkout not found")
	}

	checkout.UpdatedAt = time.Now()
	r.checkouts[checkout.ID] = checkout

	// Update in the appropriate map
	if checkout.SessionID != "" {
		r.sessionCheckouts[checkout.SessionID] = checkout
	}

	if checkout.UserID > 0 {
		r.userCheckouts[checkout.UserID] = checkout
	}

	return nil
}

// Delete deletes a checkout
func (r *MockCheckoutRepository) Delete(checkoutID uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	checkout, found := r.checkouts[checkoutID]
	if !found {
		return errors.New("checkout not found")
	}

	// Remove from the maps
	delete(r.checkouts, checkoutID)

	if checkout.SessionID != "" {
		delete(r.sessionCheckouts, checkout.SessionID)
	}

	if checkout.UserID > 0 {
		delete(r.userCheckouts, checkout.UserID)
	}

	return nil
}

// ConvertGuestCheckoutToUserCheckout converts a guest checkout to a user checkout
func (r *MockCheckoutRepository) ConvertGuestCheckoutToUserCheckout(sessionID string, userID uint) (*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	checkout, found := r.sessionCheckouts[sessionID]
	if !found {
		return nil, errors.New("guest checkout not found")
	}

	// Update the checkout
	checkout.UserID = userID
	checkout.UpdatedAt = time.Now()

	// Store in user checkouts map
	r.userCheckouts[userID] = checkout

	return checkout, nil
}

// GetExpiredCheckouts retrieves all checkouts that have expired
func (r *MockCheckoutRepository) GetExpiredCheckouts() ([]*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var expiredCheckouts []*entity.Checkout
	now := time.Now()

	for _, checkout := range r.checkouts {
		if checkout.Status == entity.CheckoutStatusActive && checkout.ExpiresAt.Before(now) {
			expiredCheckouts = append(expiredCheckouts, checkout)
		}
	}

	return expiredCheckouts, nil
}

// GetCheckoutsByStatus retrieves checkouts by status
func (r *MockCheckoutRepository) GetCheckoutsByStatus(status entity.CheckoutStatus, offset, limit int) ([]*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var matchingCheckouts []*entity.Checkout

	for _, checkout := range r.checkouts {
		if checkout.Status == status {
			matchingCheckouts = append(matchingCheckouts, checkout)
		}
	}

	// Apply offset and limit
	start := offset
	if start >= len(matchingCheckouts) {
		return []*entity.Checkout{}, nil
	}

	end := offset + limit
	if end > len(matchingCheckouts) {
		end = len(matchingCheckouts)
	}

	return matchingCheckouts[start:end], nil
}

// GetActiveCheckoutsByUserID retrieves all active checkouts for a user
func (r *MockCheckoutRepository) GetActiveCheckoutsByUserID(userID uint) ([]*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var activeCheckouts []*entity.Checkout

	for _, checkout := range r.checkouts {
		if checkout.UserID == userID && checkout.Status == entity.CheckoutStatusActive {
			activeCheckouts = append(activeCheckouts, checkout)
		}
	}

	return activeCheckouts, nil
}

// GetCompletedCheckoutsByUserID retrieves all completed checkouts for a user
func (r *MockCheckoutRepository) GetCompletedCheckoutsByUserID(userID uint, offset, limit int) ([]*entity.Checkout, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var completedCheckouts []*entity.Checkout

	for _, checkout := range r.checkouts {
		if checkout.UserID == userID && checkout.Status == entity.CheckoutStatusCompleted {
			completedCheckouts = append(completedCheckouts, checkout)
		}
	}

	// Apply offset and limit
	start := offset
	if start >= len(completedCheckouts) {
		return []*entity.Checkout{}, nil
	}

	end := offset + limit
	if end > len(completedCheckouts) {
		end = len(completedCheckouts)
	}

	return completedCheckouts[start:end], nil
}
