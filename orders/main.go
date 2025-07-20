package main

import (
	"context"
	"log"
	"net"

	"github.com/zenfulcode/commercify/common"

	"google.golang.org/grpc"
)

func main() {
	grpcAddress := common.GetEnv("GRPC_ADDRESS", ":6001")

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
	defer listener.Close()

	store := NewOrderStore()
	service := NewOrderService(store)

	NewHandler(grpcServer, service)

	log.Printf("gRPC server is running on %s", grpcAddress)

	o := Order{
		ID:         "123",
		Amount:     100.0,
		Items:      []string{"item1", "item2"},
		CustomerID: "cust1",
		Status:     "pending",
		CreatedAt:  "2023-10-01T12:00:00Z",
		UpdatedAt:  "2023-10-01T12:00:00Z",
	}

	service.CreateOrder(context.Background(), o)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
