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
func (h *Handlers) handleEcho(conn net.Conn, request *http.Request, content string) {
	bodyBytes := []byte(content)

	// Check if client accepts gzip encoding
	if request.AcceptsEncoding(http.EncodingGzip) {
		// Compress the content with gzip
		compressedBytes, err := http.CompressGzip(bodyBytes)
		if err != nil {
			// Fall back to uncompressed response if compression fails
			h.writeResponse(conn, http.StatusOK, http.ContentTypePlain, bodyBytes, len(bodyBytes))
			return
		}

		// Send response with Content-Encoding header and compressed body
		h.writeResponseWithEncoding(
			conn,
			http.StatusOK,
			http.ContentTypePlain,
			http.EncodingGzip,
			compressedBytes,
		)
	} else {
		// Standard response without encoding
		h.writeResponse(
			conn,
			http.StatusOK,
			http.ContentTypePlain,
			bodyBytes,
			len(bodyBytes),
		)
	}
}

// handleUserAgent handles requests to the /user-agent endpoint
func (h *Handlers) handleUserAgent(conn net.Conn, request *http.Request, userAgent string) {
	bodyBytes := []byte(userAgent)

	// Check if client accepts gzip encoding
	if request.AcceptsEncoding(http.EncodingGzip) {
		// Compress the content with gzip
		compressedBytes, err := http.CompressGzip(bodyBytes)
		if err != nil {
			// Fall back to uncompressed response if compression fails
			h.writeResponse(conn, http.StatusOK, http.ContentTypePlain, bodyBytes, len(bodyBytes))
			return
		}

		// Send response with Content-Encoding header and compressed body
		h.writeResponseWithEncoding(
			conn,
			http.StatusOK,
			http.ContentTypePlain,
			http.EncodingGzip,
			compressedBytes,
		)
	} else {
		// Standard response without encoding
		h.writeResponse(conn, http.StatusOK, http.ContentTypePlain, bodyBytes, len(bodyBytes))
	}
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
