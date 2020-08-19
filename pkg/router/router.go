// Package router is a thin wrapper over httprouter that
// improves maintainability and integrates alice middleware chaining
package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Router implementation
type Router struct {
	*httprouter.Router
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

// GET is a shortcut for router.Handle(http.MethodGet, path, handle)
func (rt *Router) GET(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodGet, path, rt.Middlewares.Then(handler))
}

// HEAD is a shortcut for router.Handle(http.MethodHead, path, handle)
func (rt *Router) HEAD(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodHead, path, rt.Middlewares.Then(handler))
}

// OPTIONS is a shortcut for router.Handle(http.MethodOptions, path, handle)
func (rt *Router) OPTIONS(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodOptions, path, rt.Middlewares.Then(handler))
}

// POST is a shortcut for router.Handle(http.MethodPost, path, handle)
func (rt *Router) POST(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodPost, path, rt.Middlewares.Then(handler))
}

// PUT is a shortcut for router.Handle(http.MethodPut, path, handle)
func (rt *Router) PUT(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodPut, path, rt.Middlewares.Then(handler))
}

// PATCH is a shortcut for router.Handle(http.MethodPatch, path, handle)
func (rt *Router) PATCH(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodPatch, path, rt.Middlewares.Then(handler))
}

// DELETE is a shortcut for router.Handle(http.MethodDelete, path, handle)
func (rt *Router) DELETE(path string, handler http.Handler) {
	rt.Router.Handler(http.MethodDelete, path, rt.Middlewares.Then(handler))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rt.Router.ServeHTTP(w, req)
}
