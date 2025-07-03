package order

import (
	"context"
	"encoding/json"
	customerr "github.com/aliskhannn/order-service/internal/errors"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "github.com/aliskhannn/order-service/internal/mocks/api/order"
	"github.com/aliskhannn/order-service/internal/model"
)

func TestGetOrderByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockService := mocks.NewMockorderService(ctrl)
	h := NewGetHandler(logger, mockService)

	orderID := uuid.New()
	expectedOrder := &model.Order{OrderID: orderID}

	mockService.EXPECT().GetOrderByID(gomock.Any(), orderID).Return(expectedOrder, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", orderID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	h.GetOrderByID(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)

	var order model.Order
	err := json.NewDecoder(rr.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.OrderID, order.OrderID)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockService := mocks.NewMockorderService(ctrl)
	h := NewGetHandler(logger, mockService)

	orderID := uuid.New()

	mockService.EXPECT().GetOrderByID(gomock.Any(), orderID).Return(nil, customerr.ErrOrderNotFound)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", orderID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	h.GetOrderByID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetOrderByID_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockService := mocks.NewMockorderService(ctrl)
	h := NewGetHandler(logger, mockService)

	req := httptest.NewRequest(http.MethodGet, "/orders/", nil)

	rr := httptest.NewRecorder()

	h.GetOrderByID(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
