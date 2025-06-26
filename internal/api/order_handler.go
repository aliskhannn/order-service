package api

import (
	"context"
	"encoding/json"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type orderService interface {
	GetOrderByID(ctx context.Context, orderUID uuid.UUID) (*model.Order, error)
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
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
