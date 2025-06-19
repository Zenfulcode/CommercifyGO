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
	SKU        string
	Price      float64
	Stock      int
	Attributes []entity.VariantAttribute
	Images     []string
	IsDefault  bool
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

	// Create product
	product, err := entity.NewProduct(
		input.Name,
		input.Description,
		input.Currency,
		input.CategoryID,
		input.Images,
	)
	if err != nil {
		return nil, err
	}

	// Save product
	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	// If product has variants, create them
	if len(input.Variants) > 0 {
		variants := make([]*entity.ProductVariant, 0, len(input.Variants))
		for _, variantInput := range input.Variants {

			variant, err := entity.NewProductVariant(
				product.ID,
				variantInput.SKU,
				variantInput.Price,
				product.CurrencyCode,
				variantInput.Stock,
				variantInput.Attributes,
				variantInput.Images,
				variantInput.IsDefault,
			)
			if err != nil {
				return nil, err
			}

			variants = append(variants, variant)
		}

		// Save each variant individually to process their currency prices too
		for _, variant := range variants {
			if err := uc.productVariantRepo.Create(variant); err != nil {
				return nil, err
			}
		}

		// Add variants to product
		product.Variants = variants
		// Only set has_variants=true if there are multiple variants
		product.HasVariants = len(variants) > 1
		product.Active = input.Active
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (uc *ProductUseCase) GetProductByID(id uint, currencyCode string) (*entity.Product, error) {
	if currencyCode == "" {
		return nil, errors.New("currency code is required")
	}

	// First get the product with all its data
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	product.Variants, err = uc.productVariantRepo.GetByProduct(id)
	if err != nil {
		return nil, err
	}

	// Validate currency exists
	currency, err := uc.currencyRepo.GetByCode(currencyCode)
	if err != nil {
		return nil, errors.New("invalid currency code: " + currencyCode)
	}

	currencyPrice, found := product.GetPriceInCurrency(currency.Code)
	if found {
		product.Price = currencyPrice
	} else {
		product.Price = uc.defaultCurrency.ConvertAmount(currencyPrice, currency)
	}

	product.CurrencyCode = currency.Code

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
	product, err := uc.productRepo.GetByIDWithVariants(id)
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

	// Update product fields
	if input.Name != "" {
		product.Name = input.Name
	}
	if input.Description != "" {
		product.Description = input.Description
	}

	if len(input.Images) > 0 {
		product.Images = input.Images
	}
	if input.Active != product.Active {
		product.Active = input.Active
	}

	// Update product in repository
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateVariantInput contains the data needed to update a product variant (prices in dollars)
type UpdateVariantInput struct {
	SKU        string
	Price      float64
	Stock      int
	Attributes []entity.VariantAttribute
	Images     []string
	IsDefault  bool
}

// UpdateVariant updates a product variant
func (uc *ProductUseCase) UpdateVariant(productID uint, variantID uint, input UpdateVariantInput) (*entity.ProductVariant, error) {
	// Get variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return nil, err
	}

	// Check if variant belongs to the product
	if variant.ProductID != productID {
		return nil, errors.New("variant does not belong to this product")
	}

	// Update variant fields
	if input.SKU != "" {
		variant.SKU = input.SKU
	}
	if input.Price > 0 {
		variant.Price = money.ToCents(input.Price) // Convert to cents
	}
	if input.Stock >= 0 {
		variant.Stock = input.Stock
	}
	if len(input.Attributes) > 0 {
		variant.Attributes = input.Attributes
	}
	if len(input.Images) > 0 {
		variant.Images = input.Images
	}

	// Handle default status
	if input.IsDefault != variant.IsDefault {
		// If setting this variant as default, unset any other default variants
		if input.IsDefault {
			variants, err := uc.productVariantRepo.GetByProduct(productID)
			if err != nil {
				return nil, err
			}

			for _, v := range variants {
				if v.ID != variantID && v.IsDefault {
					v.IsDefault = false
					if err := uc.productVariantRepo.Update(v); err != nil {
						return nil, err
					}
				}
			}
		}

		variant.IsDefault = input.IsDefault
	}

	// Update variant in repository
	if err := uc.productVariantRepo.Update(variant); err != nil {
		return nil, err
	}

	return variant, nil
}

// AddVariantInput contains the data needed to add a variant to a product
type AddVariantInput struct {
	ProductID  uint
	SKU        string
	Price      float64
	Stock      int
	Attributes []entity.VariantAttribute
	Images     []string
	IsDefault  bool
}

// AddVariant adds a new variant to a product
func (uc *ProductUseCase) AddVariant(input AddVariantInput) (*entity.ProductVariant, error) {
	product, err := uc.productRepo.GetByIDWithVariants(input.ProductID)
	if err != nil {
		return nil, err
	}

	// Create variant
	variant, err := entity.NewProductVariant(
		input.ProductID,
		input.SKU,
		input.Price, // Use cents
		product.CurrencyCode,
		input.Stock,
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
				if err := uc.productVariantRepo.Update(v); err != nil {
					return nil, err
				}
			}
		}
	}

	// Save variant
	if err := uc.productVariantRepo.Create(variant); err != nil {
		return nil, err
	}

	return variant, nil
}

// DeleteVariant deletes a product variant
func (uc *ProductUseCase) DeleteVariant(productID uint, variantID uint) error {
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return err
	}

	// Check if variant belongs to the product
	if variant.ProductID != productID {
		return errors.New("variant does not belong to this product")
	}

	// Delete variant
	return uc.productVariantRepo.Delete(variantID)
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
