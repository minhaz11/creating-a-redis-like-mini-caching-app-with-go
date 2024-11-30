package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (c *Cache) clearExpiredKeys() {
	for {
		time.Sleep(1 * time.Second)

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

func (c *Cache) loadFromFile() {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.filePath)

	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error opening cache file:", err)
		}
		return
	}

	defer file.Close()

	data := struct {
		Store      map[string]string
		Expiration map[string]time.Time
	}{}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding cache file:", err)
		return
	}

	c.store = data.Store
	c.expiration = data.Expiration
}


func (c *Cache) periodicSave() {
	for {
		time.Sleep(10 * time.Second)
		c.saveToFile()
	}
}
