package cache

import (
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/patrickmn/go-cache"
	"time"
)

type GoCache struct {
	c *cache.Cache
}

func NewGoCache(defaultExpiration, cleanupInterval time.Duration) *GoCache {
	return &GoCache{
		c: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (g *GoCache) Get(orderID string) (*model.Order, bool) {
	val, found := g.c.Get(orderID)
	if !found {
		return nil, false
	}

	order, ok := val.(*model.Order)

	return order, ok
}

func (g *GoCache) Set(orderID string, order *model.Order) {
	g.c.Set(orderID, &order, cache.DefaultExpiration)
}
