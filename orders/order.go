package main

type Order struct {
	ID         string
	Amount     float64
	Items      []string
	CustomerID string
	Status     string
	CreatedAt  string
	UpdatedAt  string
}
