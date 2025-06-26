package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/google/uuid"
)

type orderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
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
