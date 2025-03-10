package handlers

import (
	"net"
	"os"
	"path/filepath"
	"strings"
)

// handleRoot handles requests to the root path
func (h *Handlers) handleRoot(conn net.Conn) {
	h.writeResponse(conn, StatusOK, "", nil, 0)
}

// handleEcho handles requests to the /echo/ endpoint
func (h *Handlers) handleEcho(conn net.Conn, content string) {
	bodyBytes := []byte(content)
	h.writeResponse(conn, StatusOK, ContentTypePlain, bodyBytes, len(bodyBytes))
}

// handleUserAgent handles requests to the /user-agent endpoint
func (h *Handlers) handleUserAgent(conn net.Conn, userAgent string) {
	bodyBytes := []byte(userAgent)
	h.writeResponse(conn, StatusOK, ContentTypePlain, bodyBytes, len(bodyBytes))
}

// handleFiles handles requests to the /files/ endpoint
func (h *Handlers) handleFiles(conn net.Conn, filename string) {
	if h.config.FilesDirectory == "" {
		h.writeResponse(conn, StatusNotFound, "", nil, 0)
		return
	}

	// Prevent path traversal attacks by cleaning the path
	cleanFilename := filepath.Clean(filename)
	if strings.Contains(cleanFilename, "..") {
		h.writeResponse(conn, StatusNotFound, "", nil, 0)
		return
	}

	filePath := filepath.Join(h.config.FilesDirectory, cleanFilename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		h.writeResponse(conn, StatusNotFound, "", nil, 0)
		return
	}

	h.writeResponse(conn, StatusOK, ContentTypeOctetStream, content, len(content))
}
