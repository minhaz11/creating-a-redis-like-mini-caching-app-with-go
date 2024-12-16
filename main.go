package main

import (
	"fmt"
	"os"

	"github.com/minhaz11/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		fmt.Println("Failed to create server:", err)
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		server.Logger.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
