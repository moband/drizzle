package http

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// ParseRequest reads and parses an HTTP request
func ParseRequest(conn net.Conn) (map[string]string, error) {
	request := make(map[string]any)
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
