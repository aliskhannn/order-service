package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mocks "github.com/aliskhannn/order-service/internal/mocks/service/order"
	"github.com/aliskhannn/order-service/internal/model"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(mockCache, mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	order := &model.Order{
		OrderID:     orderID,
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

	mockRepo.EXPECT().SaveOrder(ctx, order).Return(orderID, nil)

	resOrderID, err := orderService.CreateOrder(ctx, order)
	assert.NoError(t, err)
	assert.Equal(t, orderID, resOrderID)
}

func TestGetOrderByID_FromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(mockCache, mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	expectedOrder := &model.Order{OrderID: orderID}

	mockCache.EXPECT().Get(orderID).Return(expectedOrder, true)

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrderByID_FromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(mockCache, mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	expectedOrder := &model.Order{OrderID: orderID}
	expectedItems := []model.Item{
		{RID: "item1"},
		{RID: "item2"},
	}

	mockCache.EXPECT().Get(orderID).Return(nil, false)
	mockRepo.EXPECT().GetOrderById(ctx, orderID).Return(expectedOrder, nil)
	mockRepo.EXPECT().GetItemsByOrderID(ctx, orderID).Return(expectedItems, nil)
	mockCache.EXPECT().Set(orderID, gomock.Any())

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	assert.Equal(t, expectedItems, order.Items)
}

func TestGetOrderById_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(mockCache, mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	expectedErr := errors.New("not found")

	mockCache.EXPECT().Get(orderID).Return(nil, false)
	mockRepo.EXPECT().GetOrderById(ctx, orderID).Return(nil, expectedErr)

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.Nil(t, order)
	assert.ErrorIs(t, err, expectedErr)
}

func TestGetOrderByID_ItemsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockorderCache(ctrl)
	mockRepo := mocks.NewMockorderRepository(ctrl)

	orderService := New(mockCache, mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedOrder := &model.Order{OrderID: orderID}
	expectedErr := errors.New("items fetch failed")

	mockCache.EXPECT().Get(orderID).Return(nil, false)
	mockRepo.EXPECT().GetOrderById(ctx, orderID).Return(expectedOrder, nil)
	mockRepo.EXPECT().GetItemsByOrderID(ctx, orderID).Return(nil, expectedErr)

	order, err := orderService.GetOrderByID(ctx, orderID)
	assert.Nil(t, order)
	assert.ErrorIs(t, err, expectedErr)
}
