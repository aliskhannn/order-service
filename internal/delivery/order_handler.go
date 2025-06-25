package delivery

import (
	"encoding/json"
	"github.com/aliskhannn/order-service/internal/infra/kafka"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/aliskhannn/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type OrderHTTPHandler struct {
	producer     *kafka.Producer
	orderService service.OrderService
}

func NewOrderHTTPHandler(producer *kafka.Producer, s service.OrderService) *OrderHTTPHandler {
	return &OrderHTTPHandler{
		producer:     producer,
		orderService: s,
	}
}

func (h *OrderHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.producer.ProduceMessage(r.Context(), order)
	if err != nil {
		http.Error(w, "failed to send order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "order received", "order_id": order.OrderID.String()})
}

func (h *OrderHTTPHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
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
