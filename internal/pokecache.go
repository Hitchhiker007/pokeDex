package pokecache

import (
	"sync"
	"time"
)

// I used a Cache struct to hold a map[string]cacheEntry and a mutex to protect the map across goroutines
type Cache struct {
	items    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration // configurable interval (e.g. TTL)
}

type cacheEntry struct {
	createdAt time.Time // createdAt - A time.Time that represents when the entry was created.
	val       []byte    // val - A []byte that represents the raw data we're caching.
}

// expose a NewCache() function that creates a new cache with a configurable interval (time.Duration).
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		items:    make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop() // start background cleanup
	return c
}

// Create a cache.Add() method that adds a new entry to the cache. It should take a key (a string) and a val (a []byte).
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock() // lock the mutex before writing to the map so multiple goroutines dont corrupt it
	defer c.mu.Unlock()

	// store a cacheEntry with the current timestamp and the value
	c.items[key] = cacheEntry{
		createdAt: time.Now(), // later createdAt + interval let us expire entries if too old
		val:       val,
	}
	// fmt.Printf("[Cache] Added key: %s\n", key)
}

// reate a cache.Get() method that gets an entry from the cache. It should take a key (a string) and return a []byte and a bool.
// The bool should be true if the entry was found and false if it wasn't.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.items[key]
	if !found {
		return nil, false
	}
	// fmt.Printf("[Cache] Key exists! Key: %s\n", key)
	return entry.val, true
}

// Create a cache.reapLoop() method that is called when the cache is created (by the NewCache function).
// Each time an interval (the time.Duration passed to NewCache) passes it should remove any entries that are older than the interval.
// This makes sure that the cache doesn't grow too large over time.
// For example, if the interval is 5 seconds, and an entry was added 7 seconds ago, that entry should be removed.
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.items {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.items, key)
				// fmt.Printf("[Cache] Reaped expired key: %s\n", key)
			}
		}
		c.mu.Unlock()
	}
}
