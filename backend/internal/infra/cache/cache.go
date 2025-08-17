package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/patrickmn/go-cache"

	"github.com/aliskhannn/order-service/internal/model"
)

type orderRepository interface {
	GetLastOrders(ctx context.Context, limit int) ([]model.Order, error)
}

// GoCache is an in-memory cache implementation for orders,
// built on top of github.com/patrickmn/go-cache. It supports
// basic Get/Set operations and preloading from a repository.
type GoCache struct {
	c      *cache.Cache
	logger *zap.Logger
	repo   orderRepository
}

// New creates a new GoCache instance with the given expiration
// and cleanup intervals. A zap.Logger is used for structured logging,
// and an orderRepository is required for cache preloading.
func New(defaultExpiration, cleanupInterval time.Duration, l *zap.Logger, r orderRepository) *GoCache {
	return &GoCache{
		c:      cache.New(defaultExpiration, cleanupInterval),
		logger: l,
		repo:   r,
	}
}

// Get retrieves an order from the cache by its UUID.
// It returns the order (if found and type-asserted) and a boolean flag.
func (g *GoCache) Get(orderID uuid.UUID) (*model.Order, bool) {
	val, found := g.c.Get(orderID.String())
	if !found {
		return nil, false
	}

	order, ok := val.(*model.Order)

	return order, ok
}

// Set adds an order to the cache with the default expiration time.
func (g *GoCache) Set(orderID uuid.UUID, order *model.Order) {
	g.c.Set(orderID.String(), order, cache.DefaultExpiration)
}

// Preload fetches the latest orders from the repository (up to limit)
// and loads them into the cache. This is typically used during service
// startup to warm the cache with recent data.
//
// Returns an error if fetching from the repository fails. If no orders
// are found, it logs an informational message and returns nil.
func (g *GoCache) Preload(ctx context.Context, limit int) error {
	orders, err := g.repo.GetLastOrders(ctx, limit)
	if err != nil {
		g.logger.Error("failed to preload cache", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrCachePreload, err)
	}

	if len(orders) == 0 {
		g.logger.Info("no orders found to preload cache")
		return nil
	}

	for _, order := range orders {
		o := order // Create a copy to avoid issues with the loop variable
		g.Set(order.OrderID, &o)
	}

	g.logger.Info("cache preloaded successfully", zap.Int("orders_count", len(orders)))
	return nil
}
