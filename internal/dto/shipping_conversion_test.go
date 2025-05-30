package dto

import (
	"testing"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

func TestConvertToShippingMethodDetailDTO(t *testing.T) {
	now := time.Now()
	method := &entity.ShippingMethod{
		ID:                    1,
		Name:                  "Express Shipping",
		Description:           "Fast delivery",
		EstimatedDeliveryDays: 2,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	dto := ConvertToShippingMethodDetailDTO(method)

	if dto.ID != method.ID {
		t.Errorf("Expected ID %d, got %d", method.ID, dto.ID)
	}
	if dto.Name != method.Name {
		t.Errorf("Expected Name %s, got %s", method.Name, dto.Name)
	}
	if dto.Description != method.Description {
		t.Errorf("Expected Description %s, got %s", method.Description, dto.Description)
	}
	if dto.EstimatedDeliveryDays != method.EstimatedDeliveryDays {
		t.Errorf("Expected EstimatedDeliveryDays %d, got %d", method.EstimatedDeliveryDays, dto.EstimatedDeliveryDays)
	}
	if dto.Active != method.Active {
		t.Errorf("Expected Active %t, got %t", method.Active, dto.Active)
	}
}

func TestConvertToShippingMethodDetailDTONil(t *testing.T) {
	dto := ConvertToShippingMethodDetailDTO(nil)

	if dto.ID != 0 {
		t.Errorf("Expected ID 0 for nil input, got %d", dto.ID)
	}
	if dto.Name != "" {
		t.Errorf("Expected empty Name for nil input, got %s", dto.Name)
	}
}

func TestConvertToShippingZoneDTO(t *testing.T) {
	now := time.Now()
	zone := &entity.ShippingZone{
		ID:          1,
		Name:        "North America",
		Description: "US and Canada",
		Countries:   []string{"US", "CA"},
		States:      []string{"NY", "CA"},
		ZipCodes:    []string{"10001", "90210"},
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	dto := ConvertToShippingZoneDTO(zone)

	if dto.ID != zone.ID {
		t.Errorf("Expected ID %d, got %d", zone.ID, dto.ID)
	}
	if dto.Name != zone.Name {
		t.Errorf("Expected Name %s, got %s", zone.Name, dto.Name)
	}
	if dto.Description != zone.Description {
		t.Errorf("Expected Description %s, got %s", zone.Description, dto.Description)
	}
	if len(dto.Countries) != len(zone.Countries) {
		t.Errorf("Expected Countries length %d, got %d", len(zone.Countries), len(dto.Countries))
	}
	if len(dto.States) != len(zone.States) {
		t.Errorf("Expected States length %d, got %d", len(zone.States), len(dto.States))
	}
	if len(dto.ZipCodes) != len(zone.ZipCodes) {
		t.Errorf("Expected ZipCodes length %d, got %d", len(zone.ZipCodes), len(dto.ZipCodes))
	}
	if dto.Active != zone.Active {
		t.Errorf("Expected Active %t, got %t", zone.Active, dto.Active)
	}
}

func TestConvertToShippingZoneDTONil(t *testing.T) {
	dto := ConvertToShippingZoneDTO(nil)

	if dto.ID != 0 {
		t.Errorf("Expected ID 0 for nil input, got %d", dto.ID)
	}
	if dto.Name != "" {
		t.Errorf("Expected empty Name for nil input, got %s", dto.Name)
	}
}

func TestConvertToShippingRateDTO(t *testing.T) {
	now := time.Now()
	freeShippingThreshold := int64(10000) // $100.00 in cents

	rate := &entity.ShippingRate{
		ID:                    1,
		ShippingMethodID:      1,
		ShippingZoneID:        1,
		BaseRate:              999,  // $9.99 in cents
		MinOrderValue:         2500, // $25.00 in cents
		FreeShippingThreshold: &freeShippingThreshold,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	dto := ConvertToShippingRateDTO(rate)

	if dto.ID != rate.ID {
		t.Errorf("Expected ID %d, got %d", rate.ID, dto.ID)
	}
	if dto.ShippingMethodID != rate.ShippingMethodID {
		t.Errorf("Expected ShippingMethodID %d, got %d", rate.ShippingMethodID, dto.ShippingMethodID)
	}
	if dto.ShippingZoneID != rate.ShippingZoneID {
		t.Errorf("Expected ShippingZoneID %d, got %d", rate.ShippingZoneID, dto.ShippingZoneID)
	}
	expectedBaseRate := money.FromCents(rate.BaseRate)
	if dto.BaseRate != expectedBaseRate {
		t.Errorf("Expected BaseRate %f, got %f", expectedBaseRate, dto.BaseRate)
	}
	expectedMinOrderValue := money.FromCents(rate.MinOrderValue)
	if dto.MinOrderValue != expectedMinOrderValue {
		t.Errorf("Expected MinOrderValue %f, got %f", expectedMinOrderValue, dto.MinOrderValue)
	}
	expectedFreeShippingThreshold := money.FromCents(*rate.FreeShippingThreshold)
	if dto.FreeShippingThreshold == nil || *dto.FreeShippingThreshold != expectedFreeShippingThreshold {
		t.Errorf("Expected FreeShippingThreshold %f, got %v", expectedFreeShippingThreshold, dto.FreeShippingThreshold)
	}
	if dto.Active != rate.Active {
		t.Errorf("Expected Active %t, got %t", rate.Active, dto.Active)
	}
}

func TestConvertToShippingRateDTONil(t *testing.T) {
	dto := ConvertToShippingRateDTO(nil)

	if dto.ID != 0 {
		t.Errorf("Expected ID 0 for nil input, got %d", dto.ID)
	}
	if dto.BaseRate != 0 {
		t.Errorf("Expected BaseRate 0 for nil input, got %f", dto.BaseRate)
	}
}

func TestConvertToWeightBasedRateDTO(t *testing.T) {
	now := time.Now()
	rate := &entity.WeightBasedRate{
		ID:             1,
		ShippingRateID: 1,
		MinWeight:      0.0,
		MaxWeight:      5.0,
		Rate:           299, // $2.99 in cents
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	dto := ConvertToWeightBasedRateDTO(rate)

	if dto.ID != rate.ID {
		t.Errorf("Expected ID %d, got %d", rate.ID, dto.ID)
	}
	if dto.ShippingRateID != rate.ShippingRateID {
		t.Errorf("Expected ShippingRateID %d, got %d", rate.ShippingRateID, dto.ShippingRateID)
	}
	if dto.MinWeight != rate.MinWeight {
		t.Errorf("Expected MinWeight %f, got %f", rate.MinWeight, dto.MinWeight)
	}
	if dto.MaxWeight != rate.MaxWeight {
		t.Errorf("Expected MaxWeight %f, got %f", rate.MaxWeight, dto.MaxWeight)
	}
	expectedRate := money.FromCents(rate.Rate)
	if dto.Rate != expectedRate {
		t.Errorf("Expected Rate %f, got %f", expectedRate, dto.Rate)
	}
}

func TestConvertToValueBasedRateDTO(t *testing.T) {
	now := time.Now()
	rate := &entity.ValueBasedRate{
		ID:             1,
		ShippingRateID: 1,
		MinOrderValue:  0,    // $0.00 in cents
		MaxOrderValue:  5000, // $50.00 in cents
		Rate:           999,  // $9.99 in cents
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	dto := ConvertToValueBasedRateDTO(rate)

	if dto.ID != rate.ID {
		t.Errorf("Expected ID %d, got %d", rate.ID, dto.ID)
	}
	if dto.ShippingRateID != rate.ShippingRateID {
		t.Errorf("Expected ShippingRateID %d, got %d", rate.ShippingRateID, dto.ShippingRateID)
	}
	expectedMinOrderValue := money.FromCents(rate.MinOrderValue)
	if dto.MinOrderValue != expectedMinOrderValue {
		t.Errorf("Expected MinOrderValue %f, got %f", expectedMinOrderValue, dto.MinOrderValue)
	}
	expectedMaxOrderValue := money.FromCents(rate.MaxOrderValue)
	if dto.MaxOrderValue != expectedMaxOrderValue {
		t.Errorf("Expected MaxOrderValue %f, got %f", expectedMaxOrderValue, dto.MaxOrderValue)
	}
	expectedRate := money.FromCents(rate.Rate)
	if dto.Rate != expectedRate {
		t.Errorf("Expected Rate %f, got %f", expectedRate, dto.Rate)
	}
}

func TestConvertToShippingOptionDTO(t *testing.T) {
	option := &entity.ShippingOption{
		ShippingRateID:        1,
		ShippingMethodID:      1,
		Name:                  "Standard Shipping",
		Description:           "5-7 business days",
		EstimatedDeliveryDays: 6,
		Cost:                  999, // $9.99 in cents
		FreeShipping:          false,
	}

	dto := ConvertToShippingOptionDTO(option)

	if dto.ShippingRateID != option.ShippingRateID {
		t.Errorf("Expected ShippingRateID %d, got %d", option.ShippingRateID, dto.ShippingRateID)
	}
	if dto.ShippingMethodID != option.ShippingMethodID {
		t.Errorf("Expected ShippingMethodID %d, got %d", option.ShippingMethodID, dto.ShippingMethodID)
	}
	if dto.Name != option.Name {
		t.Errorf("Expected Name %s, got %s", option.Name, dto.Name)
	}
	if dto.Description != option.Description {
		t.Errorf("Expected Description %s, got %s", option.Description, dto.Description)
	}
	if dto.EstimatedDeliveryDays != option.EstimatedDeliveryDays {
		t.Errorf("Expected EstimatedDeliveryDays %d, got %d", option.EstimatedDeliveryDays, dto.EstimatedDeliveryDays)
	}
	expectedCost := money.FromCents(option.Cost)
	if dto.Cost != expectedCost {
		t.Errorf("Expected Cost %f, got %f", expectedCost, dto.Cost)
	}
	if dto.FreeShipping != option.FreeShipping {
		t.Errorf("Expected FreeShipping %t, got %t", option.FreeShipping, dto.FreeShipping)
	}
}

func TestCreateShippingMethodRequestToUseCaseInput(t *testing.T) {
	request := CreateShippingMethodRequest{
		Name:                  "Express Shipping",
		Description:           "Fast delivery",
		EstimatedDeliveryDays: 2,
	}

	input := request.ToCreateShippingMethodInput()

	if input.Name != request.Name {
		t.Errorf("Expected Name %s, got %s", request.Name, input.Name)
	}
	if input.Description != request.Description {
		t.Errorf("Expected Description %s, got %s", request.Description, input.Description)
	}
	if input.EstimatedDeliveryDays != request.EstimatedDeliveryDays {
		t.Errorf("Expected EstimatedDeliveryDays %d, got %d", request.EstimatedDeliveryDays, input.EstimatedDeliveryDays)
	}
}

func TestUpdateShippingMethodRequestToUseCaseInput(t *testing.T) {
	id := uint(1)
	request := UpdateShippingMethodRequest{
		Name:                  "Updated Express",
		Description:           "Updated description",
		EstimatedDeliveryDays: 3,
		Active:                false,
	}

	input := request.ToUpdateShippingMethodInput(id)

	if input.ID != id {
		t.Errorf("Expected ID %d, got %d", id, input.ID)
	}
	if input.Name != request.Name {
		t.Errorf("Expected Name %s, got %s", request.Name, input.Name)
	}
	if input.Description != request.Description {
		t.Errorf("Expected Description %s, got %s", request.Description, input.Description)
	}
	if input.EstimatedDeliveryDays != request.EstimatedDeliveryDays {
		t.Errorf("Expected EstimatedDeliveryDays %d, got %d", request.EstimatedDeliveryDays, input.EstimatedDeliveryDays)
	}
	if input.Active != request.Active {
		t.Errorf("Expected Active %t, got %t", request.Active, input.Active)
	}
}

func TestCreateShippingZoneRequestToUseCaseInput(t *testing.T) {
	request := CreateShippingZoneRequest{
		Name:        "Europe",
		Description: "European countries",
		Countries:   []string{"DE", "FR"},
		States:      []string{"Bavaria"},
		ZipCodes:    []string{"80331"},
	}

	input := request.ToCreateShippingZoneInput()

	if input.Name != request.Name {
		t.Errorf("Expected Name %s, got %s", request.Name, input.Name)
	}
	if input.Description != request.Description {
		t.Errorf("Expected Description %s, got %s", request.Description, input.Description)
	}
	if len(input.Countries) != len(request.Countries) {
		t.Errorf("Expected Countries length %d, got %d", len(request.Countries), len(input.Countries))
	}
	if len(input.States) != len(request.States) {
		t.Errorf("Expected States length %d, got %d", len(request.States), len(input.States))
	}
	if len(input.ZipCodes) != len(request.ZipCodes) {
		t.Errorf("Expected ZipCodes length %d, got %d", len(request.ZipCodes), len(input.ZipCodes))
	}
}

func TestAddressDTOToEntityAddress(t *testing.T) {
	dto := AddressDTO{
		AddressLine1: "123 Main St",
		AddressLine2: "Apt 4B",
		City:         "New York",
		State:        "NY",
		PostalCode:   "10001",
		Country:      "US",
	}

	address := dto.ToEntityAddress()

	if address.Street != dto.AddressLine1 {
		t.Errorf("Expected Street %s, got %s", dto.AddressLine1, address.Street)
	}
	if address.City != dto.City {
		t.Errorf("Expected City %s, got %s", dto.City, address.City)
	}
	if address.State != dto.State {
		t.Errorf("Expected State %s, got %s", dto.State, address.State)
	}
	if address.PostalCode != dto.PostalCode {
		t.Errorf("Expected PostalCode %s, got %s", dto.PostalCode, address.PostalCode)
	}
	if address.Country != dto.Country {
		t.Errorf("Expected Country %s, got %s", dto.Country, address.Country)
	}
}

func TestAddressDTOToDomainAddress(t *testing.T) {
	dto := AddressDTO{
		AddressLine1: "456 Oak Ave",
		City:         "Boston",
		State:        "MA",
		PostalCode:   "02101",
		Country:      "US",
	}

	address := dto.ToDomainAddress()

	// Should be the same as ToEntityAddress
	if address.Street != dto.AddressLine1 {
		t.Errorf("Expected Street %s, got %s", dto.AddressLine1, address.Street)
	}
	if address.City != dto.City {
		t.Errorf("Expected City %s, got %s", dto.City, address.City)
	}
}

func TestConvertShippingMethodListToDTO(t *testing.T) {
	now := time.Now()
	methods := []*entity.ShippingMethod{
		{
			ID:                    1,
			Name:                  "Standard",
			EstimatedDeliveryDays: 5,
			Active:                true,
			CreatedAt:             now,
			UpdatedAt:             now,
		},
		{
			ID:                    2,
			Name:                  "Express",
			EstimatedDeliveryDays: 2,
			Active:                true,
			CreatedAt:             now,
			UpdatedAt:             now,
		},
	}

	dtos := ConvertShippingMethodListToDTO(methods)

	if len(dtos) != len(methods) {
		t.Errorf("Expected DTOs length %d, got %d", len(methods), len(dtos))
	}
	if dtos[0].Name != "Standard" {
		t.Errorf("Expected first DTO Name 'Standard', got %s", dtos[0].Name)
	}
	if dtos[1].Name != "Express" {
		t.Errorf("Expected second DTO Name 'Express', got %s", dtos[1].Name)
	}
}

func TestConvertShippingZoneListToDTO(t *testing.T) {
	now := time.Now()
	zones := []*entity.ShippingZone{
		{
			ID:        1,
			Name:      "Domestic",
			Countries: []string{"US"},
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "International",
			Countries: []string{"CA", "MX"},
			Active:    true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	dtos := ConvertShippingZoneListToDTO(zones)

	if len(dtos) != len(zones) {
		t.Errorf("Expected DTOs length %d, got %d", len(zones), len(dtos))
	}
	if dtos[0].Name != "Domestic" {
		t.Errorf("Expected first DTO Name 'Domestic', got %s", dtos[0].Name)
	}
	if dtos[1].Name != "International" {
		t.Errorf("Expected second DTO Name 'International', got %s", dtos[1].Name)
	}
}

func TestConvertShippingRateListToDTO(t *testing.T) {
	now := time.Now()
	rates := []*entity.ShippingRate{
		{
			ID:               1,
			ShippingMethodID: 1,
			ShippingZoneID:   1,
			BaseRate:         999,
			Active:           true,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               2,
			ShippingMethodID: 2,
			ShippingZoneID:   1,
			BaseRate:         1999,
			Active:           true,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	dtos := ConvertShippingRateListToDTO(rates)

	if len(dtos) != len(rates) {
		t.Errorf("Expected DTOs length %d, got %d", len(rates), len(dtos))
	}
	expectedFirstRate := money.FromCents(rates[0].BaseRate)
	if dtos[0].BaseRate != expectedFirstRate {
		t.Errorf("Expected first DTO BaseRate %f, got %f", expectedFirstRate, dtos[0].BaseRate)
	}
	expectedSecondRate := money.FromCents(rates[1].BaseRate)
	if dtos[1].BaseRate != expectedSecondRate {
		t.Errorf("Expected second DTO BaseRate %f, got %f", expectedSecondRate, dtos[1].BaseRate)
	}
}

func TestConvertShippingOptionListToDTO(t *testing.T) {
	options := []*entity.ShippingOption{
		{
			ShippingRateID:   1,
			ShippingMethodID: 1,
			Name:             "Standard",
			Cost:             999,
			FreeShipping:     false,
		},
		{
			ShippingRateID:   2,
			ShippingMethodID: 2,
			Name:             "Express",
			Cost:             1999,
			FreeShipping:     false,
		},
	}

	dtos := ConvertShippingOptionListToDTO(options)

	if len(dtos) != len(options) {
		t.Errorf("Expected DTOs length %d, got %d", len(options), len(dtos))
	}
	if dtos[0].Name != "Standard" {
		t.Errorf("Expected first DTO Name 'Standard', got %s", dtos[0].Name)
	}
	if dtos[1].Name != "Express" {
		t.Errorf("Expected second DTO Name 'Express', got %s", dtos[1].Name)
	}
	expectedFirstCost := money.FromCents(options[0].Cost)
	if dtos[0].Cost != expectedFirstCost {
		t.Errorf("Expected first DTO Cost %f, got %f", expectedFirstCost, dtos[0].Cost)
	}
}
