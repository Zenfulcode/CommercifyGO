package entity

import (
	"errors"
	"slices"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/money"
	"gorm.io/gorm"
)

// DiscountType represents the type of discount
type DiscountType string

const (
	// DiscountTypeBasket applies to the entire order
	DiscountTypeBasket DiscountType = "basket"
	// DiscountTypeProduct applies to specific products
	DiscountTypeProduct DiscountType = "product"
)

// DiscountMethod represents how the discount is calculated
type DiscountMethod string

const (
	// DiscountMethodFixed is a fixed amount discount
	DiscountMethodFixed DiscountMethod = "fixed"
	// DiscountMethodPercentage is a percentage discount
	DiscountMethodPercentage DiscountMethod = "percentage"
)

// Discount represents a discount in the system
type Discount struct {
	gorm.Model
	Code             string         `gorm:"uniqueIndex;not null;size:100"`
	Type             DiscountType   `gorm:"not null;size:50"`
	Method           DiscountMethod `gorm:"not null;size:50"`
	Value            float64        `gorm:"not null"`
	MinOrderValue    int64          `gorm:"default:0"`
	MaxDiscountValue int64          `gorm:"default:0"`
	ProductIDs       []uint         `gorm:"type:jsonb"`
	CategoryIDs      []uint         `gorm:"type:jsonb"`
	StartDate        time.Time      `gorm:"index"`
	EndDate          time.Time      `gorm:"index"`
	UsageLimit       int            `gorm:"default:0"`
	CurrentUsage     int            `gorm:"default:0"`
	Active           bool           `gorm:"default:true"`
}

// NewDiscount creates a new discount
func NewDiscount(
	code string,
	discountType DiscountType,
	method DiscountMethod,
	value float64,
	minOrderValue int64,
	maxDiscountValue int64,
	productIDs []uint,
	categoryIDs []uint,
	startDate time.Time,
	endDate time.Time,
	usageLimit int,
) (*Discount, error) {
	if code == "" {
		return nil, errors.New("discount code cannot be empty")
	}

	if value <= 0 {
		return nil, errors.New("discount value must be greater than zero")
	}

	if method == DiscountMethodPercentage && value > 100 {
		return nil, errors.New("percentage discount cannot exceed 100%")
	}

	if discountType == DiscountTypeProduct && len(productIDs) == 0 && len(categoryIDs) == 0 {
		return nil, errors.New("product discount must specify at least one product or category")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	return &Discount{
		Code:             code,
		Type:             discountType,
		Method:           method,
		Value:            value,
		MinOrderValue:    minOrderValue,
		MaxDiscountValue: maxDiscountValue,
		ProductIDs:       productIDs,
		CategoryIDs:      categoryIDs,
		StartDate:        startDate,
		EndDate:          endDate,
		UsageLimit:       usageLimit,
		CurrentUsage:     0,
		Active:           true,
	}, nil
}

// IsValid checks if the discount is valid for the current time and usage
func (d *Discount) IsValid() bool {
	now := time.Now().Local()
	return d.Active &&
		now.After(d.StartDate.Local()) &&
		now.Before(d.EndDate.Local()) &&
		(d.UsageLimit == 0 || d.CurrentUsage < d.UsageLimit)
}

// IsApplicableToOrder checks if the discount is applicable to the given order
func (d *Discount) IsApplicableToOrder(order *Order) bool {
	if !d.IsValid() {
		return false
	}

	// Check minimum order value
	if d.MinOrderValue > 0 && order.TotalAmount < d.MinOrderValue {
		return false
	}

	switch d.Type {
	case DiscountTypeBasket:
		return true
	case DiscountTypeProduct:
		for _, item := range order.Items {
			// Check if the product is directly included
			if slices.Contains(d.ProductIDs, item.ProductID) {
				return true
			}
			// Note: Category check is handled separately in the CalculateDiscount method
			// since we need product details from the repository
		}
		// If we have category IDs but no direct product matches,
		// we still need to check if any product belongs to those categories
		// This is handled in the use case layer
		if len(d.CategoryIDs) > 0 {
			return true
		}
		return false
	}

	return false
}

// CalculateDiscount calculates the discount amount for an order
func (d *Discount) CalculateDiscount(order *Order) int64 {
	if !d.IsApplicableToOrder(order) {
		return 0
	}

	var discountAmount int64

	switch d.Type {
	case DiscountTypeBasket:
		// Calculate discount for the entire order
		switch d.Method {
		case DiscountMethodFixed:
			// For fixed amount method, the value is in dollars and needs to be converted to cents
			// But since we updated the structure, the database will provide the value already in cents
			discountAmount = money.ToCents(d.Value)
		case DiscountMethodPercentage:
			// For percentage, apply the percentage to the total amount
			discountAmount = money.ApplyPercentage(order.TotalAmount, d.Value)
		}
	case DiscountTypeProduct:
		// Calculate discount for eligible products only
		for _, item := range order.Items {
			isEligible := slices.Contains(d.ProductIDs, item.ProductID)

			if isEligible {
				itemTotal := item.Subtotal
				switch d.Method {
				case DiscountMethodFixed:
					// For fixed discount, apply once per item (not per quantity)
					// This matches with the current implementation in ApplyDiscountToOrder
					fixedDiscountInCents := money.ToCents(d.Value)
					itemDiscount := min(fixedDiscountInCents, itemTotal)
					discountAmount += itemDiscount
				case DiscountMethodPercentage:
					// For percentage discount, apply percentage to item total
					discountAmount += money.ApplyPercentage(itemTotal, d.Value)
				}
			}
		}
	}

	// Apply maximum discount cap if specified
	if d.MaxDiscountValue > 0 && discountAmount > d.MaxDiscountValue {
		discountAmount = d.MaxDiscountValue
	}

	// Ensure discount doesn't exceed order total
	if discountAmount > order.TotalAmount {
		discountAmount = order.TotalAmount
	}

	return discountAmount
}

// IncrementUsage increments the usage count of the discount
func (d *Discount) IncrementUsage() {
	d.CurrentUsage++

}
