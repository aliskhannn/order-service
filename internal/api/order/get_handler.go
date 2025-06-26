package order

import (
	"context"
	"encoding/json"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type orderService interface {
	GetOrderByID(ctx context.Context, orderID string) (*model.Order, error)
}

type OrderHTTPHandler struct {
	orderService orderService
}

func NewOrderHTTPHandler(s orderService) *OrderHTTPHandler {
	return &OrderHTTPHandler{
		orderService: s,
	}
}

func (h *OrderHTTPHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	order, err := h.orderService.GetOrderByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
