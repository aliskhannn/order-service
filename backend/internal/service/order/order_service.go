package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/aliskhannn/order-service/internal/model"
)

//go:generate mockgen -source=order_service.go -destination=../../mocks/service/order/mock_order.go -package=mocks
type orderRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
	GetOrderById(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Item, error)
}

type orderCache interface {
	Get(orderID uuid.UUID) (*model.Order, bool)
	Set(orderID uuid.UUID, order *model.Order)
}

// Service provides business logic for managing orders.
// It interacts with the repository for data persistence and can use a cache
// to optimize read operations. The service is responsible for creating orders,
// retrieving orders by ID, and potentially other order-related operations in the future.
type Service struct {
	cache orderCache
	repo  orderRepository
}

// New creates a new Service instance with the provided cache and repository.
func New(c orderCache, repo orderRepository) *Service {
	return &Service{
		cache: c,
		repo:  repo,
	}
}

// CreateOrder creates a new order in the system by delegating persistence to the repository.
// Returns the generated order UUID or an error if saving fails.
func (s *Service) CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	orderID, err := s.repo.SaveOrder(ctx, order)
	if err != nil {
		return uuid.Nil, fmt.Errorf("save order: %w", err)
	}

	return orderID, nil
}

// GetOrderByID retrieves an order by its UUID.
// The method first checks the cache (if enabled). If the order is not cached,
// it queries the repository, fetches the items, stores the result in the cache, and returns it.
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
