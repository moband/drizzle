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
	path := request.Path
	method := request.Method

	switch {
	case path == "/":
		h.handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		h.handleEcho(conn, path[len("/echo/"):])
	case path == "/user-agent":
		h.handleUserAgent(conn, request.Headers[http.HeaderUserAgent])
	case strings.HasPrefix(path, "/files/"):
		filename := path[len("/files/"):]
		if method == http.GET {
			h.handleFilesGet(conn, filename)
		} else if method == http.POST {
			h.handleFilesPost(conn, filename, request.Body)
		} else {
			h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		}
	default:
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
	}
}

func (h *Handlers) handleGet(conn net.Conn, request *http.Request) {
	path := request.Path
	switch path {
	case path == "/":
		h.handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		h.handleEcho(conn, request.Path[len("/echo/"):])
	case path == "/user-agent":
		h.handleUserAgent(conn, request.Headers[http.HeaderUserAgent])

	}
}

func (h *Handlers) handlePost(conn net.Conn, request *http.Request) {

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
