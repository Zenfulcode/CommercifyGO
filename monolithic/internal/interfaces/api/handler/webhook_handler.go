package handler

import (
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// WebhookHandlerProvider provides webhook handlers for different payment providers
type WebhookHandlerProvider struct {
	stripeHandler    *StripeWebhookHandler
	mobilePayHandler *MobilePayWebhookHandler
}

// NewWebhookHandlerProvider creates a new WebhookHandlerProvider
func NewWebhookHandlerProvider(
	orderUseCase *usecase.OrderUseCase,
	paymentProviderService service.PaymentProviderService,
	cfg *config.Config,
	logger logger.Logger,
) *WebhookHandlerProvider {
	return &WebhookHandlerProvider{
		stripeHandler:    NewStripeWebhookHandler(orderUseCase, cfg, logger),
		mobilePayHandler: NewMobilePayWebhookHandler(orderUseCase, paymentProviderService, cfg, logger),
	}
}

// StripeHandler returns the Stripe webhook handler
func (p *WebhookHandlerProvider) StripeHandler() *StripeWebhookHandler {
	return p.stripeHandler
}

// MobilePayHandler returns the MobilePay webhook handler
func (p *WebhookHandlerProvider) MobilePayHandler() *MobilePayWebhookHandler {
	return p.mobilePayHandler
}
