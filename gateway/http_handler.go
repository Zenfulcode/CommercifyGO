package main

import (
	"net/http"

	"github.com/zenfulcode/commercify/common"
	pb "github.com/zenfulcode/commercify/common/api"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{
		client,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/orders", h.handleCreateOrder)
}

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var items []*pb.OrderItem
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// todo: Validate items and other fields as necessary

	result, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderCommand{
		CustomerId: "12345",
		Items:      items,
	})
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusCreated, result)
}
