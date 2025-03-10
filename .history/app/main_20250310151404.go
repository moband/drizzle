package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// Server represents the HTTP server configuration
type Server struct {
	directory string
	address   string
}

// NewServer creates a new HTTP server instance
func NewServer(directory, address string) *Server {
	return &Server{
		directory: directory,
		address:   address,
	}
}

// Start begins listening and serving HTTP requests
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %v", s.address, err)
	}

	log.Printf("Server started on %s", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection processes a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Set read deadline to prevent hanging connections
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Parse the request
	request, err := parseRequest(conn)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return
	}

	// Handle the request based on the path
	s.routeRequest(conn, request)
}

// parseRequest reads and parses an HTTP request
func parseRequest(conn net.Conn) (map[string]string, error) {
	request := make(map[string]string)
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

// routeRequest routes the request to the appropriate handler
func (s *Server) routeRequest(conn net.Conn, request map[string]string) {
	path := request["path"]

	switch {
	case path == "/":
		s.handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		s.handleEcho(conn, path[len("/echo/"):])
	case path == "/user-agent":
		headers := request["headers"].(map[string]string)
		s.handleUserAgent(conn, headers["user-agent"])
	case strings.HasPrefix(path, "/files/"):
		s.handleFiles(conn, path[len("/files/"):])
	default:
		s.writeResponse(conn, StatusNotFound, "", nil, 0)
	}
}

// handleRoot handles requests to the root path
func (s *Server) handleRoot(conn net.Conn) {
	s.writeResponse(conn, StatusOK, "", nil, 0)
}

// handleEcho handles requests to the /echo/ endpoint
func (s *Server) handleEcho(conn net.Conn, content string) {
	s.writeResponse(conn, StatusOK, ContentTypePlain, []byte(content), len(content))
}

// handleUserAgent handles requests to the /user-agent endpoint
func (s *Server) handleUserAgent(conn net.Conn, userAgent string) {
	s.writeResponse(conn, StatusOK, ContentTypePlain, []byte(userAgent), len(userAgent))
}

// handleFiles handles requests to the /files/ endpoint
func (s *Server) handleFiles(conn net.Conn, filename string) {
	if s.directory == "" {
		s.writeResponse(conn, StatusNotFound, "", nil, 0)
		return
	}

	// Prevent path traversal attacks by cleaning the path
	cleanFilename := filepath.Clean(filename)
	if strings.Contains(cleanFilename, "..") {
		s.writeResponse(conn, StatusNotFound, "", nil, 0)
		return
	}

	filePath := filepath.Join(s.directory, cleanFilename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		s.writeResponse(conn, StatusNotFound, "", nil, 0)
		return
	}

	s.writeResponse(conn, StatusOK, ContentTypeOctetStream, content, len(content))
}

// writeResponse writes a response to the client
func (s *Server) writeResponse(conn net.Conn, status, contentType string, body []byte, contentLength int) {
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

func main() {
	// Parse command line flags
	directory := flag.String("directory", "", "Directory to serve files from")
	flag.Parse()

	// Create and start the server
	server := NewServer(*directory, "0.0.0.0:4221")

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
