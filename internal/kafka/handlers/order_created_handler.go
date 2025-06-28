package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"log"
)

type orderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
}

type validator interface {
	Validate(i interface{}) error
}

type OrderCreatedHandler struct {
	orderService orderService
	validator    validator
}

func NewOrderCreatedHandler(s orderService, v validator) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		orderService: s,
		validator:    v,
	}
}

func (h *OrderCreatedHandler) HandleMessage(ctx context.Context, msg []byte) error {
	var order *model.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		log.Printf("invalid message format: %v", err)
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if order == nil {
		log.Println("received nil order")
		return fmt.Errorf("invalid JSON: order is nil")
	}

	if err := h.validator.Validate(order); err != nil {
		log.Printf("validation failed: %v", err)
		return fmt.Errorf("validation error: %w", err)
	}

	if _, err := h.orderService.CreateOrder(ctx, order); err != nil {
		log.Printf("error creating order: %v", err)
		return fmt.Errorf("error creating order: %w", err)
	}

	return nil
}
