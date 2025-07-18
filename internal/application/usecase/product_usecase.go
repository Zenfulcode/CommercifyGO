package usecase

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductUseCase implements product-related use cases
type ProductUseCase struct {
	productRepo        repository.ProductRepository
	categoryRepo       repository.CategoryRepository
	productVariantRepo repository.ProductVariantRepository
	currencyRepo       repository.CurrencyRepository
	orderRepo          repository.OrderRepository
	checkoutRepo       repository.CheckoutRepository
	defaultCurrency    *entity.Currency
}

// NewProductUseCase creates a new ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	productVariantRepo repository.ProductVariantRepository,
	currencyRepo repository.CurrencyRepository,
	orderRepo repository.OrderRepository,
	checkoutRepo repository.CheckoutRepository,
) *ProductUseCase {
	uc := &ProductUseCase{
		productRepo:        productRepo,
		categoryRepo:       categoryRepo,
		productVariantRepo: productVariantRepo,
		currencyRepo:       currencyRepo,
		orderRepo:          orderRepo,
		checkoutRepo:       checkoutRepo,
	}

	// Try to get default currency but don't fail if it doesn't exist
	defaultCurrency, err := currencyRepo.GetDefault()
	if err == nil {
		uc.defaultCurrency = defaultCurrency
	}
	// If no default currency exists, defaultCurrency will be nil
	// This should be handled in methods that need it

	return uc
}

type VariantInput struct {
	SKU        string
	Stock      int
	Weight     float64
	Images     []string
	Attributes entity.VariantAttributes
	Price      int64
	IsDefault  bool
}

// CreateProductInput contains the data needed to create a product
type CreateProductInput struct {
	Name        string
	Description string
	Currency    string
	CategoryID  uint
	Images      []string
	Variants    []CreateVariantInput
	Active      bool
}

