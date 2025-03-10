package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Read the HTTP request
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	// Parse the request line to extract the path
	// Request line format: METHOD PATH HTTP-VERSION
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	var path string
	if len(parts) >= 2 {
		path = parts[1]
	}

	// Read all headers
	headers := make(map[string]string)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading headers: ", err.Error())
			os.Exit(1)
		}

		// Trim the trailing CRLF
		headerLine = strings.TrimSpace(headerLine)

		// Empty line signifies end of headers
		if headerLine == "" {
			break
		}

		// Parse header (format: "Name: Value")
		colonIndex := strings.Index(headerLine, ":")
		if colonIndex > 0 {
			headerName := strings.TrimSpace(headerLine[:colonIndex])
			headerValue := strings.TrimSpace(headerLine[colonIndex+1:])
			// Store header in case-insensitive way
			headers[strings.ToLower(headerName)] = headerValue
		}
	}

	// Determine response based on path
	var response string
	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path, "/echo/") {
		// Extract the string part after "/echo/"
		echoStr := path[len("/echo/"):]

		// Construct response with headers and body
		contentLength := len(echoStr)
		response = fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/plain\r\n"+
				"Content-Length: %d\r\n\r\n"+
				"%s",
			contentLength,
			echoStr,
		)
	} else if path == "/user-agent" {
		// Get the User-Agent header
		userAgent := headers["user-agent"]

		// Construct response with headers and body
		contentLength := len(userAgent)
		response = fmt.Sprintf(
			"HTTP/1.1 200 OK\r\n"+
				"Content-Type: text/plain\r\n"+
				"Content-Length: %d\r\n\r\n"+
				"%s",
			contentLength,
			userAgent,
		)
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	// Send HTTP response
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

	// Close the connection
	conn.Close()
}
