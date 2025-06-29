package order

import (
	"context"
	"encoding/json"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

//go:generate mockgen -source=get_handler.go -destination=../../mocks/api/order/mock_order_service.go -package=mocks
type orderService interface {
	GetOrderByID(ctx context.Context, orderID string) (*model.Order, error)
}

type GetHandler struct {
	orderService orderService
}

func NewGetHandler(s orderService) *GetHandler {
	return &GetHandler{
		orderService: s,
	}
}

func (h *GetHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	if orderID == "" {
		http.Error(w, "order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	log.Println(order)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
