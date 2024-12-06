# A simple Mini Cache Server

This project is a simple in-memory key-value store with basic functionalities. It supports commands like `SET`, `GET`, `DEL`, and `EXPIRE`.

## Project Structure


## Getting Started

### Prerequisites

- Go 1.23.3 or later

### Installation

1. Clone the repository:
    ```sh
    [git clone https://github.com/minhaz11/redis-clone.git](https://github.com/minhaz11/A-mini-caching-app-with-Go.git)
    cd to cloned folder
    ```

2. Build the project:
    ```sh
    go build -o 'folder name' main.go
    ```

3. Run the server:
    ```sh
    ./folder name
    ```

## Usage

The server listens on port `6369`. You can connect to it using `telnet` or any TCP client.

### Commands

- **SET**: Set a key-value pair.
    ```
    SET key value [ttl]
    ```

- **GET**: Get the value of a key.
    ```
    GET key
    ```

- **DEL**: Delete a key.
    ```
    DEL key
    ```

- **EXPIRE**: Set a time-to-live (TTL) for a key.
    ```
    EXPIRE key ttl
    ```

### Example

```sh
$ telnet localhost 6369
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
SET foo bar
OK Key 'foo' set successfully.
GET foo
bar
EXPIRE foo 10s
OK Key 'foo' expiration set to '10s'.
GET foo
bar
