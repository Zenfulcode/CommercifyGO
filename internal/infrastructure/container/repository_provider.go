package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/gorm"
)

// RepositoryProvider provides access to all repositories
type RepositoryProvider interface {
	UserRepository() repository.UserRepository
	ProductRepository() repository.ProductRepository
	ProductVariantRepository() repository.ProductVariantRepository
	CategoryRepository() repository.CategoryRepository
	OrderRepository() repository.OrderRepository
	CheckoutRepository() repository.CheckoutRepository
	DiscountRepository() repository.DiscountRepository
	PaymentProviderRepository() repository.PaymentProviderRepository
	PaymentTransactionRepository() repository.PaymentTransactionRepository
	CurrencyRepository() repository.CurrencyRepository

	// Shipping related repository
	ShippingMethodRepository() repository.ShippingMethodRepository
	ShippingZoneRepository() repository.ShippingZoneRepository
	ShippingRateRepository() repository.ShippingRateRepository
}

// repositoryProvider is the concrete implementation of RepositoryProvider
type repositoryProvider struct {
	container Container
	mu        sync.Mutex

	userRepo            repository.UserRepository
	productVariantRepo  repository.ProductVariantRepository
	productRepo         repository.ProductRepository
	categoryRepo        repository.CategoryRepository
	orderRepo           repository.OrderRepository
	checkoutRepo        repository.CheckoutRepository
	discountRepo        repository.DiscountRepository
	paymentProviderRepo repository.PaymentProviderRepository
	paymentTrxRepo      repository.PaymentTransactionRepository
	currencyRepo        repository.CurrencyRepository

	shippingMethodRepo repository.ShippingMethodRepository
	shippingZoneRepo   repository.ShippingZoneRepository
	shippingRateRepo   repository.ShippingRateRepository
}

// NewRepositoryProvider creates a new repository provider
func NewRepositoryProvider(container Container) RepositoryProvider {
	return &repositoryProvider{
		container: container,
	}
}

// UserRepository returns the user repository
func (p *repositoryProvider) UserRepository() repository.UserRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userRepo == nil {
		p.userRepo = gorm.NewUserRepository(p.container.DB())
	}
	return p.userRepo
}

// ProductRepository returns the product repository
func (p *repositoryProvider) ProductRepository() repository.ProductRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productRepo == nil {
		// Initialize both repositories under the same lock
		if p.productVariantRepo == nil {
			p.productVariantRepo = gorm.NewProductVariantRepository(p.container.DB())
		}
		p.productRepo = gorm.NewProductRepository(p.container.DB())
	}
	return p.productRepo
}

// ProductVariantRepository returns the product variant repository
func (p *repositoryProvider) ProductVariantRepository() repository.ProductVariantRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productVariantRepo == nil {
		p.productVariantRepo = gorm.NewProductVariantRepository(p.container.DB())
	}
	return p.productVariantRepo
}

// CategoryRepository returns the category repository
func (p *repositoryProvider) CategoryRepository() repository.CategoryRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.categoryRepo == nil {
		p.categoryRepo = gorm.NewCategoryRepository(p.container.DB())
	}
	return p.categoryRepo
}

// OrderRepository returns the order repository
func (p *repositoryProvider) OrderRepository() repository.OrderRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.orderRepo == nil {
		p.orderRepo = gorm.NewOrderRepository(p.container.DB())
	}
	return p.orderRepo
}

// CheckoutRepository returns the checkout repository
func (p *repositoryProvider) CheckoutRepository() repository.CheckoutRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.checkoutRepo == nil {
		p.checkoutRepo = gorm.NewCheckoutRepository(p.container.DB())
	}
	return p.checkoutRepo
}

// DiscountRepository returns the discount repository
func (p *repositoryProvider) DiscountRepository() repository.DiscountRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.discountRepo == nil {
		p.discountRepo = gorm.NewDiscountRepository(p.container.DB())
	}
	return p.discountRepo
}

// PaymentProviderRepository returns the payment provider repository
func (p *repositoryProvider) PaymentProviderRepository() repository.PaymentProviderRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.paymentProviderRepo == nil {
		p.paymentProviderRepo = gorm.NewPaymentProviderRepository(p.container.DB())
	}
	return p.paymentProviderRepo
}

// PaymentTransactionRepository returns the payment transaction repository
func (p *repositoryProvider) PaymentTransactionRepository() repository.PaymentTransactionRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.paymentTrxRepo == nil {
		p.paymentTrxRepo = gorm.NewTransactionRepository(p.container.DB())
	}
	return p.paymentTrxRepo
}

// ShippingMethodRepository returns the shipping method repository
func (p *repositoryProvider) ShippingMethodRepository() repository.ShippingMethodRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shippingMethodRepo == nil {
		p.shippingMethodRepo = gorm.NewShippingMethodRepository(p.container.DB())
	}
	return p.shippingMethodRepo
}

// ShippingZoneRepository returns the shipping zone repository
func (p *repositoryProvider) ShippingZoneRepository() repository.ShippingZoneRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shippingZoneRepo == nil {
		p.shippingZoneRepo = gorm.NewShippingZoneRepository(p.container.DB())
	}
	return p.shippingZoneRepo
}

// ShippingRateRepository returns the shipping rate repository
func (p *repositoryProvider) ShippingRateRepository() repository.ShippingRateRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shippingRateRepo == nil {
		p.shippingRateRepo = gorm.NewShippingRateRepository(p.container.DB())
	}
	return p.shippingRateRepo
}

// CurrencyRepository returns the currency repository
func (p *repositoryProvider) CurrencyRepository() repository.CurrencyRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.currencyRepo == nil {
		p.currencyRepo = gorm.NewCurrencyRepository(p.container.DB())
	}
	return p.currencyRepo
}
