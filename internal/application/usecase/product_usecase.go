package usecase

import (
	"errors"

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
	defaultCurrency    *entity.Currency
}

// NewProductUseCase creates a new ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	productVariantRepo repository.ProductVariantRepository,
	currencyRepo repository.CurrencyRepository,
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
		defaultCurrency:    defaultCurrency,
	}
}

// CurrencyPriceInput represents a price in a specific currency
type CurrencyPriceInput struct {
	CurrencyCode string
	Price        float64
}

// CreateProductInput contains the data needed to create a product (prices in dollars)
type CreateProductInput struct {
	Name           string
	Description    string
	Price          float64
	Stock          int
	Weight         float64
	CategoryID     uint
	Images         []string
	Variants       []CreateVariantInput
	CurrencyPrices []CurrencyPriceInput
}

// CreateVariantInput contains the data needed to create a product variant
type CreateVariantInput struct {
	SKU            string
	Price          float64
	Stock          int
	Attributes     []entity.VariantAttribute
	Images         []string
	IsDefault      bool
	CurrencyPrices []CurrencyPriceInput
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(input CreateProductInput) (*entity.Product, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Convert price to cents
	priceCents := money.ToCents(input.Price)

	// Create product
	product, err := entity.NewProduct(
		input.Name,
		input.Description,
		priceCents, // Use cents
		uc.defaultCurrency.Code,
		input.Stock,
		input.Weight,
		input.CategoryID,
		input.Images,
	)
	if err != nil {
		return nil, err
	}

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		product.Prices = make([]entity.ProductPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)

			product.Prices = append(product.Prices, entity.ProductPrice{
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
			})
		}
	}

	product.HasVariants = false

	// Save product
	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	// If product has variants, create them
	if len(input.Variants) > 0 {
		variants := make([]*entity.ProductVariant, 0, len(input.Variants))
		for _, variantInput := range input.Variants {
			// Convert variant prices to cents
			variantPriceCents := money.ToCents(variantInput.Price)

			variant, err := entity.NewProductVariant(
				product.ID,
				variantInput.SKU,
				variantPriceCents, // Use cents
				product.CurrencyCode,
				variantInput.Stock,
				variantInput.Attributes,
				variantInput.Images,
				variantInput.IsDefault,
			)
			if err != nil {
				return nil, err
			}

			// Process currency-specific prices for variant, if any
			if len(variantInput.CurrencyPrices) > 0 {
				variant.Prices = make([]entity.ProductVariantPrice, 0, len(variantInput.CurrencyPrices))

				for _, currPrice := range variantInput.CurrencyPrices {
					// Validate currency exists
					_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
					if err != nil {
						return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
					}

					// Convert price to cents
					priceCents := money.ToCents(currPrice.Price)

					variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
						CurrencyCode: currPrice.CurrencyCode,
						Price:        priceCents,
					})
				}
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
		product.HasVariants = true
	} else {
		// ALL PRODUCTS MUST HAVE AT LEAST ONE VARIANT
		// Create a default variant using the product's basic information
		defaultVariant, err := entity.NewDefaultProductVariant(
			product.ID,
			product.ProductNumber, // Use product number as SKU
			product.Price,         // Use product price
			product.CurrencyCode,
			product.Stock, // Use product stock
		)
		if err != nil {
			return nil, err
		}

		// Save the default variant
		if err := uc.productVariantRepo.Create(defaultVariant); err != nil {
			return nil, err
		}

		// Add variant to product
		product.Variants = []*entity.ProductVariant{defaultVariant}
		product.HasVariants = true
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
	Name           string
	Description    string
	Price          float64
	Stock          int
	CategoryID     uint
	Images         []string
	CurrencyPrices []CurrencyPriceInput
	Active         bool
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
	if input.Price > 0 && !product.HasVariants {
		product.Price = money.ToCents(input.Price) // Convert to cents
	}
	if input.Stock >= 0 && !product.HasVariants {
		product.Stock = input.Stock
	}
	if len(input.Images) > 0 {
		product.Images = input.Images
	}
	if input.Active != product.Active {
		product.Active = input.Active
	}

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		// Clear existing prices
		product.Prices = make([]entity.ProductPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)

			product.Prices = append(product.Prices, entity.ProductPrice{
				ProductID:    product.ID,
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
			})
		}
	}

	// Update product in repository
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateVariantInput contains the data needed to update a product variant (prices in dollars)
type UpdateVariantInput struct {
	SKU            string
	Price          float64
	Stock          int
	Attributes     []entity.VariantAttribute
	Images         []string
	IsDefault      bool
	CurrencyPrices []CurrencyPriceInput
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

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		// Clear existing prices
		variant.Prices = make([]entity.ProductVariantPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)

			variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
				VariantID:    variant.ID,
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
			})
		}
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
	ProductID      uint
	SKU            string
	Price          float64
	Stock          int
	Attributes     []entity.VariantAttribute
	Images         []string
	IsDefault      bool
	CurrencyPrices []CurrencyPriceInput
}

