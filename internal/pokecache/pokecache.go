package pokecache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	Entry map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{Entry: make(map[string]cacheEntry)}

	// call reap loop with interval
	go cache.reapLoop(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entry[key] = cacheEntry{createdAt: time.Now(), val: val}

	return nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.Entry[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return entry.val, nil
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.Entry {
			if time.Since(v.createdAt) > interval {
				delete(c.Entry, k)
			}
		}
		c.mu.Unlock()
	}
}
