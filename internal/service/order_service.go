package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/aliskhannn/order-service/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
	GetOrderByID(ctx context.Context, orderUID uuid.UUID) (*model.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	if order == nil {
		return uuid.Nil, errors.New("order is required")
	}

	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating order: %w", err)
	}

	return orderID, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	if orderID == uuid.Nil {
		return nil, errors.New("order id is required")
	}

	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting order by id %v: %w", orderID, err)
	}

	return order, nil
}
