package main

import (
	"log"
	"net/http"

	"github.com/zenfulcode/commercify/common"

	pb "github.com/zenfulcode/commercify/common/api"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orderServiceAddress := common.GetEnv("ORDER_SERVICE_ADDRESS", "localhost:6001")

	conn, err := grpc.NewClient(orderServiceAddress, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to gRPC server at %s", orderServiceAddress)
	orderSvcClient := pb.NewOrderServiceClient(conn)

	httpAddress := common.GetEnv("HTTP_ADDRESS", ":6091")
	mux := http.NewServeMux()
	handler := NewHandler(orderSvcClient)
	handler.registerRoutes(mux)
	log.Printf("Server is running on %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		panic(err)
	}
}
