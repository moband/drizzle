package handlers

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/http"
)

// handleRoot handles requests to the root path
func (h *Handlers) handleRoot(conn net.Conn) {
	h.writeResponse(conn, http.StatusOK, "", nil, 0)
}

// handleEcho handles requests to the /echo/ endpoint
func (h *Handlers) handleEcho(conn net.Conn, content string) {
	bodyBytes := []byte(content)
	h.writeResponse(conn, http.StatusOK, http.ContentTypePlain, bodyBytes, len(bodyBytes))
}

// handleUserAgent handles requests to the /user-agent endpoint
func (h *Handlers) handleUserAgent(conn net.Conn, userAgent string) {
	bodyBytes := []byte(userAgent)
	h.writeResponse(conn, http.StatusOK, http.ContentTypePlain, bodyBytes, len(bodyBytes))
}

// handleFilesGet handles GET requests to the /files/{filename} endpoint
func (h *Handlers) handleFilesGet(conn net.Conn, filename string) {
	if h.config.FilesDirectory == "" {
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		return
	}

	// Prevent path traversal attacks by cleaning the path
	cleanFilename := filepath.Clean(filename)
	if strings.Contains(cleanFilename, "..") {
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		return
	}

	filePath := filepath.Join(h.config.FilesDirectory, cleanFilename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		return
	}

	h.writeResponse(conn, http.StatusOK, http.ContentTypeOctetStream, content, len(content))
}

// handleFilesPost handles POST requests to the /files/{filename} endpoint
func (h *Handlers) handleFilesPost(conn net.Conn, filename string, body []byte) {
	if h.config.FilesDirectory == "" {
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		return
	}

	// Prevent path traversal attacks by cleaning the path
	cleanFilename := filepath.Clean(filename)
	if strings.Contains(cleanFilename, "..") {
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		return
	}

	filePath := filepath.Join(h.config.FilesDirectory, cleanFilename)

	// Create the file and write the request body to it
	err := os.WriteFile(filePath, body, 0644)
	if err != nil {
		h.writeResponse(conn, http.StatusNotFound, "", nil, 0)
		return
	}

	// Return 201 Created status code
	h.writeResponse(conn, http.StatusCreated, "", nil, 0)
}
