package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	store      map[string]string
	expiration map[string]time.Time
	mu         sync.Mutex
	filePath   string
}

func NewCache(filePath string) *Cache {
	cache := &Cache{
		store:      make(map[string]string),
		expiration: make(map[string]time.Time),
		filePath: filePath,
	}

	cache.loadFromFile()
	go cache.clearExpiredKeys()
	go cache.periodicSave()

	return cache
}




func (c *Cache) Set(key, value string, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value

	if len(ttl) > 0 && ttl[0] > 0 {

		c.expiration[key] = time.Now().Add(ttl[0])
		fmt.Printf("Key '%s' set to '%s' with TTL: %s\n", key, value, ttl[0])

	} else {

		delete(c.expiration, key)
		fmt.Printf("Key '%s' set to '%s' without expiration\n", key, value)
	}
}


func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.store[key]

	if !ok {
		return "", false
	}

	if expTime, expExists := c.expiration[key]; expExists && time.Now().After(expTime) {
		delete(c.store, key)
		delete(c.expiration, key)
		fmt.Printf("Key '%s' has expired and was removed\n", key)
		return "", false
	}

	return value, true
}


func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	delete(c.expiration, key)

	fmt.Printf("key '%s' has been deleted \n", key)
}
