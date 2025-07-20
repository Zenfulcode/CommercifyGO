package main

import "context"

type service struct {
	storage OrderStorage
}

func NewOrderService(storage OrderStorage) *service {
	return &service{storage: storage}
}

func (s *service) CreateOrder(ctx context.Context, order Order) error {
	return s.storage.Save(ctx, order)
}
