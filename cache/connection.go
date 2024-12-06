package cache

import (
	"bufio"
	"fmt"
	"net"
)

func (c *Cache) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
	
		if err != nil {
			return
		}

		response := c.CommandParser(command)

		_, err = conn.Write([]byte(response))

		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	}
}