// CreateVariantInput contains the data needed to create a product variant
type CreateVariantInput struct {
	VariantInput
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(input CreateProductInput) (*entity.Product, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Validate currency exists
	_, err = uc.currencyRepo.GetByCode(input.Currency)
	if err != nil {
		return nil, errors.New("invalid currency code: " + input.Currency)
	}

	variants := make([]*entity.ProductVariant, 0, len(input.Variants))

	// If product has variants, create them
	if len(input.Variants) > 0 {
		defaultVariantCount := 0

		// First pass: count default variants and validate there's only one
		for _, variantInput := range input.Variants {
			if variantInput.IsDefault {
				defaultVariantCount++
			}
		}

		// Ensure only one default variant
		if defaultVariantCount > 1 {
			return nil, errors.New("only one variant can be set as default")
		}

		// If no default variant specified, set the first one as default
		if defaultVariantCount == 0 && len(input.Variants) > 0 {
			input.Variants[0].IsDefault = true
		}

		for _, variantInput := range input.Variants {
			// Create variant with new schema - weight defaults to 0 if not provided
			weight := variantInput.Weight
			if weight == 0 {
				weight = 0.0 // default weight
			}

			variant, err := entity.NewProductVariant(
				variantInput.SKU,
				variantInput.Stock,
				variantInput.Price,
				weight,
				variantInput.Attributes,
				variantInput.Images,
				variantInput.IsDefault,
			)
			if err != nil {
				return nil, err
			}

			variants = append(variants, variant)
		}
	}

	// Create product
	product, err := entity.NewProduct(
		input.Name,
		input.Description,
		input.Currency,
		input.CategoryID,
		input.Images,
		variants,
		input.Active,
	)
	if err != nil {
		return nil, err
	}

	// Save product
	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (uc *ProductUseCase) GetProductByID(id uint) (*entity.Product, error) {
	// Simply get the product with all its data - no currency filtering needed
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProductInput contains the data needed to update a product (prices in dollars)
type UpdateProductInput struct {
	Name        *string
	Description *string
	Currency    *string
	CategoryID  *uint
	Images      *[]string
	Active      *bool
	Variants    *[]UpdateVariantInput
}

// UpdateProduct updates a product (admin only)
func (uc *ProductUseCase) UpdateProduct(id uint, input UpdateProductInput) (*entity.Product, error) {
	// Get product
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate category exists if changing
	if input.CategoryID != nil && *input.CategoryID != product.CategoryID {
		_, err := uc.categoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
	}

	// Validate currency exists if changing
	if input.Currency != nil && *input.Currency != product.Currency {
		_, err := uc.currencyRepo.GetByCode(*input.Currency)
		if err != nil {
			return nil, errors.New("invalid currency code: " + *input.Currency)
		}
	}

	// Update basic product fields
	updated := product.Update(input.Name, input.Description, input.Currency, input.Images, input.Active, input.CategoryID)

	// Handle variant updates if provided
	if input.Variants != nil {
		for _, variantUpdate := range *input.Variants {
			// Find the variant to update by SKU or ID
			var targetVariant *entity.ProductVariant
			for _, variant := range product.Variants {
				// If SKU is provided, match by SKU; otherwise this is a new variant
				if variantUpdate.SKU != "" && variant.SKU == variantUpdate.SKU {
					targetVariant = variant
					break
				}
			}

			if targetVariant != nil {
				// Handle IsDefault logic before updating
				var isDefaultPtr *bool
				if variantUpdate.IsDefault {
					// If setting this variant as default, unset any other default variants
					for _, v := range product.Variants {
						if v.ID != targetVariant.ID && v.IsDefault {
							v.IsDefault = false
						}
					}
					isDefaultPtr = &variantUpdate.IsDefault
				} else {
					isDefaultPtr = &variantUpdate.IsDefault
				}

				// Update existing variant
				variantUpdated, err := targetVariant.Update(
					variantUpdate.SKU,
					variantUpdate.Stock,
					variantUpdate.Price,
					variantUpdate.Weight,
					variantUpdate.Images,
					variantUpdate.Attributes,
					isDefaultPtr,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to update variant: %w", err)
				}
				if variantUpdated {
					updated = true
				}
			} else {
				// Add new variant if SKU is provided and not found
				if variantUpdate.SKU != "" {
					// If this new variant is set as default, unset any existing default variants
					if variantUpdate.IsDefault {
						for _, v := range product.Variants {
							if v.IsDefault {
								v.IsDefault = false
							}
						}
					}

					newVariant, err := entity.NewProductVariant(
						variantUpdate.SKU,
						variantUpdate.Stock,
						variantUpdate.Price,
						variantUpdate.Weight,
						variantUpdate.Attributes,
						variantUpdate.Images,
						variantUpdate.IsDefault,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to create variant: %w", err)
					}

					err = product.AddVariant(newVariant)
					if err != nil {
						return nil, fmt.Errorf("failed to add variant to product: %w", err)
					}
					updated = true
				}
			}
		}
	}

	if !updated {
		return product, nil // No changes to update
	}

	// Update product in repository
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateVariantInput contains the data needed to update a product variant (prices in dollars)
type UpdateVariantInput struct {
	VariantInput
}

// UpdateVariant updates a product variant (admin only)
func (uc *ProductUseCase) UpdateVariant(productId, variantId uint, input UpdateVariantInput) (*entity.ProductVariant, error) {
	product, err := uc.productRepo.GetByID(productId)
	if err != nil {
		return nil, err
	}

	// Get the variant by SKU
	variant := product.GetVariantByID(variantId)
	if variant == nil {
		return nil, errors.New("variant not found")
	}

	// Update variant fields
	isDefaultPtr := &input.IsDefault
	updated, err := variant.Update(
		input.SKU,
		input.Stock,
		input.Price,
		input.Weight,
		input.Images,
		input.Attributes,
		isDefaultPtr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	// Handle default status if changed
	if updated && input.IsDefault != variant.IsDefault {
		// If setting this variant as default, unset any other default variants
		if input.IsDefault {
			for _, v := range product.Variants {
				if v.ID != variantId && v.IsDefault {
					v.IsDefault = false
				}
			}
		}
	}

	// Update variant in repository
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	return variant, nil
}

// AddVariant adds a new variant to a product (admin only)
func (uc *ProductUseCase) AddVariant(productID uint, input CreateVariantInput) (*entity.ProductVariant, error) {
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// Create variant
	variant, err := entity.NewProductVariant(
		input.SKU,
		input.Stock,
		input.Price,
		input.Weight,
		input.Attributes,
		input.Images,
		input.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	err = product.AddVariant(variant)
	if err != nil {
		return nil, err
	}

	if input.IsDefault {
		variants := product.Variants

		for _, v := range variants {
			if v.ID != variant.ID && v.IsDefault {
				v.IsDefault = false
			}
		}
	}

	// Update the product to persist the recalculated stock
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	return variant, nil
}

// DeleteVariant deletes a product variant (admin only)
func (uc *ProductUseCase) DeleteVariant(productID, variantID uint) error {
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	// Check if the variant exists in the product
	variant := product.GetVariantByID(variantID)
	if variant == nil {
		return errors.New("variant not found")
	}

	// Check if this is the last variant (products must have at least one variant)
	if len(product.Variants) <= 1 {
		return errors.New("cannot delete the last variant of a product")
	}

	// TODO: Add checks for orders and checkouts with this specific variant
	// For now, we'll check at the product level which is safer
	hasOrders, err := uc.orderRepo.HasOrdersWithProduct(productID)
	if err != nil {
		return fmt.Errorf("failed to check for product orders: %w", err)
	}
	if hasOrders {
		return errors.New("cannot delete variant from product with existing orders")
	}

	hasActiveCheckouts, err := uc.checkoutRepo.HasActiveCheckoutsWithProduct(productID)
	if err != nil {
		return fmt.Errorf("failed to check for active checkouts: %w", err)
	}
	if hasActiveCheckouts {
		return errors.New("cannot delete variant from product with active checkouts")
	}

	// Remove the variant from the product's variants slice first
	err = product.RemoveVariant(variantID)
	if err != nil {
		return fmt.Errorf("failed to remove variant from product: %w", err)
	}

	// Actually delete the variant from the database
	err = uc.productVariantRepo.Delete(variantID)
	if err != nil {
		return fmt.Errorf("failed to delete variant from database: %w", err)
	}

	// Update the product to reflect the change (this will sync the variants relationship)
	if err := uc.productRepo.Update(product); err != nil {
		return fmt.Errorf("failed to update product after removing variant: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product after checking it has no associated orders or active checkouts (admin only)
func (uc *ProductUseCase) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("product ID is required")
	}

	// Check if product has any associated orders
	hasOrders, err := uc.orderRepo.HasOrdersWithProduct(id)
	if err != nil {
		return fmt.Errorf("failed to check for product orders: %w", err)
	}

	if hasOrders {
		return errors.New("cannot delete product that has existing orders")
	}

	// Check if product has any active checkouts
	hasActiveCheckouts, err := uc.checkoutRepo.HasActiveCheckoutsWithProduct(id)
	if err != nil {
		return fmt.Errorf("failed to check for active checkouts: %w", err)
	}

	if hasActiveCheckouts {
		return errors.New("cannot delete product that is in active checkouts. Please wait for checkouts to complete or expire")
	}

	return uc.productRepo.Delete(id)
}

// SearchProductsInput contains the data needed to search for products (prices in dollars)
type SearchProductsInput struct {
	Query        string  `json:"query"`
	CurrencyCode string  `json:"currency_code"` // Optional currency code for prices
	MaxPrice     float64 `json:"max_price"`     // Price in dollars
	MinPrice     float64 `json:"min_price"`     // Price in dollars
	CategoryID   uint    `json:"category_id"`
	Offset       uint    `json:"offset"`
	Limit        uint    `json:"limit"`
	ActiveOnly   bool    `json:"active_only"` // Whether to filter active products only
}

// ListProducts lists all products with pagination and returns total count
func (uc *ProductUseCase) ListProducts(input SearchProductsInput) ([]*entity.Product, int, error) {
	minPriceCents := money.ToCents(input.MinPrice)
	maxPriceCents := money.ToCents(input.MaxPrice)

	products, err := uc.productRepo.List(
		input.Query,
		input.CurrencyCode,
		input.CategoryID,
		input.Offset,
		input.Limit,
		minPriceCents, // Convert to cents
		maxPriceCents, // Convert to cents
		input.ActiveOnly,
	)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.productRepo.Count(
		input.Query,
		input.CurrencyCode,
		input.CategoryID,
		minPriceCents, // Pass cents
		maxPriceCents, // Pass cents
		input.ActiveOnly,
	)
	if err != nil {
		return products, 0, err
	}

	return products, total, nil
}

// ListCategories lists all product categories
func (uc *ProductUseCase) ListCategories() ([]*entity.Category, error) {
	return uc.categoryRepo.List()
}
