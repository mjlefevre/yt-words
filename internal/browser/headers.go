package browser

import "net/http"

// DefaultHeaders returns a map of default headers used for HTTP requests
func DefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "keep-alive",
		"Cache-Control":   "max-age=0",
	}
}

// SetDefaultHeaders sets default headers on the given request
// If customHeaders is provided, they will override the default values
func SetDefaultHeaders(req *http.Request, customHeaders ...map[string]string) {
	// Set default headers
	for key, value := range DefaultHeaders() {
		req.Header.Set(key, value)
	}

	// If custom headers are provided, override defaults
	if len(customHeaders) > 0 {
		for key, value := range customHeaders[0] {
			req.Header.Set(key, value)
		}
	}
}
