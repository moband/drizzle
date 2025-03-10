package http

// HTTP Methods
const (
	GET  = "GET"
	POST = "POST"
)

// HTTP Status Codes
const (
	StatusOK               = "200 OK"
	StatusCreated          = "201 Created"
	StatusNotFound         = "404 Not Found"
	StatusMethodNotAllowed = "405 Method Not Allowed"
)

// Content Types
const (
	ContentTypePlain       = "text/plain"
	ContentTypeOctetStream = "application/octet-stream"
)

// Header names
const (
	HeaderContentType   = "content-type"
	HeaderContentLength = "content-length"
	HeaderUserAgent     = "user-agent"
)
