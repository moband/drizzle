package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/config"
	"github.com/codecrafters-io/http-server-starter-go/app/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	config   *config.Config
	handlers *handlers.Handlers
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	return &Server{
		config:   cfg,
		handlers: handlers.New(cfg),
	}
}

// Start begins listening and serving HTTP requests
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %v", s.config.Address, err)
	}

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

	// Route and handle the request
	s.handlers.HandleRequest(conn, request)
}

// parseRequest reads and parses an HTTP request
func parseRequest(conn net.Conn) (map[string]interface{}, error) {

	return http.ParseRequest(conn)
}
