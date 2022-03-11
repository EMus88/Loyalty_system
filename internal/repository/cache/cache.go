package cache

import (
	"fmt"
	"sync"
)

type Cache struct {
	cache *sync.Map
}

func NewCache() *Cache {
	return &Cache{
		cache: &sync.Map{},
	}
}

//add to cash ============================================================
func (c *Cache) AddToCache(key string, value string) {
	c.cache.Store(key, value)
}

//get from cash ==========================================================
func (c *Cache) GetFromCache(key string) (string, bool) {
	value, ok := c.cache.Load(key)
	return fmt.Sprint(value), ok
}

//remove from cash =======================================================
func (c *Cache) RemoveFromCache(key string) {
	c.cache.Delete(key)
}

//print cash =============================================================
func (c *Cache) PrintCache() string {
	return fmt.Sprintf("cash: %v", c.cache)
}
