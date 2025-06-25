package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/aliskhannn/order-service/internal/service"
)

type OrderKafkaHandler struct {
	orderService service.OrderService
}

func NewOrderKafkaHandler(s service.OrderService) *OrderKafkaHandler {
	return &OrderKafkaHandler{
		orderService: s,
	}
}

func (h *OrderKafkaHandler) HandleMessage(msg []byte) error {
	var order model.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		return fmt.Errorf("error parsing data: %w", err)
	}

	ctx := context.Background()

	_, err := h.orderService.CreateOrder(ctx, &order)
	return err
}
