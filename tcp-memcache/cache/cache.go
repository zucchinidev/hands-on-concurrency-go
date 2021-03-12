package cache

import "sync"

type Cache struct {
	mu     sync.RWMutex
	values map[string]interface{}
}

func New() *Cache {
	return &Cache{
		mu:     sync.RWMutex{},
		values: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
}

func (c *Cache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if v, ok := c.values[key]; ok {
		return v
	}
	return ""
}
