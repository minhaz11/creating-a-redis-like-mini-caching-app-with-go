package cache

import (
	"bufio"
	"io"
	"net"
	"strings"
)

func (c *Cache) HandleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			c.logger.Errorf("Error closing connection: %v", err)
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {

				c.logger.Info("Client connection closed")
				return
			}

			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					c.logger.Infof("Connection timeout: %v", netErr)
					return
				}
			}

			c.logger.Errorf("Unexpected connection error: %v", err)
			return
		}

		command = strings.TrimSpace(strings.Trim(command, "\r\n"))

		if command == "" {
			continue
		}

		response := c.CommandParser(command)

		fullResponse := response + "\n"

		_, err = conn.Write([]byte(fullResponse))
		if err != nil {
			c.logger.Errorf("Error writing response: %v", err)
			return
		}
	}
}
