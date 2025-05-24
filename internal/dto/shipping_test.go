package dto

import (
	"testing"
	"time"
)

func TestShippingMethodDetailDTO(t *testing.T) {
	now := time.Now()
	method := ShippingMethodDetailDTO{
		ID:                    1,
		Name:                  "Standard Shipping",
		Description:           "Standard delivery service",
		EstimatedDeliveryDays: 5,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if method.ID != 1 {
		t.Errorf("Expected ID 1, got %d", method.ID)
	}
	if method.Name != "Standard Shipping" {
		t.Errorf("Expected Name 'Standard Shipping', got %s", method.Name)
	}
	if method.Description != "Standard delivery service" {
		t.Errorf("Expected Description 'Standard delivery service', got %s", method.Description)
	}
	if method.EstimatedDeliveryDays != 5 {
		t.Errorf("Expected EstimatedDeliveryDays 5, got %d", method.EstimatedDeliveryDays)
	}
	if !method.Active {
		t.Errorf("Expected Active true, got %t", method.Active)
	}
}

func TestCreateShippingMethodRequest(t *testing.T) {
	request := CreateShippingMethodRequest{
		Name:                  "Express Shipping",
		Description:           "Fast delivery service",
		EstimatedDeliveryDays: 2,
	}

	if request.Name != "Express Shipping" {
		t.Errorf("Expected Name 'Express Shipping', got %s", request.Name)
	}
	if request.Description != "Fast delivery service" {
		t.Errorf("Expected Description 'Fast delivery service', got %s", request.Description)
	}
	if request.EstimatedDeliveryDays != 2 {
		t.Errorf("Expected EstimatedDeliveryDays 2, got %d", request.EstimatedDeliveryDays)
	}
}

func TestUpdateShippingMethodRequest(t *testing.T) {
	request := UpdateShippingMethodRequest{
		Name:                  "Updated Express",
		Description:           "Updated description",
		EstimatedDeliveryDays: 3,
		Active:                false,
	}

	if request.Name != "Updated Express" {
		t.Errorf("Expected Name 'Updated Express', got %s", request.Name)
	}
	if request.Description != "Updated description" {
		t.Errorf("Expected Description 'Updated description', got %s", request.Description)
	}
	if request.EstimatedDeliveryDays != 3 {
		t.Errorf("Expected EstimatedDeliveryDays 3, got %d", request.EstimatedDeliveryDays)
	}
	if request.Active {
		t.Errorf("Expected Active false, got %t", request.Active)
	}
}

func TestShippingZoneDTO(t *testing.T) {
	now := time.Now()
	zone := ShippingZoneDTO{
		ID:          1,
		Name:        "North America",
		Description: "United States and Canada",
		Countries:   []string{"US", "CA"},
		States:      []string{"NY", "CA", "ON", "BC"},
		ZipCodes:    []string{"10001", "90210", "M5V3A8"},
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if zone.ID != 1 {
		t.Errorf("Expected ID 1, got %d", zone.ID)
	}
	if zone.Name != "North America" {
		t.Errorf("Expected Name 'North America', got %s", zone.Name)
	}
	if zone.Description != "United States and Canada" {
		t.Errorf("Expected Description 'United States and Canada', got %s", zone.Description)
	}
	if len(zone.Countries) != 2 {
		t.Errorf("Expected Countries length 2, got %d", len(zone.Countries))
	}
	if zone.Countries[0] != "US" {
		t.Errorf("Expected Countries[0] 'US', got %s", zone.Countries[0])
	}
	if len(zone.States) != 4 {
		t.Errorf("Expected States length 4, got %d", len(zone.States))
	}
	if len(zone.ZipCodes) != 3 {
		t.Errorf("Expected ZipCodes length 3, got %d", len(zone.ZipCodes))
	}
	if !zone.Active {
		t.Errorf("Expected Active true, got %t", zone.Active)
	}
}

func TestCreateShippingZoneRequest(t *testing.T) {
	request := CreateShippingZoneRequest{
		Name:        "Europe",
		Description: "European Union countries",
		Countries:   []string{"DE", "FR", "IT"},
		States:      []string{"Bavaria", "Ile-de-France"},
		ZipCodes:    []string{"80331", "75001"},
	}

	if request.Name != "Europe" {
		t.Errorf("Expected Name 'Europe', got %s", request.Name)
	}
	if request.Description != "European Union countries" {
		t.Errorf("Expected Description 'European Union countries', got %s", request.Description)
	}
	if len(request.Countries) != 3 {
		t.Errorf("Expected Countries length 3, got %d", len(request.Countries))
	}
	if request.Countries[0] != "DE" {
		t.Errorf("Expected Countries[0] 'DE', got %s", request.Countries[0])
	}
	if len(request.States) != 2 {
		t.Errorf("Expected States length 2, got %d", len(request.States))
	}
	if len(request.ZipCodes) != 2 {
		t.Errorf("Expected ZipCodes length 2, got %d", len(request.ZipCodes))
	}
}

func TestUpdateShippingZoneRequest(t *testing.T) {
	request := UpdateShippingZoneRequest{
		Name:        "Updated Europe",
		Description: "Updated description",
		Countries:   []string{"DE", "FR"},
		States:      []string{"Bavaria"},
		ZipCodes:    []string{"80331"},
		Active:      false,
	}

	if request.Name != "Updated Europe" {
		t.Errorf("Expected Name 'Updated Europe', got %s", request.Name)
	}
	if request.Description != "Updated description" {
		t.Errorf("Expected Description 'Updated description', got %s", request.Description)
	}
	if len(request.Countries) != 2 {
		t.Errorf("Expected Countries length 2, got %d", len(request.Countries))
	}
	if request.Active {
		t.Error("Expected Active false, got true")
	}
}

func TestShippingRateDTO(t *testing.T) {
	now := time.Now()
	freeShippingThreshold := 100.0

	shippingMethod := &ShippingMethodDetailDTO{
		ID:   1,
		Name: "Standard",
	}

	shippingZone := &ShippingZoneDTO{
		ID:   1,
		Name: "Domestic",
	}

	rate := ShippingRateDTO{
		ID:                    1,
		ShippingMethodID:      1,
		ShippingMethod:        shippingMethod,
		ShippingZoneID:        1,
		ShippingZone:          shippingZone,
		BaseRate:              9.99,
		MinOrderValue:         25.00,
		FreeShippingThreshold: &freeShippingThreshold,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if rate.ID != 1 {
		t.Errorf("Expected ID 1, got %d", rate.ID)
	}
	if rate.ShippingMethodID != 1 {
		t.Errorf("Expected ShippingMethodID 1, got %d", rate.ShippingMethodID)
	}
	if rate.ShippingZoneID != 1 {
		t.Errorf("Expected ShippingZoneID 1, got %d", rate.ShippingZoneID)
	}
	if rate.BaseRate != 9.99 {
		t.Errorf("Expected BaseRate 9.99, got %f", rate.BaseRate)
	}
	if rate.MinOrderValue != 25.00 {
		t.Errorf("Expected MinOrderValue 25.00, got %f", rate.MinOrderValue)
	}
	if rate.FreeShippingThreshold == nil || *rate.FreeShippingThreshold != 100.0 {
		t.Errorf("Expected FreeShippingThreshold 100.0, got %v", rate.FreeShippingThreshold)
	}
	if !rate.Active {
		t.Errorf("Expected Active true, got %t", rate.Active)
	}
	if rate.ShippingMethod.Name != "Standard" {
		t.Errorf("Expected ShippingMethod.Name 'Standard', got %s", rate.ShippingMethod.Name)
	}
	if rate.ShippingZone.Name != "Domestic" {
		t.Errorf("Expected ShippingZone.Name 'Domestic', got %s", rate.ShippingZone.Name)
	}
}

func TestCreateShippingRateRequest(t *testing.T) {
	freeShippingThreshold := 75.0

	request := CreateShippingRateRequest{
		ShippingMethodID:      1,
		ShippingZoneID:        2,
		BaseRate:              12.99,
		MinOrderValue:         30.00,
		FreeShippingThreshold: &freeShippingThreshold,
		Active:                true,
	}

	if request.ShippingMethodID != 1 {
		t.Errorf("Expected ShippingMethodID 1, got %d", request.ShippingMethodID)
	}
	if request.ShippingZoneID != 2 {
		t.Errorf("Expected ShippingZoneID 2, got %d", request.ShippingZoneID)
	}
	if request.BaseRate != 12.99 {
		t.Errorf("Expected BaseRate 12.99, got %f", request.BaseRate)
	}
	if request.MinOrderValue != 30.00 {
		t.Errorf("Expected MinOrderValue 30.00, got %f", request.MinOrderValue)
	}
	if request.FreeShippingThreshold == nil || *request.FreeShippingThreshold != 75.0 {
		t.Errorf("Expected FreeShippingThreshold 75.0, got %v", request.FreeShippingThreshold)
	}
	if !request.Active {
		t.Errorf("Expected Active true, got %t", request.Active)
	}
}

func TestUpdateShippingRateRequest(t *testing.T) {
	freeShippingThreshold := 50.0

	request := UpdateShippingRateRequest{
		BaseRate:              8.99,
		MinOrderValue:         20.00,
		FreeShippingThreshold: &freeShippingThreshold,
		Active:                false,
	}

	if request.BaseRate != 8.99 {
		t.Errorf("Expected BaseRate 8.99, got %f", request.BaseRate)
	}
	if request.MinOrderValue != 20.00 {
		t.Errorf("Expected MinOrderValue 20.00, got %f", request.MinOrderValue)
	}
	if request.FreeShippingThreshold == nil || *request.FreeShippingThreshold != 50.0 {
		t.Errorf("Expected FreeShippingThreshold 50.0, got %v", request.FreeShippingThreshold)
	}
	if request.Active {
		t.Errorf("Expected Active false, got %t", request.Active)
	}
}

func TestWeightBasedRateDTO(t *testing.T) {
	now := time.Now()
	rate := WeightBasedRateDTO{
		ID:             1,
		ShippingRateID: 1,
		MinWeight:      0.0,
		MaxWeight:      5.0,
		Rate:           2.99,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if rate.ID != 1 {
		t.Errorf("Expected ID 1, got %d", rate.ID)
	}
	if rate.ShippingRateID != 1 {
		t.Errorf("Expected ShippingRateID 1, got %d", rate.ShippingRateID)
	}
	if rate.MinWeight != 0.0 {
		t.Errorf("Expected MinWeight 0.0, got %f", rate.MinWeight)
	}
	if rate.MaxWeight != 5.0 {
		t.Errorf("Expected MaxWeight 5.0, got %f", rate.MaxWeight)
	}
	if rate.Rate != 2.99 {
		t.Errorf("Expected Rate 2.99, got %f", rate.Rate)
	}
}

func TestCreateWeightBasedRateRequest(t *testing.T) {
	request := CreateWeightBasedRateRequest{
		ShippingRateID: 2,
		MinWeight:      5.0,
		MaxWeight:      10.0,
		Rate:           5.99,
	}

	if request.ShippingRateID != 2 {
		t.Errorf("Expected ShippingRateID 2, got %d", request.ShippingRateID)
	}
	if request.MinWeight != 5.0 {
		t.Errorf("Expected MinWeight 5.0, got %f", request.MinWeight)
	}
	if request.MaxWeight != 10.0 {
		t.Errorf("Expected MaxWeight 10.0, got %f", request.MaxWeight)
	}
	if request.Rate != 5.99 {
		t.Errorf("Expected Rate 5.99, got %f", request.Rate)
	}
}

func TestValueBasedRateDTO(t *testing.T) {
	now := time.Now()
	rate := ValueBasedRateDTO{
		ID:             1,
		ShippingRateID: 1,
		MinOrderValue:  0.0,
		MaxOrderValue:  50.0,
		Rate:           9.99,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if rate.ID != 1 {
		t.Errorf("Expected ID 1, got %d", rate.ID)
	}
	if rate.ShippingRateID != 1 {
		t.Errorf("Expected ShippingRateID 1, got %d", rate.ShippingRateID)
	}
	if rate.MinOrderValue != 0.0 {
		t.Errorf("Expected MinOrderValue 0.0, got %f", rate.MinOrderValue)
	}
	if rate.MaxOrderValue != 50.0 {
		t.Errorf("Expected MaxOrderValue 50.0, got %f", rate.MaxOrderValue)
	}
	if rate.Rate != 9.99 {
		t.Errorf("Expected Rate 9.99, got %f", rate.Rate)
	}
}

func TestCreateValueBasedRateRequest(t *testing.T) {
	request := CreateValueBasedRateRequest{
		ShippingRateID: 2,
		MinOrderValue:  50.0,
		MaxOrderValue:  100.0,
		Rate:           7.99,
	}

	if request.ShippingRateID != 2 {
		t.Errorf("Expected ShippingRateID 2, got %d", request.ShippingRateID)
	}
	if request.MinOrderValue != 50.0 {
		t.Errorf("Expected MinOrderValue 50.0, got %f", request.MinOrderValue)
	}
	if request.MaxOrderValue != 100.0 {
		t.Errorf("Expected MaxOrderValue 100.0, got %f", request.MaxOrderValue)
	}
	if request.Rate != 7.99 {
		t.Errorf("Expected Rate 7.99, got %f", request.Rate)
	}
}

func TestShippingOptionDTO(t *testing.T) {
	option := ShippingOptionDTO{
		ShippingRateID:        1,
		ShippingMethodID:      1,
		Name:                  "Standard Shipping",
		Description:           "5-7 business days",
		EstimatedDeliveryDays: 6,
		Cost:                  9.99,
		FreeShipping:          false,
	}

	if option.ShippingRateID != 1 {
		t.Errorf("Expected ShippingRateID 1, got %d", option.ShippingRateID)
	}
	if option.ShippingMethodID != 1 {
		t.Errorf("Expected ShippingMethodID 1, got %d", option.ShippingMethodID)
	}
	if option.Name != "Standard Shipping" {
		t.Errorf("Expected Name 'Standard Shipping', got %s", option.Name)
	}
	if option.Description != "5-7 business days" {
		t.Errorf("Expected Description '5-7 business days', got %s", option.Description)
	}
	if option.EstimatedDeliveryDays != 6 {
		t.Errorf("Expected EstimatedDeliveryDays 6, got %d", option.EstimatedDeliveryDays)
	}
	if option.Cost != 9.99 {
		t.Errorf("Expected Cost 9.99, got %f", option.Cost)
	}
	if option.FreeShipping {
		t.Errorf("Expected FreeShipping false, got %t", option.FreeShipping)
	}
}

func TestCalculateShippingOptionsRequest(t *testing.T) {
	address := AddressDTO{
		AddressLine1: "123 Test St",
		City:         "Test City",
		State:        "TS",
		PostalCode:   "12345",
		Country:      "US",
	}

	request := CalculateShippingOptionsRequest{
		Address:     address,
		OrderValue:  99.99,
		OrderWeight: 2.5,
	}

	if request.OrderValue != 99.99 {
		t.Errorf("Expected OrderValue 99.99, got %f", request.OrderValue)
	}
	if request.OrderWeight != 2.5 {
		t.Errorf("Expected OrderWeight 2.5, got %f", request.OrderWeight)
	}
	if request.Address.City != "Test City" {
		t.Errorf("Expected Address.City 'Test City', got %s", request.Address.City)
	}
}

func TestCalculateShippingOptionsResponse(t *testing.T) {
	options := []ShippingOptionDTO{
		{
			ShippingRateID: 1,
			Name:           "Standard",
			Cost:           9.99,
		},
		{
			ShippingRateID: 2,
			Name:           "Express",
			Cost:           19.99,
		},
	}

	response := CalculateShippingOptionsResponse{
		Options: options,
	}

	if len(response.Options) != 2 {
		t.Errorf("Expected Options length 2, got %d", len(response.Options))
	}
	if response.Options[0].Name != "Standard" {
		t.Errorf("Expected Options[0].Name 'Standard', got %s", response.Options[0].Name)
	}
	if response.Options[1].Cost != 19.99 {
		t.Errorf("Expected Options[1].Cost 19.99, got %f", response.Options[1].Cost)
	}
}

func TestCalculateShippingCostRequest(t *testing.T) {
	request := CalculateShippingCostRequest{
		OrderValue:  75.50,
		OrderWeight: 3.2,
	}

	if request.OrderValue != 75.50 {
		t.Errorf("Expected OrderValue 75.50, got %f", request.OrderValue)
	}
	if request.OrderWeight != 3.2 {
		t.Errorf("Expected OrderWeight 3.2, got %f", request.OrderWeight)
	}
}

func TestCalculateShippingCostResponse(t *testing.T) {
	response := CalculateShippingCostResponse{
		Cost: 12.99,
	}

	if response.Cost != 12.99 {
		t.Errorf("Expected Cost 12.99, got %f", response.Cost)
	}
}
