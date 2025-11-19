# HTTP Server in Go

A simple HTTP server implementation in Go, along with a TCP listener for debugging HTTP requests.

## Prerequisites

- Go 1.25.3 or later

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/guruorgoru/http-server-in-go.git
   cd http-server-in-go
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the HTTP Server

The HTTP server serves basic routes and requires a `PORT` environment variable.

1. Ensure `.env` file exists with `PORT=8414` (or your preferred port).

2. Run the server:
   ```bash
   go run ./cmd/httpserver
   ```

   Or build and run:
   ```bash
   go build -o bin/httpserver ./cmd/httpserver
   ./bin/httpserver
   ```

The server will start on the specified port and handle the following routes:
- `GET /` - Returns "GuruOrGoru"
- `GET /skillissues` - Returns 400 Bad Request with "You got some skill issues"
- `GET /myissues` - Returns 500 Internal Server Error with "Sorry, I got some skill issues"
- Any other path - Returns 404 Not Found

## Running the TCP Listener

The TCP listener accepts raw TCP connections on port 42069 and parses HTTP requests for debugging.

Run the listener:
```bash
go run ./cmd/tcplistener
```

Or build and run:
```bash
go build -o bin/tcplistener ./cmd/tcplistener
./bin/tcplistener
```

## Usage Examples

### Testing the HTTP Server

Use curl to test the server:

```bash
# Root path
curl http://localhost:8414/

# Skill issues
curl http://localhost:8414/skillissues

# My issues
curl http://localhost:8414/myissues

# Not found
curl http://localhost:8414/unknown
```

### Testing the TCP Listener

Send raw HTTP requests to the TCP listener:

```bash
# Using netcat or similar
echo -e "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n" | nc localhost 42069
```

The listener will print the parsed request details to stdout.

## Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```
