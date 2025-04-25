package services

import (
	"sync"
	"time"
	"ws-middle-cache/internal/middleware"
)

// CacheItem represents a single cache item with a value and expiration time
type CacheItem struct {
	Value      interface{}
	Expiration int64 // Unix timestamp in nanoseconds
}

// Cache is a simple in-memory cache structure with expiration control
type Cache struct {
	mu    sync.RWMutex
	store map[string]CacheItem
}

// NewCache creates a new instance of Cache
func NewCache() *Cache {
	return &Cache{
		store: make(map[string]CacheItem),
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.store[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// Remove expired item
		c.mu.RUnlock()
		c.Delete(key)
		c.mu.RLock()
		return nil, false
	}

	return item.Value, true
}

// Set stores a value in the cache with an optional expiration time
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.store[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Clear removes all values from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]CacheItem)
}

// CleanExpired removes all expired items from the cache
func (c *Cache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.store {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.store, key)
		}
	}
}

// Keys returns a slice of all keys currently stored in the cache
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	logger := middleware.NewLogger()

	keys := make([]string, 0, len(c.store))
	for key := range c.store {
		logger.Debug("Key:", key)
		keys = append(keys, key)
	}

	return keys
}

func (c *Cache) StartCleanupTask(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.CleanExpired()
		}
	}()
}
