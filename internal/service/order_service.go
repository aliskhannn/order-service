package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
)

type orderRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrderById(ctx context.Context, orderID string) (*model.Order, error)
}
type OrderService struct {
	repo orderRepository
}

func NewOrderService(repo orderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	if order == nil {
		return "", errors.New("order is required")
	}

	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		return "", fmt.Errorf("error creating order: %w", err)
	}

	return orderID, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*model.Order, error) {
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting order by id %v: %w", orderID, err)
	}

	return order, nil
}
