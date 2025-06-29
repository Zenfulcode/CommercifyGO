package dto

import (
	"time"
)

// CurrencyDTO represents a currency entity
type CurrencyDTO struct {
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	ExchangeRate float64   `json:"exchange_rate"`
	IsEnabled    bool      `json:"is_enabled"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
