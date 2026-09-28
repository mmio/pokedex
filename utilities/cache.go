package utilities

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       any
}

type Cache struct {
	mutex    sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
	ticker   time.Ticker
}

func cacheCleaner(cache *Cache) {
	for range cache.ticker.C {
		cache.reapLoop()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		mutex:    sync.Mutex{},
		entries:  map[string]cacheEntry{},
		interval: interval,
		ticker:   *time.NewTicker(interval),
	}

	go cacheCleaner(&cache)

	return &cache
}

func (c *Cache) Add(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, entry := range c.entries {
		if entry.createdAt.Add(c.interval).Before(time.Now()) {
			delete(c.entries, key)
		}
	}
}
