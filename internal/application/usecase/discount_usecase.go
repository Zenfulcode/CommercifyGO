package usecase

import (
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// DiscountUseCase implements discount-related use cases
type DiscountUseCase struct {
	discountRepo repository.DiscountRepository
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	orderRepo    repository.OrderRepository
}

// NewDiscountUseCase creates a new DiscountUseCase
func NewDiscountUseCase(
	discountRepo repository.DiscountRepository,
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	orderRepo repository.OrderRepository,
) *DiscountUseCase {
	return &DiscountUseCase{
		discountRepo: discountRepo,
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		orderRepo:    orderRepo,
	}
}

// CreateDiscountInput contains the data needed to create a discount
type CreateDiscountInput struct {
	Code             string    `json:"code"`
	Type             string    `json:"type"`
	Method           string    `json:"method"`
	Value            float64   `json:"value"`
	MinOrderValue    float64   `json:"min_order_value"`
	MaxDiscountValue float64   `json:"max_discount_value"`
	ProductIDs       []uint    `json:"product_ids"`
	CategoryIDs      []uint    `json:"category_ids"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	UsageLimit       int       `json:"usage_limit"`
}

// CreateDiscount creates a new discount
func (uc *DiscountUseCase) CreateDiscount(input CreateDiscountInput) (*entity.Discount, error) {
	// Validate discount type
	var discountType entity.DiscountType
	switch input.Type {
	case string(entity.DiscountTypeBasket):
		discountType = entity.DiscountTypeBasket
	case string(entity.DiscountTypeProduct):
		discountType = entity.DiscountTypeProduct
	default:
		return nil, errors.New("invalid discount type")
	}

	// Validate discount method
	var discountMethod entity.DiscountMethod
	switch input.Method {
	case string(entity.DiscountMethodFixed):
		discountMethod = entity.DiscountMethodFixed
	case string(entity.DiscountMethodPercentage):
		discountMethod = entity.DiscountMethodPercentage
	default:
		return nil, errors.New("invalid discount method")
	}

	// Check if discount code already exists
	existingDiscount, err := uc.discountRepo.GetByCode(input.Code)
	if err == nil && existingDiscount != nil {
		return nil, errors.New("discount code already exists")
	}

	// Validate product IDs if it's a product discount
	if discountType == entity.DiscountTypeProduct && len(input.ProductIDs) > 0 {
		for _, productID := range input.ProductIDs {
			_, err := uc.productRepo.GetByID(productID)
			if err != nil {
				return nil, errors.New("invalid product ID: " + err.Error())
			}
		}
	}

	// Validate category IDs if it's a product discount
	if discountType == entity.DiscountTypeProduct && len(input.CategoryIDs) > 0 {
		for _, categoryID := range input.CategoryIDs {
			_, err := uc.categoryRepo.GetByID(categoryID)
			if err != nil {
				return nil, errors.New("invalid category ID: " + err.Error())
			}
		}
	}

	// Create discount
	discount, err := entity.NewDiscount(
		input.Code,
		discountType,
		discountMethod,
		input.Value,
		money.ToCents(input.MinOrderValue),
		money.ToCents(input.MaxDiscountValue),
		input.ProductIDs,
		input.CategoryIDs,
		input.StartDate,
		input.EndDate,
		input.UsageLimit,
	)
	if err != nil {
		return nil, err
	}

	// Save discount
	if err := uc.discountRepo.Create(discount); err != nil {
		return nil, err
	}

	return discount, nil
}

// GetDiscountByID retrieves a discount by ID
func (uc *DiscountUseCase) GetDiscountByID(id uint) (*entity.Discount, error) {
	return uc.discountRepo.GetByID(id)
}

// GetDiscountByCode retrieves a discount by code
func (uc *DiscountUseCase) GetDiscountByCode(code string) (*entity.Discount, error) {
	return uc.discountRepo.GetByCode(code)
}

// UpdateDiscountInput contains the data needed to update a discount
type UpdateDiscountInput struct {
	Code             string    `json:"code"`
	Type             string    `json:"type"`
	Method           string    `json:"method"`
	Value            float64   `json:"value"`
	MinOrderValue    float64   `json:"min_order_value"`
	MaxDiscountValue float64   `json:"max_discount_value"`
	ProductIDs       []uint    `json:"product_ids"`
	CategoryIDs      []uint    `json:"category_ids"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	UsageLimit       int       `json:"usage_limit"`
	Active           bool      `json:"active"`
}

// UpdateDiscount updates a discount
func (uc *DiscountUseCase) UpdateDiscount(id uint, input UpdateDiscountInput) (*entity.Discount, error) {
	// Get existing discount
	discount, err := uc.discountRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate discount type
	if input.Type != "" {
		switch input.Type {
		case string(entity.DiscountTypeBasket):
			discount.Type = entity.DiscountTypeBasket
		case string(entity.DiscountTypeProduct):
			discount.Type = entity.DiscountTypeProduct
		default:
			return nil, errors.New("invalid discount type")
		}
	}

	// Validate discount method
	if input.Method != "" {
		switch input.Method {
		case string(entity.DiscountMethodFixed):
			discount.Method = entity.DiscountMethodFixed
		case string(entity.DiscountMethodPercentage):
			discount.Method = entity.DiscountMethodPercentage
		default:
			return nil, errors.New("invalid discount method")
		}
	}

	// Update fields
	if input.Code != "" && input.Code != discount.Code {
		// Check if new code already exists
		existingDiscount, err := uc.discountRepo.GetByCode(input.Code)
		if err == nil && existingDiscount != nil && existingDiscount.ID != id {
			return nil, errors.New("discount code already exists")
		}
		discount.Code = input.Code
	}

	if input.Value > 0 {
		discount.Value = input.Value
	}

	if input.MinOrderValue >= 0 {
		discount.MinOrderValue = money.ToCents(input.MinOrderValue)
	}

	if input.MaxDiscountValue >= 0 {
		discount.MaxDiscountValue = money.ToCents(input.MaxDiscountValue)
	}

	if len(input.ProductIDs) > 0 {
		// Validate product IDs
		for _, productID := range input.ProductIDs {
			_, err := uc.productRepo.GetByID(productID)
			if err != nil {
				return nil, errors.New("invalid product ID: " + err.Error())
			}
		}
		discount.ProductIDs = input.ProductIDs
	}

	if len(input.CategoryIDs) > 0 {
		// Validate category IDs
		for _, categoryID := range input.CategoryIDs {
			_, err := uc.categoryRepo.GetByID(categoryID)
			if err != nil {
				return nil, errors.New("invalid category ID: " + err.Error())
			}
		}
		discount.CategoryIDs = input.CategoryIDs
	}

	if !input.StartDate.IsZero() {
		discount.StartDate = input.StartDate
	}

	if !input.EndDate.IsZero() {
		discount.EndDate = input.EndDate
	}

	if input.UsageLimit >= 0 {
		discount.UsageLimit = input.UsageLimit
	}

	discount.Active = input.Active
	discount.UpdatedAt = time.Now()

	// Save discount
	if err := uc.discountRepo.Update(discount); err != nil {
		return nil, err
	}

	return discount, nil
}

// DeleteDiscount deletes a discount
func (uc *DiscountUseCase) DeleteDiscount(id uint) error {
	// check if orders are using this discount
	inUse, err := uc.orderRepo.IsDiscountIdUsed(id)
	if err != nil {
		return err
	}

	if inUse {
		return errors.New("discount is in use by an order")
	}

	return uc.discountRepo.Delete(id)
}

// ListDiscounts lists all discounts with pagination
func (uc *DiscountUseCase) ListDiscounts(offset, limit int) ([]*entity.Discount, error) {
	return uc.discountRepo.List(offset, limit)
}

// ListActiveDiscounts lists all active discounts with pagination
func (uc *DiscountUseCase) ListActiveDiscounts(offset, limit int) ([]*entity.Discount, error) {
	return uc.discountRepo.ListActive(offset, limit)
}

// ApplyDiscountToOrderInput contains the data needed to apply a discount to an order
type ApplyDiscountToOrderInput struct {
	OrderID      uint   `json:"order_id"`
	DiscountCode string `json:"discount_code"`
}

// ApplyDiscountToOrder applies a discount to an order
func (uc *DiscountUseCase) ApplyDiscountToOrder(input ApplyDiscountToOrderInput, order *entity.Order) (*entity.Order, error) {
	// Get discount by code
	discount, err := uc.discountRepo.GetByCode(input.DiscountCode)
	if err != nil {
		return nil, errors.New("invalid discount code")
	}

	// For category-based discounts, we need to modify the ProductIDs to include products from those categories
	if discount.Type == entity.DiscountTypeProduct && len(discount.CategoryIDs) > 0 {
		// Create a map to track which items in the order need discounts
		eligibleProducts := make(map[uint]bool)

		// First add directly specified products
		for _, productID := range discount.ProductIDs {
			eligibleProducts[productID] = true
		}

		// Then find products that belong to the specified categories
		for _, categoryID := range discount.CategoryIDs {
			// Get all products in this category
			products, err := uc.productRepo.List("", "", categoryID, 0, 1000, 0, 0, true)
			if err == nil && len(products) > 0 {
				// Add these products to our eligibility map
				for _, product := range products {
					eligibleProducts[product.ID] = true
				}
			}
		}

		// Now calculate the discount based on eligible products
		var discountAmount int64

		for _, item := range order.Items {
			if eligibleProducts[item.ProductID] {
				itemTotal := int64(item.Quantity) * item.Price

				switch discount.Method {
				case entity.DiscountMethodFixed:
					// Apply fixed discount per item
					itemDiscount := min(money.ToCents(discount.Value), itemTotal)
					discountAmount += itemDiscount
				case entity.DiscountMethodPercentage:
					// Apply percentage discount to the item
					// itemTotal * (discount.Value / 100)
					itemDiscount := money.ApplyPercentage(itemTotal, discount.Value)
					discountAmount += itemDiscount
				}
			}
		}

		// Apply maximum discount cap if specified
		if discount.MaxDiscountValue > 0 && discountAmount > discount.MaxDiscountValue {
			discountAmount = discount.MaxDiscountValue
		}

		// Ensure discount doesn't exceed order total
		if discountAmount > order.TotalAmount {
			discountAmount = order.TotalAmount
		}

		// Apply calculated discount amount
		order.DiscountAmount = discountAmount
		order.FinalAmount = order.TotalAmount - discountAmount

		// Record the applied discount using JSON helper method
		appliedDiscount := &entity.AppliedDiscount{
			DiscountID:     discount.ID,
			DiscountCode:   discount.Code,
			DiscountAmount: discountAmount,
		}
		order.SetAppliedDiscountJSON(appliedDiscount)

		order.UpdatedAt = time.Now()
	} else {
		// For non-category specific discounts, use the standard entity method
		if err := order.ApplyDiscount(discount); err != nil {
			return nil, err
		}
	}

	uc.orderRepo.Update(order)

	// Increment discount usage
	if err := uc.discountRepo.IncrementUsage(discount.ID); err != nil {
		return nil, err
	}

	return order, nil
}

// RemoveDiscountFromOrder removes a discount from an order
func (uc *DiscountUseCase) RemoveDiscountFromOrder(order *entity.Order) {
	order.RemoveDiscount()
	uc.orderRepo.Update(order)
}
