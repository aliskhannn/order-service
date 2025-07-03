package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"

	"go.uber.org/zap"

	customerr "github.com/aliskhannn/order-service/internal/errors"
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
		h.logger.Warn("ivalid JSON format", zap.Error(err))
		return fmt.Errorf("%w: %v", customerr.ErrInvalidJSON, err)
	}

	if order == nil {
		h.logger.Warn("received nil order in message")
		return customerr.ErrNilOrder
	}

	if err := h.validator.Validate(order); err != nil {
		h.logger.Warn("validation error", zap.Error(err))
		return fmt.Errorf("%w: %v", customerr.ErrValidation, err)
	}

	if _, err := h.orderService.CreateOrder(ctx, order); err != nil {
		h.logger.Error("failed to create order", zap.String("orderID", order.OrderID.String()), zap.Error(err))
		return fmt.Errorf("%w: %v", customerr.ErrCreateOrder, err)
	}

	h.logger.Info("order created successfully", zap.String("orderID", order.OrderID.String()))
	return nil
}
