package order

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	handlers "github.com/aliskhannn/order-service/internal/mocks/kafka/handlers/order"
	"github.com/aliskhannn/order-service/internal/model"
)

func TestHandleMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockService := handlers.NewMockorderService(ctrl)
	mockValidator := handlers.NewMockvalidator(ctrl)

	handler := NewCreateHandler(logger, mockValidator, mockService)

	orderID := uuid.New()
	order := &model.Order{
		OrderID: orderID,
	}

	msg, _ := json.Marshal(order)

	mockValidator.EXPECT().Validate(order).Return(nil)
	mockService.EXPECT().CreateOrder(gomock.Any(), order).Return(orderID, nil)

	err := handler.HandleMessage(context.Background(), msg)
	assert.NoError(t, err)
}

func TestHandleMessage_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	handler := NewCreateHandler(logger, nil, nil)

	invalidMsg := []byte(`{invalid json}`)

	err := handler.HandleMessage(context.Background(), invalidMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestHandleMessage_NilOrder(t *testing.T) {
	logger := zap.NewNop()
	handler := NewCreateHandler(logger, nil, nil)

	nilOrder := []byte(`null`)

	err := handler.HandleMessage(context.Background(), nilOrder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil order")
}

func TestHandleMessage_ValidationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockValidator := handlers.NewMockvalidator(ctrl)

	order := &model.Order{OrderID: uuid.New()}
	msg, _ := json.Marshal(order)

	mockValidator.EXPECT().Validate(order).Return(ErrValidation)

	handler := NewCreateHandler(logger, mockValidator, nil)

	err := handler.HandleMessage(context.Background(), msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate order")
}

func TestHandleMessage_CreateOrderFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockValidator := handlers.NewMockvalidator(ctrl)
	mockService := handlers.NewMockorderService(ctrl)

	order := &model.Order{OrderID: uuid.New()}
	msg, _ := json.Marshal(order)

	mockValidator.EXPECT().Validate(order).Return(nil)
	mockService.EXPECT().CreateOrder(gomock.Any(), order).Return(uuid.UUID{}, errors.New("db failure"))

	handler := NewCreateHandler(logger, mockValidator, mockService)

	err := handler.HandleMessage(context.Background(), msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create order")
}
