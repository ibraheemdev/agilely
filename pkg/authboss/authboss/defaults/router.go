package defaults

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Router implementation
type Router struct {
	Router      *httprouter.Router
	Middlewares alice.Chain
}

// NewRouter creates a new router
func NewRouter(r *httprouter.Router) *Router {
	return &Router{Router: r}
}

// Use : Applies an alice middleware chain to handlers
func (rt *Router) Use(c alice.Chain) {
	rt.Middlewares = c
}

// Wrapper : Acts as a compatibility layer between http.Handler and httprouter.Handle
func Wrapper(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		next.ServeHTTP(w, r)
	}
}

// Get method route
func (rt *Router) Get(path string, handler http.Handler) {
	rt.Router.GET(path, Wrapper(rt.Middlewares.Then(handler)))
}

// Post method route
func (rt *Router) Post(path string, handler http.Handler) {
	rt.Router.POST(path, Wrapper(rt.Middlewares.Then(handler)))
}

// Delete method route
func (rt *Router) Delete(path string, handler http.Handler) {
	rt.Router.DELETE(path, Wrapper(rt.Middlewares.Then(handler)))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
	case "POST":
	case "DELETE":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "method not allowed")
		return
	}
	rt.Router.ServeHTTP(w, req)
}
