package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/application/usecase"
)

// UseCaseProvider provides access to all use cases
type UseCaseProvider interface {
	UserUseCase() *usecase.UserUseCase
	ProductUseCase() *usecase.ProductUseCase
	CategoryUseCase() *usecase.CategoryUseCase
	CheckoutUseCase() *usecase.CheckoutUseCase
	OrderUseCase() *usecase.OrderUseCase
	DiscountUseCase() *usecase.DiscountUseCase
	ShippingUseCase() *usecase.ShippingUseCase
	CurrencyUsecase() *usecase.CurrencyUseCase
}

// useCaseProvider is the concrete implementation of UseCaseProvider
type useCaseProvider struct {
	container Container
	mu        sync.Mutex

	userUseCase     *usecase.UserUseCase
	productUseCase  *usecase.ProductUseCase
	categoryUseCase *usecase.CategoryUseCase
	checkoutUseCase *usecase.CheckoutUseCase
	orderUseCase    *usecase.OrderUseCase
	discountUseCase *usecase.DiscountUseCase
	shippingUseCase *usecase.ShippingUseCase
	currencyUseCase *usecase.CurrencyUseCase
}

// NewUseCaseProvider creates a new use case provider
func NewUseCaseProvider(container Container) UseCaseProvider {
	return &useCaseProvider{
		container: container,
	}
}

// UserUseCase returns the user use case
func (p *useCaseProvider) UserUseCase() *usecase.UserUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userUseCase == nil {
		p.userUseCase = usecase.NewUserUseCase(
			p.container.Repositories().UserRepository(),
		)
	}
	return p.userUseCase
}

// ProductUseCase returns the product use case
func (p *useCaseProvider) ProductUseCase() *usecase.ProductUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productUseCase == nil {
		p.productUseCase = usecase.NewProductUseCase(
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().CategoryRepository(),
			p.container.Repositories().ProductVariantRepository(),
			p.container.Repositories().CurrencyRepository(),
			p.container.Repositories().OrderRepository(),
			p.container.Repositories().CheckoutRepository(),
		)
	}
	return p.productUseCase
}

// CategoryUseCase returns the category use case
func (p *useCaseProvider) CategoryUseCase() *usecase.CategoryUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.categoryUseCase == nil {
		p.categoryUseCase = usecase.NewCategoryUseCase(
			p.container.Repositories().CategoryRepository(),
		)
	}
	return p.categoryUseCase
}

// CheckoutUseCase returns the checkout use case
func (p *useCaseProvider) CheckoutUseCase() *usecase.CheckoutUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.checkoutUseCase == nil {
		// Initialize shipping use case directly to avoid circular dependency
		if p.shippingUseCase == nil {
			p.shippingUseCase = usecase.NewShippingUseCase(
				p.container.Repositories().ShippingMethodRepository(),
				p.container.Repositories().ShippingZoneRepository(),
				p.container.Repositories().ShippingRateRepository(),
			)
		}

		p.checkoutUseCase = usecase.NewCheckoutUseCase(
			p.container.Repositories().CheckoutRepository(),
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().ProductVariantRepository(),
			p.container.Repositories().ShippingMethodRepository(),
			p.container.Repositories().ShippingRateRepository(),
			p.container.Repositories().DiscountRepository(),
			p.container.Repositories().OrderRepository(),
			p.container.Repositories().CurrencyRepository(),
			p.container.Repositories().PaymentTransactionRepository(),
			p.container.Services().PaymentService(),
			p.shippingUseCase,
		)
	}
	return p.checkoutUseCase
}

// OrderUseCase returns the order use case
func (p *useCaseProvider) OrderUseCase() *usecase.OrderUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.orderUseCase == nil {
		p.orderUseCase = usecase.NewOrderUseCase(
			p.container.Repositories().OrderRepository(),
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().ProductVariantRepository(),
			p.container.Repositories().UserRepository(),
			p.container.Services().PaymentService(),
			p.container.Services().EmailService(),
			p.container.Repositories().PaymentTransactionRepository(),
			p.container.Repositories().CurrencyRepository(),
		)
	}
	return p.orderUseCase
}

// DiscountUseCase returns the discount use case
func (p *useCaseProvider) DiscountUseCase() *usecase.DiscountUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.discountUseCase == nil {
		p.discountUseCase = usecase.NewDiscountUseCase(
			p.container.Repositories().DiscountRepository(),
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().CategoryRepository(),
			p.container.Repositories().OrderRepository(),
		)
	}
	return p.discountUseCase
}

// ShippingUseCase returns the shipping use case
func (p *useCaseProvider) ShippingUseCase() *usecase.ShippingUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shippingUseCase == nil {
		p.shippingUseCase = usecase.NewShippingUseCase(
			p.container.Repositories().ShippingMethodRepository(),
			p.container.Repositories().ShippingZoneRepository(),
			p.container.Repositories().ShippingRateRepository(),
		)
	}
	return p.shippingUseCase
}

// CurrencyUsecase returns the currency use case
func (p *useCaseProvider) CurrencyUsecase() *usecase.CurrencyUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.currencyUseCase == nil {
		p.currencyUseCase = usecase.NewCurrencyUseCase(
			p.container.Repositories().CurrencyRepository(),
		)

		var defaultCurrency usecase.CurrencyInput = usecase.CurrencyInput{
			Code:         p.container.Config().DefaultCurrency,
			Name:         "Default Currency",
			Symbol:       "$",
			ExchangeRate: 1.0,
			IsEnabled:    true,
			IsDefault:    true,
		}

		p.currencyUseCase.CreateCurrency(defaultCurrency)
	}
	return p.currencyUseCase
}
