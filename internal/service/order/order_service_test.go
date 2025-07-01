package order

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "github.com/aliskhannn/order-service/internal/mocks/service"
	"github.com/aliskhannn/order-service/internal/model"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(logger, mockCache, mockRepo)

	ctx := context.Background()
	order := &model.Order{
		OrderID:     "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    "John Doe",
			Phone:   "+1234567890",
			Zip:     "123456",
			City:    "Test City",
			Address: "123 Test St",
			Region:  "Test Region",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "txn123",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "TestProvider",
			Amount:       1000,
			PaymentDT:    1637907727,
			Bank:         "Test Bank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      123456,
				TrackNumber: "WBILMTESTTRACK",
				Price:       1000,
				RID:         "rid123",
				Name:        "Test Item",
				Sale:        30,
				Size:        "M",
				TotalPrice:  700,
				NmID:        123456789,
				Brand:       "Test Brand",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "signature",
		CustomerId:        "customer123",
		DeliveryService:   "Test Delivery Service",
		Shardkey:          "shardkey123",
		SmId:              1,
		DateCreated:       time.Now(),
		OofShard:          "oofshard123",
	}

	mockRepo.EXPECT().SaveOrder(ctx, order).Return("b563feb7b2b84b6test", nil)
	mockCache.EXPECT().Set("b563feb7b2b84b6test", order)

	orderID, err := orderService.CreateOrder(ctx, order)
	assert.NoError(t, err)
	assert.Equal(t, "b563feb7b2b84b6test", orderID)
}

func TestGetOrderByID_FromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(logger, mockCache, mockRepo)

	ctx := context.Background()
	orderID := "b563feb7b2b84b6test"
	expectedOrder := &model.Order{OrderID: orderID}

	mockCache.EXPECT().Get(orderID).Return(expectedOrder, true)

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrderByID_FromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(logger, mockCache, mockRepo)

	ctx := context.Background()
	orderID := "b563feb7b2b84b6test"
	expectedOrder := &model.Order{OrderID: orderID}

	mockCache.EXPECT().Get(orderID).Return(nil, false)
	mockRepo.EXPECT().GetOrderById(ctx, orderID).Return(expectedOrder, nil)
	mockCache.EXPECT().Set(orderID, expectedOrder)

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrderById_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewNop()
	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(logger, mockCache, mockRepo)

	ctx := context.Background()
	orderID := "b563feb7b2b84b6test"
	expectedErr := errors.New("not found")

	mockCache.EXPECT().Get(orderID).Return(nil, false)
	mockRepo.EXPECT().GetOrderById(ctx, orderID).Return(nil, expectedErr)

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.Nil(t, order)
	assert.ErrorIs(t, err, expectedErr)
}
