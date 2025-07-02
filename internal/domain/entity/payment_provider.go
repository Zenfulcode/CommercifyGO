package entity

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PaymentProvider represents a payment provider configuration in the system
type PaymentProvider struct {
	gorm.Model
	Type                common.PaymentProviderType `gorm:"uniqueIndex;not null;size:50" json:"type"`
	Name                string                     `gorm:"not null;size:100" json:"name"`
	Description         string                     `gorm:"size:500" json:"description"`
	IconURL             string                     `gorm:"size:500" json:"icon_url,omitempty"`
	Methods             []common.PaymentMethod     `json:"methods"`
	Enabled             bool                       `gorm:"default:true" json:"enabled"`
	SupportedCurrencies []string                   `json:"supported_currencies,omitempty"`
	Configuration       datatypes.JSONMap          `json:"configuration,omitempty"`
	WebhookURL          string                     `gorm:"size:500" json:"webhook_url,omitempty"`
	WebhookSecret       string                     `gorm:"size:255" json:"webhook_secret,omitempty"`
	WebhookEvents       []string                   `json:"webhook_events,omitempty"`
	ExternalWebhookID   string                     `gorm:"size:255" json:"external_webhook_id,omitempty"`
	IsTestMode          bool                       `gorm:"default:false" json:"is_test_mode"`
	Priority            int                        `gorm:"default:0" json:"priority"` // Higher priority means higher preference
}

// Validate validates the payment provider data
func (p *PaymentProvider) Validate() error {
	if p.Type == "" {
		return errors.New("payment provider type is required")
	}

	if p.Name == "" {
		return errors.New("payment provider name is required")
	}

	if len(p.Methods) == 0 {
		return errors.New("at least one payment method is required")
	}

	// Validate payment methods
	for _, method := range p.Methods {
		if method != common.PaymentMethodCreditCard &&
			method != common.PaymentMethodWallet {
			return errors.New("invalid payment method: " + string(method))
		}
	}

	return nil
}

// SetMethods sets the payment methods for this provider
func (p *PaymentProvider) SetMethods(methods []common.PaymentMethod) {
	p.Methods = methods
}

// SetSupportedCurrencies sets the supported currencies for this provider
func (p *PaymentProvider) SetSupportedCurrencies(currencies []string) {
	p.SupportedCurrencies = currencies
}

// SetWebhookEvents sets the webhook events for this provider
func (p *PaymentProvider) SetWebhookEvents(events []string) {
	p.WebhookEvents = events
}

// SetConfiguration sets the configuration for this provider
func (p *PaymentProvider) SetConfiguration(config map[string]interface{}) {
	if config == nil {
		p.Configuration = nil
		return
	}

	p.Configuration = datatypes.JSONMap(config)
}

// GetConfigurationJSON returns the configuration as a JSON string
func (p *PaymentProvider) GetConfiguration() (string, error) {
	if p.Configuration == nil {
		return "{}", nil // Return empty JSON if no configuration
	}

	jsonData, err := p.Configuration.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (p *PaymentProvider) GetConfigurationField(fieldName string) (interface{}, error) {
	if p.Configuration == nil {
		return nil, errors.New("configuration is nil")
	}

	if p.Configuration[fieldName] == nil {
		return nil, errors.New("field not found")
	}

	return p.Configuration[fieldName], nil
}

// SupportsCurrency checks if the provider supports a specific currency
func (p *PaymentProvider) SupportsCurrency(currency string) bool {
	if len(p.SupportedCurrencies) == 0 {
		return true // If no currencies specified, assume it supports all
	}

	for _, supportedCurrency := range p.SupportedCurrencies {
		if supportedCurrency == currency {
			return true
		}
	}

	return false
}

// SupportsMethod checks if the provider supports a specific payment method
func (p *PaymentProvider) SupportsMethod(method common.PaymentMethod) bool {
	for _, supportedMethod := range p.Methods {
		if supportedMethod == method {
			return true
		}
	}
	return false
}

// PaymentProviderInfo represents payment provider information for API responses
type PaymentProviderInfo struct {
	Type                common.PaymentProviderType `json:"type"`
	Name                string                     `json:"name"`
	Description         string                     `json:"description"`
	IconURL             string                     `json:"icon_url,omitempty"`
	Methods             []common.PaymentMethod     `json:"methods"`
	Enabled             bool                       `json:"enabled"`
	SupportedCurrencies []string                   `json:"supported_currencies,omitempty"`
}

// ToPaymentProviderInfo converts the entity to PaymentProviderInfo for API responses
func (p *PaymentProvider) ToPaymentProviderInfo() PaymentProviderInfo {
	return PaymentProviderInfo{
		Type:                p.Type,
		Name:                p.Name,
		Description:         p.Description,
		IconURL:             p.IconURL,
		Methods:             p.Methods,
		Enabled:             p.Enabled,
		SupportedCurrencies: p.SupportedCurrencies,
	}
}
