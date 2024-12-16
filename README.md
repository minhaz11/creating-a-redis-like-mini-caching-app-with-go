# A Simple mini cache server with GoLang

This project is a simple in-memory key-value store with basic functionalities. It supports commands like `SET`, `GET`, `DEL`, and `EXPIRE`.

NB: It supports only the string (key, value) currently.


## Getting Started

### Prerequisites

- Go 1.23.3 or later

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/minhaz11/A-mini-caching-app-with-Go.git

    cd to cloned folder
    ```

2. Build the project and run:
    ```sh
    go build .
    ```

3. Or Run the server:
    ```sh
     go run .
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
```

##=======================================


### Example JavaScript (Node.js) Client

```javascript

const net = require('net');

const client = new net.Socket();
client.connect(6369, '127.0.0.1', function() {
    console.log('Connected to cache server');
    client.write('SET foo bar 1005s\n');  
});

client.on('data', function(data) {
    console.log('Received:', data.toString());
   
    client.write('GET foo\n');
    
});

client.on('data', function(data) {
    console.log('Received:', data.toString());

    client.end();
});

client.on('close', function() {
    console.log('Connection closed');
});

client.on('error', function(err) {
    console.log('Error:', err);
});

```


### Example Python Client

```Python

import socket

# Create a socket object
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Connect to the server on 127.0.0.1:6369
client.connect(('127.0.0.1', 6369))

# Send the SET command to the server
print('Connected to cache server')
client.sendall(b'SET foo bar 1005s\n')

# Receive the response from the server
data = client.recv(1024)
print('Received:', data.decode())

# Send the GET command to the server
client.sendall(b'GET foo\n')

# Receive the response from the server
data = client.recv(1024)
print('Received:', data.decode())

# Close the connection
client.close()
print('Connection closed')

```





