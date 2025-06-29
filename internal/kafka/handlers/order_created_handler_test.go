package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aliskhannn/order-service/internal/mocks/kafka/handlers"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := handlers.NewMockorderService(ctrl)
	mockValidator := handlers.NewMockvalidator(ctrl)

	handler := NewOrderCreatedHandler(mockService, mockValidator)

	order := &model.Order{
		OrderID: "test-id",
	}

	msg, _ := json.Marshal(order)

	mockValidator.EXPECT().Validate(order).Return(nil)
	mockService.EXPECT().CreateOrder(gomock.Any(), order).Return("test-id", nil)

	err := handler.HandleMessage(context.Background(), msg)
	assert.NoError(t, err)
}

func TestHandleMessage_InvalidJSON(t *testing.T) {
	handler := NewOrderCreatedHandler(nil, nil)

	invalidMsg := []byte(`{invalid json}`)

	err := handler.HandleMessage(context.Background(), invalidMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestHandleMessage_NilOrder(t *testing.T) {
	handler := NewOrderCreatedHandler(nil, nil)

	nilOrder := []byte(`null`)

	err := handler.HandleMessage(context.Background(), nilOrder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order is nil")
}

func TestHandleMessage_ValidationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := handlers.NewMockvalidator(ctrl)

	order := &model.Order{OrderID: "bad"}
	msg, _ := json.Marshal(order)

	mockValidator.EXPECT().Validate(order).Return(errors.New("validation failed"))

	handler := NewOrderCreatedHandler(nil, mockValidator)

	err := handler.HandleMessage(context.Background(), msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
}

func TestHandleMessage_CreateOrderFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := handlers.NewMockvalidator(ctrl)
	mockService := handlers.NewMockorderService(ctrl)

	order := &model.Order{OrderID: "fail"}
	msg, _ := json.Marshal(order)

	mockValidator.EXPECT().Validate(order).Return(nil)
	mockService.EXPECT().CreateOrder(gomock.Any(), order).Return("", errors.New("db failure"))

	handler := NewOrderCreatedHandler(mockService, mockValidator)

	err := handler.HandleMessage(context.Background(), msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating order")
}
