package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

func (c *Cache) clearExpiredKeys() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for k, exp := range c.expiration {
			if time.Now().After(exp) {
				delete(c.store, k)
				delete(c.expiration, k)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) saveToFile() {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := struct {
		Store      map[string]string
		Expiration map[string]time.Time
	}{
		Store:      c.store,
		Expiration: c.expiration,
	}

	file, err := os.Create(c.filePath)
	if err != nil {
		fmt.Println("Error saving cache to file.", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err = encoder.Encode(data); err != nil {
		fmt.Println("Error encoding cache data to file:", err)
	}
}

func (c *Cache) loadFromFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.filePath)

	if err != nil {
		if !os.IsNotExist(err) {
			return errors.New("Error opening cache file.")
		}
		return errors.New(err.Error())
	}

	defer file.Close()

	data := struct {
		Store      map[string]string
		Expiration map[string]time.Time
	}{}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&data); err != nil {
		return errors.New(err.Error())
	}

	c.store = data.Store
	c.expiration = data.Expiration

	return nil
}


func (c *Cache) periodicSave() {
	for {
		time.Sleep(10 * time.Second)
		c.saveToFile()
	}
}
