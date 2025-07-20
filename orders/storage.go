package main

import "context"

type store struct {
}

func NewOrderStore() *store {
	return &store{}
}

func (s *store) Save(ctx context.Context, order Order) error {
	return nil
}
