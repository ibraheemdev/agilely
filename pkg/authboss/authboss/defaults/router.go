package defaults

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router implementation
type Router struct {
	Router *httprouter.Router
}

// NewRouter creates a new router
func NewRouter(r *httprouter.Router) *Router {
	return &Router{Router: r}
}

// Get method route
func (r *Router) Get(path string, handler http.Handler) {
	r.Router.GET(path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler.ServeHTTP(w, r)
	})
}

// Post method route
func (r *Router) Post(path string, handler http.Handler) {
	r.Router.POST(path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler.ServeHTTP(w, r)
	})
}

// Delete method route
func (r *Router) Delete(path string, handler http.Handler) {
	r.Router.DELETE(path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler.ServeHTTP(w, r)
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		r.Router.ServeHTTP(w, req)
		return
	case "POST":
		r.Router.ServeHTTP(w, req)
		return
	case "DELETE":
		r.Router.ServeHTTP(w, req)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "method not allowed")
		return
	}
}
