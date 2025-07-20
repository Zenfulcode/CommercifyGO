package main

import (
	"context"

	pb "github.com/zenfulcode/commercify/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	svc OrderService
}

func NewHandler(server *grpc.Server, svc OrderService) {
	pb.RegisterOrderServiceServer(server, &grpcHandler{svc: svc})
}

func (h *grpcHandler) CreateOrder(ctx context.Context, cmd *pb.CreateOrderCommand) (*pb.CreateOrderResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (h *grpcHandler) GetOrder(ctx context.Context, query *pb.GetOrderQuery) (*pb.GetOrderResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
}
