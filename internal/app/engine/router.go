package engine

import (
	"net/http"
)

// Router can register routes to later be used by the web application
type Router interface {
	http.Handler

	GET(path string, handler http.Handler)
	POST(path string, handler http.Handler)
	DELETE(path string, handler http.Handler)
	PUT(path string, handler http.Handler)
}
