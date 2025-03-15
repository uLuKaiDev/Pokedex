package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Entries map[string]cacheEntry
	mutEx   sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) reapLoop(expiration time.Duration) {
	ticker := time.NewTicker(expiration)
	defer ticker.Stop()
	for range ticker.C {
		c.mutEx.Lock()
		for key, entry := range c.Entries {
			if time.Since(entry.createdAt) > expiration {
				delete(c.Entries, key)
			}
		}
		c.mutEx.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mutEx.Lock()
	defer c.mutEx.Unlock()

	fmt.Printf("*** Adding key %s to cache\n", key)
	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutEx.Lock()
	defer c.mutEx.Unlock()

	entry, found := c.Entries[key]
	if !found {
		fmt.Printf("*** Cache miss for key %s\n", key)
		return nil, false
	}
	fmt.Printf("*** Cache hit for key %s\n", key)
	return entry.val, true
}
