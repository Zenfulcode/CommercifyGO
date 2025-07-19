// Package money provides utilities for handling monetary values.
// Money is stored as integers (cents) in the database to avoid floating-point precision issues.
package money

import (
	"fmt"
	"math"
)

// ToCents converts a dollar amount (float64) to cents (int64)
// This avoids floating-point precision issues when storing money values
func ToCents(dollars float64) int64 {
	// Round to nearest cent to avoid floating point issues
	// Multiply by 100 to convert dollars to cents
	return int64(math.Round(dollars * 100))
}

// FromCents converts a cent amount (int64) to dollars (float64)
func FromCents(cents int64) float64 {
	// Divide by 100 to convert cents to dollars
	return float64(cents) / 100
}

// ConvertNullableToCents converts a nullable dollar amount (*float64) to nullable cents (*int64)
func ConvertNullableToCents(dollars *float64) *int64 {
	if dollars == nil {
		return nil
	}
	cents := ToCents(*dollars)
	return &cents
}

// ConvertNullableFromCents converts a nullable cent amount (*int64) to nullable dollars (*float64)
func ConvertNullableFromCents(cents *int64) *float64 {
	if cents == nil {
		return nil
	}
	dollars := FromCents(*cents)
	return &dollars
}

// FormatCurrency formats a cents value (int64) as a currency string
func FormatCurrency(cents int64, symbol string) string {
	// Convert to dollars
	dollars := FromCents(cents)
	return symbol + formatDollars(dollars)
}

// FormatDollars formats a dollars value as a currency string
func FormatDollars(dollars float64, symbol string) string {
	return symbol + formatDollars(dollars)
}

// formatDollars is a helper function to format dollar values with 2 decimal places
func formatDollars(dollars float64) string {
	return fmt.Sprintf("%.2f", dollars)
}

// ApplyPercentage applies a percentage to a cents value
func ApplyPercentage(cents int64, percentage float64) int64 {
	return ToCents(FromCents(cents) * percentage / 100)
}
