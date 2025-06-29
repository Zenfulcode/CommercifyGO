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
	defaultCurrency, err := currencyRepo.GetDefault()
	if err != nil {
		return nil
	}

	return &ProductUseCase{
		productRepo:        productRepo,
		categoryRepo:       categoryRepo,
		productVariantRepo: productVariantRepo,
		currencyRepo:       currencyRepo,
		orderRepo:          orderRepo,
		checkoutRepo:       checkoutRepo,
		defaultCurrency:    defaultCurrency,
	}
}

type VariantInput struct {
	SKU        string
	Stock      int
	Weight     float64
	Images     []string
	Attributes entity.VariantAttributes
	Prices     map[string]int64
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
	IsDefault bool
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
				weight,
				variantInput.Attributes,
				variantInput.Prices,
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
func (uc *ProductUseCase) GetProductByID(id uint, currencyCode string) (*entity.Product, error) {
	var currency *entity.Currency
	if currencyCode == "" {
		// Use default currency if none provided
		currency = uc.defaultCurrency
	} else {
		var err error
		currency, err = uc.currencyRepo.GetByCode(currencyCode)
		if err != nil {
			return nil, errors.New("invalid currency code: " + currencyCode)
		}
	}

	// First get the product with all its data
	product, err := uc.productRepo.GetByIDAndCurrency(id, currency.Code)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProductInput contains the data needed to update a product (prices in dollars)
type UpdateProductInput struct {
	Name        string
	Description string
	CategoryID  uint
	Images      []string
	Active      bool
}

// UpdateProduct updates a product
func (uc *ProductUseCase) UpdateProduct(id uint, input UpdateProductInput) (*entity.Product, error) {
	// Get product
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate category exists if changing
	if input.CategoryID != 0 && input.CategoryID != product.CategoryID {
		_, err := uc.categoryRepo.GetByID(input.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = input.CategoryID
	}

	updated := product.Update(input.Name, input.Description, input.Images, input.Active)
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
	IsDefault bool
}

// UpdateVariant updates a product variant
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
	variant.Update(
		input.SKU,
		input.Stock,
		input.Images,
		input.Attributes,
		input.Prices,
	)

	// Handle default status
	if input.IsDefault != variant.IsDefault {
		// If setting this variant as default, unset any other default variants
		if input.IsDefault {
			for _, v := range product.Variants {
				if v.ID != variantId && v.IsDefault {
					v.IsDefault = false
				}
			}
		}

		variant.IsDefault = input.IsDefault
	}

	// Update variant in repository
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	return variant, nil
}

// AddVariantInput contains the data needed to add a variant to a product
type AddVariantInput struct {
	VariantInput
	IsDefault bool
}

// AddVariant adds a new variant to a product
func (uc *ProductUseCase) AddVariant(productID uint, input AddVariantInput) (*entity.ProductVariant, error) {
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// Create variant
	variant, err := entity.NewProductVariant(
		input.SKU,
		input.Stock,
		input.Weight,
		input.Attributes,
		input.Prices,
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

// DeleteVariant deletes a product variant
func (uc *ProductUseCase) DeleteVariant(productID, variantID uint) error {
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	err = product.RemoveVariant(variantID)
	if err != nil {
		return fmt.Errorf("failed to remove variant: %w", err)
	}

	if err := uc.productRepo.Update(product); err != nil {
		return fmt.Errorf("failed to update product after removing variant: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product after checking it has no associated orders or active checkouts
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

// SetVariantPriceInput contains the data needed to set a price for a variant in a specific currency
type SetVariantPriceInput struct {
	VariantID    uint    `json:"variant_id"`
	CurrencyCode string  `json:"currency_code"`
	Price        float64 `json:"price"`
}

// SetVariantPriceInCurrency sets or updates the price for a variant in a specific currency
func (uc *ProductUseCase) SetVariantPriceInCurrency(input SetVariantPriceInput) (*entity.ProductVariant, error) {
	// Validate input
	if input.VariantID == 0 {
		return nil, errors.New("variant ID is required")
	}
	if input.CurrencyCode == "" {
		return nil, errors.New("currency code is required")
	}
	if input.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	// Get the variant
	variant, err := uc.productVariantRepo.GetByID(input.VariantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	// Validate currency exists and is enabled
	currency, err := uc.currencyRepo.GetByCode(input.CurrencyCode)
	if err != nil {
		return nil, fmt.Errorf("currency %s not found: %w", input.CurrencyCode, err)
	}
	if !currency.IsEnabled {
		return nil, fmt.Errorf("currency %s is not enabled", input.CurrencyCode)
	}

	// Set the price in the variant
	err = variant.SetPriceInCurrency(input.CurrencyCode, money.ToCents(input.Price))
	if err != nil {
		return nil, fmt.Errorf("failed to set price: %w", err)
	}

	// Update the variant in the repository
	err = uc.productVariantRepo.Update(variant)
	if err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	return variant, nil
}

// RemoveVariantPriceInCurrency removes the price for a variant in a specific currency
func (uc *ProductUseCase) RemoveVariantPriceInCurrency(variantID uint, currencyCode string) (*entity.ProductVariant, error) {
	// Validate input
	if variantID == 0 {
		return nil, errors.New("variant ID is required")
	}
	if currencyCode == "" {
		return nil, errors.New("currency code is required")
	}

	// Get the variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	// Remove the price
	err = variant.RemovePriceInCurrency(currencyCode)
	if err != nil {
		return nil, fmt.Errorf("failed to remove price: %w", err)
	}

	// Update the variant in the repository
	err = uc.productVariantRepo.Update(variant)
	if err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	return variant, nil
}

// GetVariantPrices returns all prices for a variant across all currencies
func (uc *ProductUseCase) GetVariantPrices(variantID uint) (map[string]float64, error) {
	// Validate input
	if variantID == 0 {
		return nil, errors.New("variant ID is required")
	}

	// Get the variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	// Get all prices in cents
	pricesInCents := variant.GetAllPrices()

	// Convert to float64 (dollars/euros/etc.)
	prices := make(map[string]float64)
	for currency, priceInCents := range pricesInCents {
		prices[currency] = money.FromCents(priceInCents)
	}

	return prices, nil
}

// SetMultipleVariantPricesInput contains the data needed to set multiple prices for a variant
type SetMultipleVariantPricesInput struct {
	VariantID uint               `json:"variant_id"`
	Prices    map[string]float64 `json:"prices"` // currency_code -> price
}

// SetMultipleVariantPrices sets multiple prices for a variant at once
func (uc *ProductUseCase) SetMultipleVariantPrices(input SetMultipleVariantPricesInput) (*entity.ProductVariant, error) {
	// Validate input
	if input.VariantID == 0 {
		return nil, errors.New("variant ID is required")
	}
	if len(input.Prices) == 0 {
		return nil, errors.New("at least one price must be provided")
	}

	// Get the variant
	variant, err := uc.productVariantRepo.GetByID(input.VariantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	// Validate all currencies and prices
	for currencyCode, price := range input.Prices {
		if currencyCode == "" {
			return nil, errors.New("currency code cannot be empty")
		}
		if price <= 0 {
			return nil, fmt.Errorf("price for %s must be greater than zero", currencyCode)
		}

		// Validate currency exists and is enabled
		currency, err := uc.currencyRepo.GetByCode(currencyCode)
		if err != nil {
			return nil, fmt.Errorf("currency %s not found: %w", currencyCode, err)
		}
		if !currency.IsEnabled {
			return nil, fmt.Errorf("currency %s is not enabled", currencyCode)
		}
	}

	// Set all prices
	for currencyCode, price := range input.Prices {
		err = variant.SetPriceInCurrency(currencyCode, money.ToCents(price))
		if err != nil {
			return nil, fmt.Errorf("failed to set price for %s: %w", currencyCode, err)
		}
	}

	// Update the variant in the repository
	err = uc.productVariantRepo.Update(variant)
	if err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	return variant, nil
}
