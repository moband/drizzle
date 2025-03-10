package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// ParseRequest reads and parses an HTTP request
func ParseRequest(conn net.Conn) (map[string]any, error) {
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
	contentLength := 0
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
			headerName = strings.ToLower(headerName)
			headers[headerName] = headerValue

			// Check for Content-Length header
			if headerName == "content-length" {
				contentLength, err = strconv.Atoi(headerValue)
				if err != nil {
					return nil, fmt.Errorf("invalid Content-Length: %v", err)
				}
			}
		}
	}

	request["headers"] = headers

	// Read request body if Content-Length is present
	if contentLength > 0 {
		body := make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %v", err)
		}
		request["body"] = body
	} else {
		request["body"] = []byte{}
	}

	return request, nil
}