// AddVariant adds a new variant to a product
func (uc *ProductUseCase) AddVariant(input AddVariantInput) (*entity.ProductVariant, error) {
	product, err := uc.productRepo.GetByIDWithVariants(input.ProductID)
	if err != nil {
		return nil, err
	}

	// Convert prices to cents
	priceCents := money.ToCents(input.Price)

	// Create variant
	variant, err := entity.NewProductVariant(
		input.ProductID,
		input.SKU,
		priceCents, // Use cents
		product.CurrencyCode,
		input.Stock,
		input.Attributes,
		input.Images,
		input.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	// Process currency-specific prices, if any
	if len(input.CurrencyPrices) > 0 {
		variant.Prices = make([]entity.ProductVariantPrice, 0, len(input.CurrencyPrices))

		for _, currPrice := range input.CurrencyPrices {
			// Validate currency exists
			_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
			if err != nil {
				return nil, errors.New("invalid currency code: " + currPrice.CurrencyCode)
			}

			// Convert price to cents
			priceCents := money.ToCents(currPrice.Price)

			variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
				CurrencyCode: currPrice.CurrencyCode,
				Price:        priceCents,
			})
		}
	}

	// If this is the first variant or it's set as default, update product
	isFirstVariant := !product.HasVariants
	if isFirstVariant || input.IsDefault {
		// If this is the first variant, set product to have variants
		if isFirstVariant {
			product.HasVariants = true
		}

		// If this is the default variant, update product price and weight
		// TODO: Handle this in the repository
		if input.IsDefault && isFirstVariant {
			// If there are other variants, unset their default status
			variants := product.Variants

			for _, v := range variants {
				if v.IsDefault {
					v.IsDefault = false
					if err := uc.productVariantRepo.Update(v); err != nil {
						return nil, err
					}
				}
			}
		}

		// Update product
		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
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

// DeleteProduct deletes a product
func (uc *ProductUseCase) DeleteProduct(id uint) error {
	if id == 0 {
		return errors.New("product ID is required")
	}

	// TODO: make sure no orders are associated with the product

	return uc.productRepo.Delete(id)
}

// SearchProductsInput contains the data needed to search for products (prices in dollars)
type SearchProductsInput struct {
	Query        string  `json:"query"`
	CategoryID   uint    `json:"category_id"`
	MinPrice     float64 `json:"min_price"`     // Price in dollars
	MaxPrice     float64 `json:"max_price"`     // Price in dollars
	CurrencyCode string  `json:"currency_code"` // Optional currency code for prices
	Offset       int     `json:"offset"`
	Limit        int     `json:"limit"`
}

// SearchProducts searches for products based on criteria
func (uc *ProductUseCase) SearchProducts(input SearchProductsInput) ([]*entity.Product, int, error) {
	// If currency is specified and not the default, convert price ranges
	var minPriceCents, maxPriceCents int64

	if input.CurrencyCode != "" && input.CurrencyCode != uc.defaultCurrency.Code {
		// Get the currency
		currency, err := uc.currencyRepo.GetByCode(input.CurrencyCode)
		if err != nil {
			return nil, 0, errors.New("invalid currency code: " + input.CurrencyCode)
		}

		// Convert min/max prices to default currency using exchange rate
		defaultPrice := input.MinPrice / currency.ExchangeRate
		minPriceCents = money.ToCents(defaultPrice)

		defaultPrice = input.MaxPrice / currency.ExchangeRate
		maxPriceCents = money.ToCents(defaultPrice)
	} else {
		// Convert min/max prices to cents for repository search
		minPriceCents = money.ToCents(input.MinPrice)
		maxPriceCents = money.ToCents(input.MaxPrice)
	}

	products, err := uc.productRepo.Search(
		input.Query,
		input.CategoryID,
		minPriceCents, // Pass cents
		maxPriceCents, // Pass cents
		input.Offset,
		input.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.productRepo.CountSearch(
		input.Query,
		input.CategoryID,
		minPriceCents, // Pass cents
		maxPriceCents, // Pass cents
	)
	if err != nil {
		return products, 0, err
	}

	return products, total, nil
}

// ListProducts lists all products with pagination and returns total count
func (uc *ProductUseCase) ListProducts(offset, limit int) ([]*entity.Product, int, error) {
	products, err := uc.productRepo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.productRepo.Count()
	if err != nil {
		return products, 0, err
	}

	return products, total, nil
}

// ListCategories lists all product categories
func (uc *ProductUseCase) ListCategories() ([]*entity.Category, error) {
	return uc.categoryRepo.List()
}

// SetProductCurrencyPrices sets currency-specific prices for a product
func (uc *ProductUseCase) SetProductCurrencyPrices(productID uint, currencyPrices []CurrencyPriceInput) error {
	// Get product to check ownership
	product, err := uc.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	// Clear existing currency prices
	product.Prices = make([]entity.ProductPrice, 0, len(currencyPrices))

	// Add new currency prices
	for _, currPrice := range currencyPrices {
		// Validate currency exists
		_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
		if err != nil {
			return errors.New("invalid currency code: " + currPrice.CurrencyCode)
		}

		// Convert prices to cents
		priceCents := money.ToCents(currPrice.Price)

		product.Prices = append(product.Prices, entity.ProductPrice{
			ProductID:    productID,
			CurrencyCode: currPrice.CurrencyCode,
			Price:        priceCents,
		})
	}

	// Update product in repository
	return uc.productRepo.Update(product)
}

// SetVariantCurrencyPrices sets currency-specific prices for a product variant
func (uc *ProductUseCase) SetVariantCurrencyPrices(productID uint, variantID uint, currencyPrices []CurrencyPriceInput) error {
	// Get variant
	variant, err := uc.productVariantRepo.GetByID(variantID)
	if err != nil {
		return err
	}

	// Clear existing currency prices
	variant.Prices = make([]entity.ProductVariantPrice, 0, len(currencyPrices))

	// Add new currency prices
	for _, currPrice := range currencyPrices {
		// Validate currency exists
		_, err := uc.currencyRepo.GetByCode(currPrice.CurrencyCode)
		if err != nil {
			return errors.New("invalid currency code: " + currPrice.CurrencyCode)
		}

		// Convert prices to cents
		priceCents := money.ToCents(currPrice.Price)

		variant.Prices = append(variant.Prices, entity.ProductVariantPrice{
			VariantID:    variantID,
			CurrencyCode: currPrice.CurrencyCode,
			Price:        priceCents,
		})
	}

	// Update variant in repository
	return uc.productVariantRepo.Update(variant)
}
