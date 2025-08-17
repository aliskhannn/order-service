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

type GoCache struct {
	c      *cache.Cache
	logger *zap.Logger
	repo   orderRepository
}

func New(defaultExpiration, cleanupInterval time.Duration, l *zap.Logger, r orderRepository) *GoCache {
	return &GoCache{
		c:      cache.New(defaultExpiration, cleanupInterval),
		logger: l,
		repo:   r,
	}
}

func (g *GoCache) Get(orderID uuid.UUID) (*model.Order, bool) {
	val, found := g.c.Get(orderID.String())
	if !found {
		return nil, false
	}

	order, ok := val.(*model.Order)

	return order, ok
}

func (g *GoCache) Set(orderID uuid.UUID, order *model.Order) {
	g.c.Set(orderID.String(), order, cache.DefaultExpiration)
}

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
