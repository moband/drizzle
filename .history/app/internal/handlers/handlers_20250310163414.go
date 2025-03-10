package handlers

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/config"
	"github.com/codecrafters-io/http-server-starter-go/app/internal/http"
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
func (h *Handlers) HandleRequest(conn net.Conn, request *http.Request) {

	switch request.Method {
	case http.GET:
		h.handleGet(conn, request)
	case http.POST:
		h.handlePost(conn, request)
	default:
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
	}
}

// handleGet handles GET requests
func (h *Handlers) handleGet(conn net.Conn, request *http.Request) {
	switch {
	case request.Path == "/":
		h.handleRoot(conn)
	case strings.HasPrefix(request.Path, "/echo/"):
		h.handleEcho(conn, request, request.Path[len("/echo/"):])
	case request.Path == "/user-agent":
		h.handleUserAgent(conn, request, request.Headers[http.HeaderUserAgent])
	case strings.HasPrefix(request.Path, "/files/"):
		filename := request.Path[len("/files/"):]
		h.handleFilesGet(conn, filename)
	default:
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
	}
}

// handlePost handles POST requests
func (h *Handlers) handlePost(conn net.Conn, request *http.Request) {

	switch {
	case strings.HasPrefix(request.Path, "/files/"):
		filename := request.Path[len("/files/"):]
		h.handleFilesPost(conn, filename, request.Body)
	default:
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
	}
}

// writeResponse writes a response to the client
func (h *Handlers) writeResponse(conn net.Conn, status, contentType string, body []byte, contentLength int) {
	// Use the HTTP response formatter if no body is provided
	if body == nil {
		response := http.FormatResponse(status, contentType, nil)
		if _, err := conn.Write([]byte(response)); err != nil {
			log.Printf("Error writing headers: %v", err)
		}
		return
	}

	// Calculate content length if not provided
	if contentLength <= 0 {
		contentLength = len(body)
	}

	// Format the response with headers
	var response string
	if contentType != "" {
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

	// Write the body
	if _, err = conn.Write(body); err != nil {
		log.Printf("Error writing body: %v", err)
	}
}

// writeResponseWithEncoding writes a response with Content-Encoding header to the client
func (h *Handlers) writeResponseWithEncoding(conn net.Conn, status, contentType, encoding string, body []byte) {
	// Format response with encoding
	response := http.FormatResponseWithEncoding(status, contentType, encoding, body)

	// Write response headers
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing headers: %v", err)
		return
	}

	// Write the body
	if _, err = conn.Write(body); err != nil {
		log.Printf("Error writing body: %v", err)
	}
}
