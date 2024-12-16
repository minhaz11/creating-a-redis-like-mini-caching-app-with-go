package cache

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Cache struct {
	store      map[string]string
	expiration map[string]time.Time
	mu         sync.RWMutex
	filePath   string
	logger     *zap.SugaredLogger
}

func NewCache(filePath string) (*Cache, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	sugar := logger.Sugar()

	cache := &Cache{
		store:      make(map[string]string),
		expiration: make(map[string]time.Time),
		filePath:   filePath,
		logger:     sugar,
	}

	if err := cache.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load cache from file: %w", err)
	}

	go cache.clearExpiredKeys()
	go cache.periodicSave()

	return cache, nil
}

func (c *Cache) Set(key, value string, ttl ...time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value

	if len(ttl) > 0 && ttl[0] > 0 {
		c.expiration[key] = time.Now().Add(ttl[0])
		c.logger.Infof("Key '%s' set with TTL: %s", key, ttl[0])
	} else {
		delete(c.expiration, key)
		c.logger.Infof("Key '%s' set without expiration", key)
	}

	return nil
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
		c.logger.Infof("Key '%s' has expired and was removed", key)
		return "", false
	}

	return value, true
}

func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	delete(c.expiration, key)

	c.logger.Infof("Key '%s' has been deleted", key)
	return nil
}
