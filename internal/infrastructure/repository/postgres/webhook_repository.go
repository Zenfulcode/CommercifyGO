package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// WebhookRepository implements the WebhookRepository interface using PostgreSQL
type WebhookRepository struct {
	db *sql.DB
}

// NewWebhookRepository creates a new WebhookRepository
func NewWebhookRepository(db *sql.DB) repository.WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

// Create creates a new webhook
func (r *WebhookRepository) Create(webhook *entity.Webhook) error {
	// Convert events to JSON string
	eventsJSON, err := json.Marshal(webhook.Events)
	if err != nil {
		return err
	}

	// Set timestamp
	now := time.Now()
	webhook.CreatedAt = now
	webhook.UpdatedAt = now

	// Insert webhook
	query := `
		INSERT INTO webhooks (
			provider, external_id, url, events, secret, active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err = r.db.QueryRow(
		query,
		webhook.Provider,
		webhook.ExternalID,
		webhook.URL,
		eventsJSON,
		webhook.Secret,
		webhook.Active,
		webhook.CreatedAt,
		webhook.UpdatedAt,
	).Scan(&webhook.ID)

	return err
}

// Update updates an existing webhook
func (r *WebhookRepository) Update(webhook *entity.Webhook) error {
	// Convert events to JSON string
	eventsJSON, err := json.Marshal(webhook.Events)
	if err != nil {
		return err
	}

	// Update timestamp
	webhook.UpdatedAt = time.Now()

	// Update webhook
	query := `
		UPDATE webhooks
		SET provider = $1, external_id = $2, url = $3, events = $4, secret = $5, active = $6, updated_at = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(
		query,
		webhook.Provider,
		webhook.ExternalID,
		webhook.URL,
		eventsJSON,
		webhook.Secret,
		webhook.Active,
		webhook.UpdatedAt,
		webhook.ID,
	)
	if err != nil {
		return err
	}

	// Check if webhook exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("webhook not found")
	}

	return nil
}

// Delete deletes a webhook
func (r *WebhookRepository) Delete(id uint) error {
	query := `DELETE FROM webhooks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	// Check if webhook exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("webhook not found")
	}

	return nil
}

// GetByID returns a webhook by ID
func (r *WebhookRepository) GetByID(id uint) (*entity.Webhook, error) {
	query := `
		SELECT id, provider, external_id, url, events, secret, active, created_at, updated_at
		FROM webhooks
		WHERE id = $1
	`

	webhook := &entity.Webhook{}
	var eventsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&webhook.ID,
		&webhook.Provider,
		&webhook.ExternalID,
		&webhook.URL,
		&eventsJSON,
		&webhook.Secret,
		&webhook.Active,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("webhook not found")
		}
		return nil, err
	}

	// Parse events JSON
	err = webhook.SetEventsFromJSON(eventsJSON)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

// GetByProvider returns all webhooks for a specific provider
func (r *WebhookRepository) GetByProvider(provider string) ([]*entity.Webhook, error) {
	query := `
		SELECT id, provider, external_id, url, events, secret, active, created_at, updated_at
		FROM webhooks
		WHERE provider = $1
	`

	rows, err := r.db.Query(query, provider)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	webhooks := []*entity.Webhook{}
	for rows.Next() {
		webhook := &entity.Webhook{}
		var eventsJSON []byte

		err := rows.Scan(
			&webhook.ID,
			&webhook.Provider,
			&webhook.ExternalID,
			&webhook.URL,
			&eventsJSON,
			&webhook.Secret,
			&webhook.Active,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse events JSON
		err = webhook.SetEventsFromJSON(eventsJSON)
		if err != nil {
			return nil, err
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// GetActive returns all active webhooks
func (r *WebhookRepository) GetActive() ([]*entity.Webhook, error) {
	query := `
		SELECT id, provider, external_id, url, events, secret, active, created_at, updated_at
		FROM webhooks
		WHERE active = true
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	webhooks := []*entity.Webhook{}
	for rows.Next() {
		webhook := &entity.Webhook{}
		var eventsJSON []byte

		err := rows.Scan(
			&webhook.ID,
			&webhook.Provider,
			&webhook.ExternalID,
			&webhook.URL,
			&eventsJSON,
			&webhook.Secret,
			&webhook.Active,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse events JSON
		err = webhook.SetEventsFromJSON(eventsJSON)
		if err != nil {
			return nil, err
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// GetByExternalID returns a webhook by external ID
func (r *WebhookRepository) GetByExternalID(provider string, externalID string) (*entity.Webhook, error) {
	query := `
		SELECT id, provider, external_id, url, events, secret, active, created_at, updated_at
		FROM webhooks
		WHERE provider = $1 AND external_id = $2
	`

	webhook := &entity.Webhook{}
	var eventsJSON []byte

	err := r.db.QueryRow(query, provider, externalID).Scan(
		&webhook.ID,
		&webhook.Provider,
		&webhook.ExternalID,
		&webhook.URL,
		&eventsJSON,
		&webhook.Secret,
		&webhook.Active,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("webhook not found")
		}
		return nil, err
	}

	// Parse events JSON
	err = webhook.SetEventsFromJSON(eventsJSON)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}
