package main

import (
	"fmt"
	"redis-clone/cache"
	"time"
)

func main() {
	c := cache.NewCache("cache.json")

	c.Set("name", "GoLang")
	c.Set("version", "1.23", 10*time.Second)

	if value, ok := c.Get("name"); ok {
		fmt.Println(value)
	} else {
		fmt.Println("key not found!")
	}

	time.Sleep(12 * time.Second)

	if value, exists := c.Get("version"); exists {
		fmt.Println("Value for 'version':", value)
	} else {
		fmt.Println("Key 'version' not found or expired")
	}
}
