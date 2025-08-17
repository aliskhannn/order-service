package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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
	cache orderCache
	repo  orderRepository
}

func New(c orderCache, repo orderRepository) *Service {
	return &Service{
		cache: c,
		repo:  repo,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		return uuid.Nil, fmt.Errorf("save order: %w", err)
	}

	//s.cache.Set(orderID, order)

	return orderID, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	// Check cache first
	if s.cache != nil {
		if order, found := s.cache.Get(orderID); found {
			return order, nil
		}
	}

	order, err := s.repo.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	items, err := s.repo.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get items by order id: %w", err)
	}

	order.Items = items

	s.cache.Set(orderID, order)

	return order, nil
}
