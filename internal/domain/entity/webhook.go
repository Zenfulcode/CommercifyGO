package entity

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// Webhook represents a registered webhook for receiving event notifications
type Webhook struct {
	gorm.Model
	Provider   string   `gorm:"not null;size:100"`
	ExternalID string   `gorm:"uniqueIndex;not null;size:255"`
	URL        string   `gorm:"not null;size:500"`
	Events     []string `gorm:"type:jsonb;not null"`
	Secret     string   `gorm:"size:255"`
	IsActive   bool     `gorm:"default:true"`
}

// Validate validates the webhook data
func (w *Webhook) Validate() error {
	if w.Provider == "" {
		return errors.New("provider is required")
	}
	if w.URL == "" {
		return errors.New("url is required")
	}
	if len(w.Events) == 0 {
		return errors.New("at least one event is required")
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
