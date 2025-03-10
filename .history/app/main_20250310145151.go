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
