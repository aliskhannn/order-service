package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/aliskhannn/order-service/internal/model"
)

//go:generate mockgen -source=order_created_handler.go -destination=../../../mocks/kafka/handlers/order/mock_order_created_handler.go -package=handlers orderService,validator
type orderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
}

type validator interface {
	Validate(i interface{}) error
}

// CreateHandler handles Kafka messages related to order creation.
// It validates incoming messages, invokes business logic,
// and logs results or errors.
type CreateHandler struct {
	logger       *zap.Logger
	validator    validator
	orderService orderService
}

// NewCreateHandler constructs a new CreateHandler with the given dependencies.
func NewCreateHandler(l *zap.Logger, v validator, s orderService) *CreateHandler {
	return &CreateHandler{
		logger:       l,
		validator:    v,
		orderService: s,
	}
}

// ProcessMessage processes an incoming Kafka message containing an order.
// It performs the following steps:
//  1. Unmarshals JSON payload into a model.Order.
//  2. Validates the order.
//  3. Persists the order using the service layer.
//  4. Logs success or returns a contextualized error.
//
// Returns an error if any step fails.
func (h *CreateHandler) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	var order *model.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		return fmt.Errorf("unmarshal order: %v", err)
	}

	if order == nil {
		return ErrNilOrder
	}

	if err := h.validator.Validate(order); err != nil {
		return fmt.Errorf("validate order: %v", err)
	}

	orderID, err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("create order: %v", err)
	}

	h.logger.Info("Order processed successfully", zap.String("orderID", orderID.String()))

	return nil
}
