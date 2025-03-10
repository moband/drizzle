package http

// HTTP Methods
const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

// HTTP Status Codes
const (
	StatusOK       = "200 OK"
	StatusCreated  = "201 Created"
	StatusNotFound = "404 Not Found"
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
