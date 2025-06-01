package dto

import (
	"testing"
)

func TestPaginationDTO(t *testing.T) {
	pagination := PaginationDTO{
		Page:     1,
		PageSize: 20,
		Total:    100,
	}

	if pagination.Page != 1 {
		t.Errorf("Expected Page 1, got %d", pagination.Page)
	}
	if pagination.PageSize != 20 {
		t.Errorf("Expected PageSize 20, got %d", pagination.PageSize)
	}
	if pagination.Total != 100 {
		t.Errorf("Expected Total 100, got %d", pagination.Total)
	}
}

func TestResponseDTO(t *testing.T) {
	data := map[string]string{"key": "value"}
	response := ResponseDTO[map[string]string]{
		Success: true,
		Message: "Operation successful",
		Data:    data,
	}

	if !response.Success {
		t.Errorf("Expected Success true, got %t", response.Success)
	}
	if response.Message != "Operation successful" {
		t.Errorf("Expected Message 'Operation successful', got %s", response.Message)
	}
	if response.Data["key"] != "value" {
		t.Errorf("Expected Data[key] 'value', got %s", response.Data["key"])
	}
	if response.Error != "" {
		t.Errorf("Expected Error empty, got %s", response.Error)
	}
}

func TestResponseDTOWithError(t *testing.T) {
	response := ResponseDTO[string]{
		Success: false,
		Error:   "Something went wrong",
	}

	if response.Success {
		t.Errorf("Expected Success false, got %t", response.Success)
	}
	if response.Error != "Something went wrong" {
		t.Errorf("Expected Error 'Something went wrong', got %s", response.Error)
	}
	if response.Message != "" {
		t.Errorf("Expected Message empty, got %s", response.Message)
	}
}

func TestListResponseDTO(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	pagination := PaginationDTO{
		Page:     1,
		PageSize: 10,
		Total:    3,
	}

	response := ListResponseDTO[string]{
		Success:    true,
		Message:    "List retrieved successfully",
		Data:       data,
		Pagination: pagination,
	}

	if !response.Success {
		t.Errorf("Expected Success true, got %t", response.Success)
	}
	if response.Message != "List retrieved successfully" {
		t.Errorf("Expected Message 'List retrieved successfully', got %s", response.Message)
	}
	if len(response.Data) != 3 {
		t.Errorf("Expected Data length 3, got %d", len(response.Data))
	}
	if response.Data[0] != "item1" {
		t.Errorf("Expected Data[0] 'item1', got %s", response.Data[0])
	}
	if response.Pagination.Total != 3 {
		t.Errorf("Expected Pagination.Total 3, got %d", response.Pagination.Total)
	}
}

func TestListResponseDTOWithError(t *testing.T) {
	response := ListResponseDTO[string]{
		Success: false,
		Error:   "Failed to retrieve list",
	}

	if response.Success {
		t.Errorf("Expected Success false, got %t", response.Success)
	}
	if response.Error != "Failed to retrieve list" {
		t.Errorf("Expected Error 'Failed to retrieve list', got %s", response.Error)
	}
	if len(response.Data) != 0 {
		t.Errorf("Expected Data length 0, got %d", len(response.Data))
	}
}

func TestAddressDTO(t *testing.T) {
	address := AddressDTO{
		AddressLine1: "123 Main St",
		AddressLine2: "Apt 4B",
		City:         "New York",
		State:        "NY",
		PostalCode:   "10001",
		Country:      "US",
	}

	if address.AddressLine1 != "123 Main St" {
		t.Errorf("Expected AddressLine1 '123 Main St', got %s", address.AddressLine1)
	}
	if address.AddressLine2 != "Apt 4B" {
		t.Errorf("Expected AddressLine2 'Apt 4B', got %s", address.AddressLine2)
	}
	if address.City != "New York" {
		t.Errorf("Expected City 'New York', got %s", address.City)
	}
	if address.State != "NY" {
		t.Errorf("Expected State 'NY', got %s", address.State)
	}
	if address.PostalCode != "10001" {
		t.Errorf("Expected PostalCode '10001', got %s", address.PostalCode)
	}
	if address.Country != "US" {
		t.Errorf("Expected Country 'US', got %s", address.Country)
	}
}

func TestAddressDTOEmpty(t *testing.T) {
	address := AddressDTO{}

	if address.AddressLine1 != "" {
		t.Errorf("Expected AddressLine1 empty, got %s", address.AddressLine1)
	}
	if address.AddressLine2 != "" {
		t.Errorf("Expected AddressLine2 empty, got %s", address.AddressLine2)
	}
	if address.City != "" {
		t.Errorf("Expected City empty, got %s", address.City)
	}
	if address.State != "" {
		t.Errorf("Expected State empty, got %s", address.State)
	}
	if address.PostalCode != "" {
		t.Errorf("Expected PostalCode empty, got %s", address.PostalCode)
	}
	if address.Country != "" {
		t.Errorf("Expected Country empty, got %s", address.Country)
	}
}
