package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// HTTP Methods
const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

// HTTP Status Codes
const (
	StatusOK       = "200 OK"
	StatusCreated  = "201 Created"
	StatusNotFound = "404 Not Found"
)

// Content Types
const (
	ContentTypePlain       = "text/plain"
	ContentTypeOctetStream = "application/octet-stream"
)

// Header names
const (
	HeaderContentType   = "content-type"
	HeaderContentLength = "content-length"
	HeaderUserAgent     = "user-agent"
)

// Request represents an HTTP request
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

// ParseRequest reads and parses an HTTP request from a connection
func ParseRequest(conn net.Conn) (map[string]any, error) {
	reader := bufio.NewReader(conn)

	// Parse the request line
	method, path, err := parseRequestLine(reader)
	if err != nil {
		return nil, err
	}

	// Parse headers
	headers, err := parseHeaders(reader)
	if err != nil {
		return nil, err
	}

	// Parse body based on Content-Length header
	body, err := parseBody(reader, headers)
	if err != nil {
		return nil, err
	}

	// Build the request map (keeping the same structure for compatibility)
	request := make(map[string]any)
	request["method"] = method
	request["path"] = path
	request["headers"] = headers
	request["body"] = body

	return request, nil
}

// parseRequestLine parses the HTTP request line
func parseRequestLine(reader *bufio.Reader) (method, path string, err error) {
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("error reading request line: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid request line: %s", requestLine)
	}

	return parts[0], parts[1], nil
}

// parseHeaders reads and parses HTTP headers
func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)

	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading headers: %w", err)
		}

		headerLine = strings.TrimSpace(headerLine)
		// Empty line signals end of headers
		if headerLine == "" {
			break
		}

		colonIndex := strings.Index(headerLine, ":")
		if colonIndex > 0 {
			name := strings.ToLower(strings.TrimSpace(headerLine[:colonIndex]))
			value := strings.TrimSpace(headerLine[colonIndex+1:])
			headers[name] = value
		}
	}

	return headers, nil
}

// parseBody reads the request body based on Content-Length header
func parseBody(reader *bufio.Reader, headers map[string]string) ([]byte, error) {
	contentLengthStr, exists := headers[HeaderContentLength]
	if !exists {
		return []byte{}, nil
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length: %w", err)
	}

	// Read the body only if Content-Length > 0
	if contentLength > 0 {
		body := make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		return body, nil
	}

	return []byte{}, nil
}

// FormatResponse formats an HTTP response
func FormatResponse(status, contentType string, body []byte) string {
	if contentType != "" && body != nil {
		return fmt.Sprintf(
			"HTTP/1.1 %s\r\n"+
				"Content-Type: %s\r\n"+
				"Content-Length: %d\r\n\r\n",
			status,
			contentType,
			len(body),
		)
	}

	return fmt.Sprintf("HTTP/1.1 %s\r\n\r\n", status)
}
