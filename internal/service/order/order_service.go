package order

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/aliskhannn/order-service/internal/model"
)

//go:generate mockgen -source=order_service.go -destination=../mocks/service/mock_order.go -package=mocks
type orderRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrderById(ctx context.Context, orderID string) (*model.Order, error)
}

type orderCache interface {
	Get(orderID string) (*model.Order, bool)
	Set(orderID string, order *model.Order)
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

func (s *Service) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		s.logger.Error("failed to save order", zap.Error(err))
		return "", fmt.Errorf("repository error: %w", err)
	}

	s.cache.Set(orderID, order)

	s.logger.Info("order saved and cached", zap.String("orderID", orderID))

	return orderID, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID string) (*model.Order, error) {
	// Check cache first
	if order, found := s.cache.Get(orderID); found {
		s.logger.Info("order found in cache", zap.String("orderID", orderID))
		return order, nil
	}

	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		s.logger.Error("failed to get order by ID", zap.String("orderID", orderID), zap.Error(err))
		return nil, fmt.Errorf("repository error: %w", err)
	}

	s.cache.Set(orderID, order)

	s.logger.Info("order retrieved from repository and cached", zap.String("orderID", orderID))

	return order, nil
}
