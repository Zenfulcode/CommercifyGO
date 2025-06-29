package entity

import (
	"encoding/json"
	"time"
)

// Webhook represents a registered webhook for receiving event notifications
type Webhook struct {
	ID         uint
	Provider   string
	ExternalID string
	URL        string
	Events     []string
	Secret     string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Validate validates the webhook data
func (w *Webhook) Validate() error {
	if w.Provider == "" {
		return ErrInvalidInput{Field: "provider", Message: "provider is required"}
	}
	if w.URL == "" {
		return ErrInvalidInput{Field: "url", Message: "url is required"}
	}
	if len(w.Events) == 0 {
		return ErrInvalidInput{Field: "events", Message: "at least one event is required"}
	}
	return nil
}

// SetEvents sets the events for this webhook
func (w *Webhook) SetEvents(events []string) {
	w.Events = events
}

// GetEventsJSON returns the events as a JSON string
func (w *Webhook) GetEventsJSON() (string, error) {
	eventsJSON, err := json.Marshal(w.Events)
	if err != nil {
		return "", err
	}
	return string(eventsJSON), nil
}

// SetEventsFromJSON sets the events from a JSON string
func (w *Webhook) SetEventsFromJSON(eventsJSON []byte) error {
	var events []string
	if err := json.Unmarshal(eventsJSON, &events); err != nil {
		return err
	}
	w.Events = events
	return nil
}
