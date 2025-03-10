package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Parse command line arguments to get the directory
	var directory string
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--directory" && i+1 < len(args) {
			directory = args[i+1]
			break
		}
	}

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Continually accept new connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue // Continue to accept other connections
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn, directory)
	}
}

// handleConnection processes a single client connection
func handleConnection(conn net.Conn, directory string) {
	defer conn.Close() // Ensure connection is closed when function returns

	// Read the HTTP request
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
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
			return
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
	} else if strings.HasPrefix(path, "/files/") {
		// Handle file requests only if directory was specified
		if directory != "" {
			// Extract the filename from the path
			filename := path[len("/files/"):]
			filepath := filepath.Join(directory, filename)

			// Check if the file exists and read its contents
			fileContent, err := ioutil.ReadFile(filepath)
			if err == nil {
				// File exists, return it
				contentLength := len(fileContent)
				response = fmt.Sprintf(
					"HTTP/1.1 200 OK\r\n"+
						"Content-Type: application/octet-stream\r\n"+
						"Content-Length: %d\r\n\r\n",
					contentLength,
				)
				// Send the response headers
				_, err = conn.Write([]byte(response))
				if err != nil {
					fmt.Println("Error writing to connection: ", err.Error())
					return
				}

				// Send the file content separately
				_, err = conn.Write(fileContent)
				if err != nil {
					fmt.Println("Error writing file content: ", err.Error())
				}
				return
			} else {
				// File doesn't exist or can't be read
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			}
		} else {
			// No directory specified
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	// Send HTTP response
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}
