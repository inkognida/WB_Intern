package cache

import (
	"WB_Intern/internal/model"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	NoSuchOrder = "no such order"
	OrderAlreadyExists = "order already exists"
)

type cache struct {
	logger *logrus.Logger
	cache  map[string]*model.Order
	mu *sync.Mutex

	// repo for sync
	repo orderRepository
}

type orderRepository interface {
	GetOrders(context.Context) (map[string]*model.Order, error)
}

type Cache interface {
	SaveOrder(*model.Order) error
	GetOrderById(string) (*model.Order, error)
}

func cacheSetup(logger *logrus.Logger, repo orderRepository) map[string]*model.Order {
	newCache, err := repo.GetOrders(context.Background())

	if err != nil {
		logger.Fatalln(err)
	}

	return newCache
}

func NewCache(logger *logrus.Logger, repo orderRepository) Cache {
	return &cache{
		cache: cacheSetup(logger, repo),
		logger: logger,
		mu: &sync.Mutex{},
		repo: repo,
	}
}

func (c *cache) SaveOrder(order *model.Order) error {
	if _, ok := c.cache[order.OrderUID]; !ok {
		c.mu.Lock()
		c.cache[order.OrderUID] = order
		c.mu.Unlock()
		return nil
	}

	return errors.New(OrderAlreadyExists)
}

func (c *cache) GetOrderById(orderId string) (*model.Order, error) {
	if order, ok := c.cache[orderId]; ok {
		return order, nil
	}

	return nil, errors.New(NoSuchOrder)
}
