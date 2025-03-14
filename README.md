# HTTP Server in Go

[![progress-banner](https://backend.codecrafters.io/progress/http-server/d5861641-0171-4e47-ae5f-39d0bd452434)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

HTTP/1.1 server implementation built in Go for the ["Build Your Own HTTP server" Challenge](https://app.codecrafters.io/courses/http-server/overview).

## Features

This HTTP server implements:

- Concurrent connection handling with goroutines
- Support for GET and POST methods
- Dynamic content endpoints (/echo/{string})
- User-agent information endpoint (/user-agent)
- File serving (/files/{filename}) with support for GET and POST
- HTTP compression with gzip encoding
- Robust error handling
- Security measures against path traversal

## Project Architecture

The server follows a clean architecture pattern with separation of concerns:

```
app/
├── cmd/                # Command-line entry points
│   └── server/         # Main server executable
├── internal/           # Private application code
│   ├── config/         # Configuration handling
│   ├── http/           # HTTP protocol implementation
│   ├── handlers/       # Request handlers
│   └── server/         # Core server implementation
```

## Best Practices

This project follows several Go best practices:

1. **Code Organization**
   - Separation of concerns with distinct packages
   - Internal packages for implementation details
   - Clear interfaces between components

2. **Error Handling**
   - Proper error propagation using `fmt.Errorf` with `%w` for wrapping
   - Contextual error messages
   - Graceful error recovery without crashing

3. **HTTP Implementation**
   - Proper header handling (case-insensitive)
   - Correct CRLF line endings
   - Content negotiation for compression

4. **Concurrency**
   - Connection timeouts to prevent resource exhaustion
   - Goroutines for concurrent request handling
   - Clean connection handling with deferred closures

5. **Security**
   - Path traversal prevention
   - Input validation
   - Content length validation

6. **Performance**
   - Efficient I/O with bufio
   - Response streaming
   - Compression support

## Running the Server

1. Ensure you have Go 1.20+ installed
2. Run the server:

```sh
# Simple run
./run.sh

# With a directory for file serving
./run.sh --directory /path/to/files
```

## Testing

You can test the different endpoints using curl:

```sh
# Test echo endpoint
curl -v http://localhost:4221/echo/hello

# Test user-agent endpoint
curl -v http://localhost:4221/user-agent

# Test compression
curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/hello

# Test file operations
# Read a file
curl -v http://localhost:4221/files/example.txt

# Create/update a file
curl -v --data "Hello, World!" -H "Content-Type: application/octet-stream" http://localhost:4221/files/example.txt
```

## Extending the Server

To extend the server with new endpoints:

1. Add a new handler function in `app/internal/handlers/routes.go`
2. Update the routing logic in `HandleRequest` or `handleGet`/`handlePost` methods
3. For new HTTP methods, add constants in `app/internal/http/constants.go`

## CodeCrafters Challenge

This project was built as part of the CodeCrafters "Build Your Own HTTP server" challenge.
