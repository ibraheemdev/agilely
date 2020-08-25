package errorhandler

import (
	"fmt"
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

// ErrorHandler wraps http handlers with errors with itself
// to provide error handling.
//
// The pieces provided to this struct must be thread-safe
// since they will be handed to many pointers to themselves.
type ErrorHandler struct {
	LogWriter engine.Logger
}

// New : ErrorHandler constructor
func New(logger engine.Logger) ErrorHandler {
	return ErrorHandler{LogWriter: logger}
}

// Wrap an http handler with an error
func (e ErrorHandler) Wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return errorHandler{
		Handler:   handler,
		LogWriter: e.LogWriter,
	}
}

type errorHandler struct {
	Handler   func(w http.ResponseWriter, r *http.Request) error
	LogWriter engine.Logger
}

// ServeHTTP handles errors
func (e errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.Handler(w, r)
	if err == nil {
		return
	}

	e.LogWriter.Error(fmt.Sprintf("request error from (%s) %s: %+v", r.RemoteAddr, r.URL.String(), err))
}