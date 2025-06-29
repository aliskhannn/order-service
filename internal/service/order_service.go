package service

import (
	"context"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"log"
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
type OrderService struct {
	repo  orderRepository
	cache orderCache
}

func NewOrderService(repo orderRepository, cache orderCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		return "", fmt.Errorf("error creating order: %w", err)
	}

	s.cache.Set(orderID, order)

	return orderID, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*model.Order, error) {
	// Check cache first
	if order, found := s.cache.Get(orderID); found {
		log.Printf("order %s found in cache", orderID)
		return order, nil
	}

	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting order by id %v: %w", orderID, err)
	}
	s.cache.Set(orderID, order)

	return order, nil
}
