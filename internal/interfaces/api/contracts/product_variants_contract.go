package contracts

// SetVariantPriceRequest represents the request to set a price for a variant in a specific currency
type SetVariantPriceRequest struct {
	CurrencyCode string  `json:"currency_code"`
	Price        float64 `json:"price"`
}

// SetMultipleVariantPricesRequest represents the request to set multiple prices for a variant
type SetMultipleVariantPricesRequest struct {
	Prices map[string]float64 `json:"prices"` // currency_code -> price
}

// VariantPricesResponse represents the response containing all prices for a variant
type VariantPricesResponse struct {
	VariantID uint               `json:"variant_id"`
	Prices    map[string]float64 `json:"prices"` // currency_code -> price
}
