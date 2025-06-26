package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
)

type orderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
}

type OrderCreatedTopicHandler struct {
	orderService orderService
}

func NewOrderCreatedTopicHandler(s orderService) *OrderCreatedTopicHandler {
	return &OrderCreatedTopicHandler{
		orderService: s,
	}
}

func (h *OrderCreatedTopicHandler) HandleMessage(ctx context.Context, msg []byte) error {
	var order model.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		return fmt.Errorf("error parsing data: %w", err)
	}

	_, err := h.orderService.CreateOrder(ctx, &order)
	return err
}
