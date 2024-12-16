package cache

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (c *Cache) CommandParser(command string) string {

	command = strings.TrimSpace(command)

	segments := strings.Fields(command)

	switch strings.ToUpper(segments[0]) {
	case "SET":
		key, err := parseSetCmd(c, segments)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("OK Key '%s' set successfully.\n", key)

	case "GET":
		value, err := parseGetCmd(c, segments)
		if err != nil {
			return err.Error()
		}
		return value

	case "DEL":
		key, err := parseDelCmd(c, segments)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("OK Key '%s' deleted.\n", key)

	case "EXPIRE":
		key, parsedTTL, err := parseExpireCmd(c, segments)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("OK Key '%s' expiration set to '%s'.\n", key, parsedTTL)
	
	default:
		return "Error: Unknown command."

	}
}

func parseSetCmd(c *Cache, segments []string) (string, error) {
	if len(segments) < 3 {
		return "", errors.New("Error: SET command requires at least a key and value.")
	}

	key := segments[1]
	value := segments[2]

	var ttl time.Duration

	if len(segments) > 3 {
		parseTTL, err := time.ParseDuration(segments[3])

		if err != nil {
			return "",  errors.New("Error: Invalid TTL format.")
		}

		ttl = parseTTL
	}

	c.Set(key, value, ttl)

	return key, nil
}

func parseGetCmd(c *Cache, segments []string) (string, error) {
	if len(segments) < 2 {
		return "", errors.New("Error: GET command requires a key.")
	}

	key := segments[1]

	value, ok := c.Get(key)

	if !ok {
		return "", errors.New("Error: Key not found or expired.")
	}

	return value, nil
}

func parseDelCmd(c *Cache, segments []string) (string, error) {
	
	if len(segments) < 2 {
		return "",errors.New("Error: DEL command requires a key.")
	}

	key := segments[1]
	c.Delete(key)

	return key, nil
}

func parseExpireCmd(c *Cache, segments []string) (string, time.Duration, error) {
	if len(segments) < 3 {
		return "", 0, errors.New("Error: EXPIRE command requires a key and TTL.")
	}

	key := segments[1]
	parsedTTL, err := time.ParseDuration(segments[2])

	if err != nil {
		return "", 0, errors.New("Error: Invalid TTL format.")
	}

	c.Set(key,c.store[key], parsedTTL)

	return key, parsedTTL, nil
}