package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/auth"
	"github.com/zenfulcode/commercify/internal/infrastructure/email"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
)

// ServiceProvider provides access to all services
type ServiceProvider interface {
	JWTService() *auth.JWTService
	PaymentService() service.PaymentService
	PaymentProviderService() service.PaymentProviderService
	EmailService() service.EmailService
	MobilePayService() *payment.MobilePayPaymentService
	InitializeMobilePay() *payment.MobilePayPaymentService
}

// serviceProvider is the concrete implementation of ServiceProvider
type serviceProvider struct {
	container Container
	mu        sync.Mutex

	jwtService             *auth.JWTService
	paymentService         service.PaymentService
	paymentProviderService service.PaymentProviderService
	emailService           service.EmailService
	mobilePayService       *payment.MobilePayPaymentService
}

// NewServiceProvider creates a new service provider
func NewServiceProvider(container Container) ServiceProvider {
	return &serviceProvider{
		container: container,
	}
}

// JWTService returns the JWT authentication service
func (p *serviceProvider) JWTService() *auth.JWTService {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.jwtService == nil {
		p.jwtService = auth.NewJWTService(p.container.Config().Auth)
	}
	return p.jwtService
}

// PaymentService returns the payment service
func (p *serviceProvider) PaymentService() service.PaymentService {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.paymentService == nil {
		multiProviderService := payment.NewMultiProviderPaymentService(
			p.container.Config(),
			p.container.Repositories().PaymentProviderRepository(),
			p.container.Logger(),
		)
		p.paymentService = multiProviderService

		// TODO: Get rid of this
		// Extract MobilePay service for webhook registration if it exists
		// We need to access the actual MultiProviderPaymentService concrete type
		// to access its GetProviders method
		for _, providerWithService := range multiProviderService.GetProviders() {
			if providerWithService.Type == common.PaymentProviderMobilePay {
				// Cast the generic service to the concrete MobilePayPaymentService type
				if mobilePayService, ok := providerWithService.Service.(*payment.MobilePayPaymentService); ok {
					p.mobilePayService = mobilePayService
					break
				}
			}
		}
	}
	return p.paymentService
}

// PaymentProviderService returns the payment provider service
func (p *serviceProvider) PaymentProviderService() service.PaymentProviderService {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.paymentProviderService == nil {
		p.paymentProviderService = payment.NewPaymentProviderService(
			p.container.Repositories().PaymentProviderRepository(),
			p.container.Config(),
			p.container.Logger(),
		)
	}
	return p.paymentProviderService
}

// InitializeMobilePay directly initializes the MobilePay service to break circular dependency
func (p *serviceProvider) InitializeMobilePay() *payment.MobilePayPaymentService {
	if !p.container.Config().MobilePay.Enabled {
		return nil
	}

	return payment.NewMobilePayPaymentService(p.container.Config().MobilePay, p.container.Logger())
}

// MobilePayService returns the MobilePay payment service
func (p *serviceProvider) MobilePayService() *payment.MobilePayPaymentService {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.mobilePayService == nil {
		// Instead of calling PaymentService() which would create a deadlock
		// We directly initialize MobilePay if needed
		p.mobilePayService = p.InitializeMobilePay()
	}
	return p.mobilePayService
}

// EmailService returns the email service
func (p *serviceProvider) EmailService() service.EmailService {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.emailService == nil {
		p.emailService = email.NewSMTPEmailService(p.container.Config().Email, p.container.Logger())
	}
	return p.emailService
}
