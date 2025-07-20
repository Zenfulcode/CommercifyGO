package main

import "context"

type OrderService interface {
	CreateOrder(context.Context, Order) error
	// GetOrder(id string) (Order, error)
	// UpdateOrder(order Order) error
	// DeleteOrder(id string) error
	// ListOrders() ([]Order, error)
}

type OrderStorage interface {
	Save(context.Context, Order) error
	// Get(context.Context, string) (Order, error)
	// Update(context.Context, Order) error
	// Delete(context.Context, string) error
	// List(context.Context) ([]Order, error)
}
