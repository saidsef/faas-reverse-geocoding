// Package cache provides caching for http requests.
package cache

import (
	"sync"
	"time"

	"github.com/saidsef/faas-reverse-geocoding/internal/utils"
)

// CacheItem represents a single item in the cache.
type CacheItem struct {
	Response   interface{} // The cached response.
	Expiration time.Time   // The expiration time of the cached item.
}

// Cache is a simple in-memory cache with expiration.
type Cache struct {
	sync.Mutex
	data map[string]CacheItem // The map holding the cached items.
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

// Set adds an item to the cache.
// key: The key under which the item is stored.
// value: The item to be cached.
// duration: The duration for which the item should be cached.
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = CacheItem{
		Response:   value,
		Expiration: time.Now().Add(duration),
	}
}

// Get retrieves an item from the cache.
// key: The key of the item to retrieve.
// Returns the cached item and a boolean indicating whether the item was found and is not expired.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	item, found := c.data[key]
	cacheExpiration := time.Now().After(item.Expiration)
	cacheTime := time.Until(item.Expiration)
	if !found || cacheExpiration {
		if utils.Verbose {
			utils.Logger.Debugf("Cache MISS for key: %s, expired: %t", key, cacheExpiration)
		}
		delete(c.data, key)
		return nil, false
	}
	if utils.Verbose {
		utils.Logger.Debugf("Cache HIT for key: %s, expired: %t, cacheTime: %s", key, cacheExpiration, cacheTime)
	}
	return item.Response, true
}
