package cache


import (
	"sync"
	"order-service/internal/model"
)

type Cache struct {
	sync.RWMutex
	orders map[string]model.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]model.Order),
	}
}

func (c *Cache) Set(order model.Order) {
	c.Lock()
	defer c.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) Get(uid string) (model.Order, bool){
	c.RLock()
	defer c.RUnlock()
	order, exists := c.orders[uid]
	return order, exists
}

func (c *Cache) GetAll() map[string]model.Order{
	c.RLock()
	defer c.RUnlock()
	res := make(map[string]model.Order, len(c.orders))
	for k, v := range c.orders{
		res[k] = v
	}
	return res
}


func (c *Cache) Restore(orders map[string]model.Order) {
	c.Lock()
	defer c.Unlock()
	c.orders = orders
}