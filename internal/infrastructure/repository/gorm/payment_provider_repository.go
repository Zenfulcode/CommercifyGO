package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PaymentProviderRepository implements repository.PaymentProviderRepository using GORM
type PaymentProviderRepository struct {
	db *gorm.DB
}

func (r *PaymentProviderRepository) buildJSONContainsQuery(column string, value string) (string, []interface{}) {
	dialect := r.db.Dialector.Name()

	switch dialect {
	case "postgres":
		// PostgreSQL uses the '?' operator for JSONB array containment (element existence)
		// Or, for string containment within a JSON array of strings, it's often more robust
		// to cast to text[] and use the ANY operator.
		// However, for datatypes.JSONSlice[string], the @> operator is suitable for checking
		// if a JSON array contains another JSON array (in this case, a single-element array).
		// The `datatypes.JSON` helper from GORM can generate this for you.
		// For an exact match within an array of strings, the `?` operator is for top-level keys.
		// For elements in a JSON array, you often use `jsonb_array_elements_text` or `@>`
		//
		// A common way to check if an element exists in a JSON array in PostgreSQL:
		// SELECT * FROM your_table WHERE your_jsonb_array_column @> '["your_value"]'::jsonb;
		//
		// GORM's `datatypes.JSONQuery` is the preferred way for cross-database JSON querying.
		// It abstracts the underlying SQL for you.

		// For checking if an array contains a specific string, we can use `datatypes.JSONQuery` with `Contains`.
		// However, `Contains` typically works for JSON objects. For array elements, a raw SQL approach
		// using `@>` or `?` with casting is often needed if `datatypes.JSONQuery` doesn't directly
		// provide a method for exact string containment in a JSON array.

		// Let's use `datatypes.JSONQuery` which is designed for this.
		// GORM's datatypes.JSONQuery("column").Contains(value, "path_to_array_element")
		// The Contains method on JSONQuery is more for checking if a JSON object contains key/value.
		// For checking if an array of strings contains a specific string, we generally need to be more explicit.

		// Option 1: Using the `@>` operator (JSON containment)
		// This checks if the array `supported_currencies` contains the array `[currency]`
		return fmt.Sprintf("%s @> ?", column), []interface{}{datatypes.JSON(fmt.Sprintf(`["%s"]`, value))}

		// Option 2: More verbose, but also works for direct element checking if @> is not desired:
		// return fmt.Sprintf("EXISTS (SELECT 1 FROM jsonb_array_elements_text(%s) AS elem WHERE elem = ?)", column), []interface{}{value}

	case "sqlite":
		// SQLite uses `json_each` or `json_extract` functions.
		// The `json_each` function can be used to iterate over array elements.
		// SQLite also has the `->>` operator for extracting a value as text.
		// To check if a JSON array contains a string in SQLite:
		// SELECT * FROM your_table WHERE json_each(your_json_array_column).value = 'your_value';

		return fmt.Sprintf("EXISTS (SELECT 1 FROM json_each(%s) WHERE json_each.value = ?)", column), []interface{}{value}

	default:
		// Fallback for other databases or if not explicitly handled
		// This might not be optimal or even work for all databases.
		// You'd ideally add specific handling for MySQL etc. if needed.
		return fmt.Sprintf("%s LIKE ?", column), []interface{}{fmt.Sprintf(`%%"%s"%%`, value)}
	}
}

// Create implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) Create(provider *entity.PaymentProvider) error {
	if err := provider.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return r.db.Create(provider).Error
}

// Update implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) Update(provider *entity.PaymentProvider) error {
	if err := provider.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return r.db.Save(provider).Error
}

// Delete implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) Delete(id uint) error {
	return r.db.Delete(&entity.PaymentProvider{}, id).Error
}

// GetByID implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetByID(id uint) (*entity.PaymentProvider, error) {
	var provider entity.PaymentProvider
	if err := r.db.First(&provider, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment provider with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch payment provider: %w", err)
	}
	return &provider, nil
}

// GetByType implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetByType(providerType common.PaymentProviderType) (*entity.PaymentProvider, error) {
	var provider entity.PaymentProvider
	if err := r.db.Where("type = ?", providerType).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment provider with type %s not found", providerType)
		}
		return nil, fmt.Errorf("failed to fetch payment provider by type: %w", err)
	}
	return &provider, nil
}

// List implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) List(offset, limit int) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	query := r.db.Order("priority DESC, created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers: %w", err)
	}
	return providers, nil
}

// GetEnabled implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabled() ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	if err := r.db.
		Where("enabled = ?", true).
		Order("priority DESC, created_at ASC").
		Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enabled payment providers: %w", err)
	}
	return providers, nil
}

// GetEnabledByMethod implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabledByMethod(method common.PaymentMethod) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider

	quer, params := r.buildJSONContainsQuery("methods", string(method))

	if err := r.db.
		Where("enabled = ?", true).
		Where(quer, params).
		Order("priority DESC, created_at ASC").
		Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers by method: %w", err)
	}

	return providers, nil
}

// GetEnabledByCurrency implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabledByCurrency(currency string) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider

	quer, params := r.buildJSONContainsQuery("supported_currencies", currency)

	if err := r.db.
		Where("enabled = ?", true).
		Where(quer, params).
		Order("priority DESC, created_at ASC").
		Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers by currency: %w", err)
	}

	return providers, nil
}

// GetEnabledByMethodAndCurrency implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabledByMethodAndCurrency(method common.PaymentMethod, currency string) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider

	currencyQuery, currencyParams := r.buildJSONContainsQuery("supported_currencies", currency)
	methodsQuery, methodsParams := r.buildJSONContainsQuery("methods", string(method))

	if err := r.db.
		Where("enabled = ?", true).
		Where(methodsQuery, methodsParams).
		Where(currencyQuery, currencyParams).
		Order("priority DESC, created_at ASC").
		Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers by method and currency: %w", err)
	}

	return providers, nil
}

// UpdateWebhookInfo implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) UpdateWebhookInfo(providerType common.PaymentProviderType, webhookURL, webhookSecret, externalWebhookID string, events []string) error {
	updates := map[string]any{
		"webhook_url":         webhookURL,
		"webhook_secret":      webhookSecret,
		"external_webhook_id": externalWebhookID,
		"webhook_events":      events,
	}

	result := r.db.Model(&entity.PaymentProvider{}).Where("type = ?", providerType).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update webhook info: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment provider with type %s not found", providerType)
	}

	return nil
}

// GetWithWebhooks implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetWithWebhooks() ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	if err := r.db.
		Where("webhook_url IS NOT NULL AND webhook_url != ''").
		Order("priority DESC, created_at ASC").
		Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers with webhooks: %w", err)
	}
	return providers, nil
}

// NewPaymentProviderRepository creates a new GORM-based PaymentProviderRepository
func NewPaymentProviderRepository(db *gorm.DB) repository.PaymentProviderRepository {
	return &PaymentProviderRepository{db: db}
}
