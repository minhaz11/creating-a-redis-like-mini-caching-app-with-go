package cache

import (
	"bufio"
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

		conn.Write([]byte(response))
	}
}