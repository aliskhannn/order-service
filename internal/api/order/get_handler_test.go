package order

import (
	"context"
	"encoding/json"
	"errors"
	mocks "github.com/aliskhannn/order-service/internal/mocks/api/order"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockorderService(ctrl)
	h := NewGetHandler(mockService)

	orderID := "b563feb7b2b84b6test"
	expectedOrder := &model.Order{OrderID: orderID}

	mockService.EXPECT().GetOrderByID(gomock.Any(), orderID).Return(expectedOrder, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", orderID)
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

	mockService := mocks.NewMockorderService(ctrl)
	h := NewGetHandler(mockService)

	orderID := "not_exist"
	expectedError := errors.New("not found")

	mockService.EXPECT().GetOrderByID(gomock.Any(), orderID).Return(nil, expectedError)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", orderID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	h.GetOrderByID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetOrderByID_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockorderService(ctrl)
	h := NewGetHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/orders/", nil)

	rr := httptest.NewRecorder()

	h.GetOrderByID(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
