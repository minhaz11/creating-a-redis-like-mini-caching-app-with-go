package main

import (
	"fmt"
	"net"
	"github.com/minhaz11/cache"
)

func main() {
	c := cache.NewCache("cache.json")

	ln, err := net.Listen("tcp", ":6369")

	if err != nil {
		fmt.Println("Error starting server:", err.Error())
		return
	}

	defer ln.Close()

	fmt.Println("Cache server is listening on port 6369...")

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		go c.HandleConnection(conn)

	}
}
