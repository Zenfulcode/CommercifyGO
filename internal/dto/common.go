package dto

// AddressDTO represents a shipping or billing address
type AddressDTO struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
}

// CustomerDetailsDTO represents customer information for a checkout
type CustomerDetailsDTO struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}
