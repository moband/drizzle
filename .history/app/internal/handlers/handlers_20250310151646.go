package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"codecrafters-http-server-go/app/internal/config"
)

// HTTP status codes
const (
	StatusOK       = "200 OK"
	StatusNotFound = "404 Not Found"
)

// Content types
const (
	ContentTypePlain       = "text/plain"
	ContentTypeOctetStream = "application/octet-stream"
)

// Handlers contains all the HTTP request handlers
type Handlers struct {
	config *config.Config
}

// New creates a new handlers instance
func New(cfg *config.Config) *Handlers {
	return &Handlers{
		config: cfg,
	}
}

// HandleRequest routes and handles an HTTP request
func (h *Handlers) HandleRequest(conn net.Conn, request map[string]interface{}) {
	path := request["path"].(string)

	switch {
	case path == "/":
		h.handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		h.handleEcho(conn, path[len("/echo/"):])
	case path == "/user-agent":
		headers := request["headers"].(map[string]string)
		h.handleUserAgent(conn, headers["user-agent"])
	case strings.HasPrefix(path, "/files/"):
		h.handleFiles(conn, path[len("/files/"):])
	default:
		h.writeResponse(conn, StatusNotFound, "", nil, 0)
	}
}

// ParseRequest reads and parses an HTTP request
func ParseRequest(conn net.Conn) (map[string]interface{}, error) {
	request := make(map[string]interface{})
	reader := bufio.NewReader(conn)

	// Read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %v", err)
	}

	// Parse request line
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}

	request["method"] = parts[0]
	request["path"] = parts[1]

	// Read all headers
	headers := make(map[string]string)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading headers: %v", err)
		}

		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break
		}

		colonIndex := strings.Index(headerLine, ":")
		if colonIndex > 0 {
			headerName := strings.TrimSpace(headerLine[:colonIndex])
			headerValue := strings.TrimSpace(headerLine[colonIndex+1:])
			headers[strings.ToLower(headerName)] = headerValue
		}
	}

	request["headers"] = headers

	return request, nil
}

// writeResponse writes a response to the client
func (h *Handlers) writeResponse(conn net.Conn, status, contentType string, body []byte, contentLength int) {
	var response string

	if contentType != "" && body != nil {
		response = fmt.Sprintf(
			"HTTP/1.1 %s\r\n"+
				"Content-Type: %s\r\n"+
				"Content-Length: %d\r\n\r\n",
			status,
			contentType,
			contentLength,
		)
	} else {
		response = fmt.Sprintf("HTTP/1.1 %s\r\n\r\n", status)
	}

	// Write the response headers
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing headers: %v", err)
		return
	}

	// Write the body if present
	if body != nil {
		_, err = conn.Write(body)
		if err != nil {
			log.Printf("Error writing body: %v", err)
		}
	}
}
