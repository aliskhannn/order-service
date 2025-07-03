package order

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/aliskhannn/order-service/internal/model"
)

//go:generate mockgen -source=order_service.go -destination=../../mocks/service/mock_order.go -package=mocks
type orderRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
	GetOrderById(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Item, error)
}

type orderCache interface {
	Get(orderID uuid.UUID) (*model.Order, bool)
	Set(orderID uuid.UUID, order *model.Order)
}
type Service struct {
	logger *zap.Logger
	cache  orderCache
	repo   orderRepository
}

func New(l *zap.Logger, c orderCache, repo orderRepository) *Service {
	return &Service{
		logger: l,
		cache:  c,
		repo:   repo,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		s.logger.Error("failed to save order", zap.Error(err))
		return uuid.Nil, fmt.Errorf("repository error: %w", err)
	}

	s.logger.Info("recieved order id", zap.String("orderID is", orderID.String()))

	s.cache.Set(orderID, order)

	s.logger.Info("order saved and cached", zap.String("orderID", orderID.String()))

	return orderID, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	// Check cache first
	if order, found := s.cache.Get(orderID); found {
		s.logger.Info("order found in cache", zap.String("orderID", orderID.String()))
		return order, nil
	}

	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		s.logger.Error("failed to get order by ID", zap.String("orderID", orderID.String()), zap.Error(err))
		return nil, fmt.Errorf("repository error: %w", err)
	}

	items, err := s.repo.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		s.logger.Error("failed to get items by order ID", zap.String("orderID", orderID.String()), zap.Error(err))
		return nil, fmt.Errorf("repository error: %w", err)
	}

	order.Items = items

	s.cache.Set(orderID, order)

	s.logger.Info("order retrieved from repository and cached", zap.String("orderID", orderID.String()))

	return order, nil
}
