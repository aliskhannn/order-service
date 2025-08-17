package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/aliskhannn/order-service/internal/model"
)

//go:generate mockgen -source=order_created_handler.go -destination=../../../mocks/kafka/handlers/mock_order_created_handler.go -package=handlers orderService,validator
type orderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
}

type validator interface {
	Validate(i interface{}) error
}

type CreateHandler struct {
	logger       *zap.Logger
	validator    validator
	orderService orderService
}

func NewCreateHandler(l *zap.Logger, v validator, s orderService) *CreateHandler {
	return &CreateHandler{
		logger:       l,
		validator:    v,
		orderService: s,
	}
}

func (h *CreateHandler) HandleMessage(ctx context.Context, msg []byte) error {
	var order *model.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	if order == nil {
		return fmt.Errorf("%w", ErrNilOrder)
	}

	if err := h.validator.Validate(order); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	orderID, err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateOrder, err)
	}

	h.logger.Info("Order processed successfully", zap.String("orderID", orderID.String()))

	return nil
}
