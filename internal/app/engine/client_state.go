package engine

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	// SessionKey is the primarily used key by engine.
	SessionKey = "uid"
	// SessionHalfAuthKey is used for sessions that have been authenticated by
	// the remember module. This serves as a way to force full authentication
	// by denying half-authed users acccess to sensitive areas.
	SessionHalfAuthKey = "halfauth"
	// SessionLastAction is the session key to retrieve the
	// last action of a user.
	SessionLastAction = "last_action"
	// SessionOAuth2State is the xsrf protection key for oauth.
	SessionOAuth2State = "oauth2_state"
	// SessionOAuth2Params is the additional settings for oauth
	// like redirection/remember.
	SessionOAuth2Params = "oauth2_params"

	// CookieRemember is used for cookies and form input names.
	CookieRemember = "rm"

	// FlashSuccessKey is used for storing success flash messages on the session
	FlashSuccessKey = "flash_success"
	// FlashErrorKey is used for storing success flash messages on the session
	FlashErrorKey = "flash_error"
)

// ClientStateAuthEventKind is an enum.
type ClientStateAuthEventKind int

// ClientStateAuthEvent kinds
const (
	// ClientStateAuthEventPut means you should put the key-value pair into the
	// client state.
	ClientStateAuthEventPut ClientStateAuthEventKind = iota
	// ClientStateAuthEventPut means you should delete the key-value pair from the
	// client state.
	ClientStateAuthEventDel
	// ClientStateAuthEventDelAll means you should delete EVERY key-value pair from
	// the client state - though a whitelist of keys that should not be deleted
	// may be passed through as a comma separated list of keys in
	// the ClientStateAuthEvent.Key field.
	ClientStateAuthEventDelAll
)

// ClientStateAuthEvent are the different events that can be recorded during
// a request.
type ClientStateAuthEvent struct {
	Kind  ClientStateAuthEventKind
	Key   string
	Value string
}

// ClientStateReadWriter is used to create a cookie storer from an http request.
// Keep in mind security considerations for your implementation, Secure,
// HTTP-Only, etc flags.
//
// There's two major uses for this. To create session storage, and remember me
// cookies.
type ClientStateReadWriter interface {
	// ReadState should return a map like structure allowing it to look up
	// any values in the current session, or any cookie in the request
	ReadState(*http.Request) (ClientState, error)
	// WriteState can sometimes be called with a nil ClientState in the event
	// that no ClientState was read in from LoadClientState
	WriteState(http.ResponseWriter, ClientState, []ClientStateAuthEvent) error
}

// UnderlyingResponseWriter retrieves the response
// writer underneath the current one. This allows us
// to wrap and later discover the particular one that we want.
// Keep in mind this should not be used to call the normal methods
// of a responsewriter, just additional ones particular to that type
// because it's possible to introduce subtle bugs otherwise.
type UnderlyingResponseWriter interface {
	UnderlyingResponseWriter() http.ResponseWriter
}

// ClientState represents the client's current state and can answer queries
// about it.
type ClientState interface {
	Get(key string) (string, bool)
}

// ClientStateResponseWriter is used to write out the client state at the last
// moment before the response code is written.
type ClientStateResponseWriter struct {
	http.ResponseWriter

	cookieStateRW  ClientStateReadWriter
	sessionStateRW ClientStateReadWriter

	cookieState  ClientState
	sessionState ClientState

	hasWritten             bool
	cookieStateAuthEvents  []ClientStateAuthEvent
	sessionStateAuthEvents []ClientStateAuthEvent
}

// LoadClientStateMiddleware wraps all requests with the
// ClientStateResponseWriter as well as loading the current client
// state into the context for use.
func (a *Engine) LoadClientStateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := a.NewResponse(w)
		request, err := a.LoadClientState(writer, r)
		if err != nil {
			logger := a.RequestLogger(r)
			logger.Errorf("failed to load client state %+v", err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(writer, request)
	})
}

// NewResponse wraps the ResponseWriter with a ClientStateResponseWriter
func (a *Engine) NewResponse(w http.ResponseWriter) *ClientStateResponseWriter {
	return &ClientStateResponseWriter{
		ResponseWriter: w,
		cookieStateRW:  a.Core.CookieState,
		sessionStateRW: a.Core.SessionState,
	}
}

// LoadClientState loads the state from sessions and cookies
// into the ResponseWriter for later use.
func (a *Engine) LoadClientState(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	if a.Core.SessionState != nil {
		state, err := a.Core.SessionState.ReadState(r)
		if err != nil {
			return nil, err
		} else if state != nil {
			c := MustClientStateResponseWriter(w)
			c.sessionState = state
			r = r.WithContext(context.WithValue(r.Context(), CTXKeySessionState, state))
		}
	}
	if a.Core.CookieState != nil {
		state, err := a.Core.CookieState.ReadState(r)
		if err != nil {
			return nil, err
		} else if state != nil {
			c := MustClientStateResponseWriter(w)
			c.cookieState = state
			r = r.WithContext(context.WithValue(r.Context(), CTXKeyCookieState, state))
		}
	}

	return r, nil
}

// MustClientStateResponseWriter tries to find a csrw inside the response
// writer by using the UnderlyingResponseWriter interface.
func MustClientStateResponseWriter(w http.ResponseWriter) *ClientStateResponseWriter {
	for {
		if c, ok := w.(*ClientStateResponseWriter); ok {
			return c
		}

		if u, ok := w.(UnderlyingResponseWriter); ok {
			w = u.UnderlyingResponseWriter()
			continue
		}

		panic(fmt.Sprintf("ResponseWriter must be a ClientStateResponseWriter or UnderlyingResponseWriter in (see: engine.LoadClientStateMiddleware): %T", w))
	}
}

// WriteHeader writes the header, but in order to handle errors from the
// underlying ClientStateReadWriter, it has to panic.
func (c *ClientStateResponseWriter) WriteHeader(code int) {
	if !c.hasWritten {
		if err := c.putClientState(); err != nil {
			panic(err)
		}
	}
	c.ResponseWriter.WriteHeader(code)
}

// Header retrieves the underlying headers
func (c ClientStateResponseWriter) Header() http.Header {
	return c.ResponseWriter.Header()
}

// Hijack implements the http.Hijacker interface by calling the
// underlying implementation if available.
func (c ClientStateResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := c.ResponseWriter.(http.Hijacker)
	if ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("engine: underlying ResponseWriter does not support hijacking")
}

// Write ensures that the client state is written before any writes
// to the body occur (before header flush to http client)
func (c *ClientStateResponseWriter) Write(b []byte) (int, error) {
	if !c.hasWritten {
		if err := c.putClientState(); err != nil {
			return 0, err
		}
	}
	return c.ResponseWriter.Write(b)
}

// UnderlyingResponseWriter for this instance
func (c *ClientStateResponseWriter) UnderlyingResponseWriter() http.ResponseWriter {
	return c.ResponseWriter
}

func (c *ClientStateResponseWriter) putClientState() error {
	if c.hasWritten {
		panic("should not call putClientState twice")
	}
	c.hasWritten = true

	if len(c.cookieStateAuthEvents) == 0 && len(c.sessionStateAuthEvents) == 0 {
		return nil
	}

	if c.sessionStateRW != nil && len(c.sessionStateAuthEvents) > 0 {
		err := c.sessionStateRW.WriteState(c, c.sessionState, c.sessionStateAuthEvents)
		if err != nil {
			return err
		}
	}
	if c.cookieStateRW != nil && len(c.cookieStateAuthEvents) > 0 {
		err := c.cookieStateRW.WriteState(c, c.cookieState, c.cookieStateAuthEvents)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsFullyAuthed returns false if the user has a SessionHalfAuth
// in his session.
func IsFullyAuthed(r *http.Request) bool {
	_, hasHalfAuth := GetSession(r, SessionHalfAuthKey)
	return !hasHalfAuth
}

// DelAllSession deletes all variables in the session except for those on
// the whitelist.
//
// The whitelist is typically provided directly from the engine config.
//
// This is the best way to ensure the session is cleaned up after use for
// a given user. An example is when a user is expired or logged out this method
// is called.
func DelAllSession(w http.ResponseWriter, whitelist []string) {
	delAllState(w, CTXKeySessionState, whitelist)
}

// DelKnownCookie deletes all known cookie variables, which can be used
// to delete remember me pieces.
func DelKnownCookie(w http.ResponseWriter) {
	DelCookie(w, CookieRemember)
}

// PutSession puts a value into the session
func PutSession(w http.ResponseWriter, key, val string) {
	putState(w, CTXKeySessionState, key, val)
}

// DelSession deletes a key-value from the session.
func DelSession(w http.ResponseWriter, key string) {
	delState(w, CTXKeySessionState, key)
}

// GetSession fetches a value from the session
func GetSession(r *http.Request, key string) (string, bool) {
	return getState(r, CTXKeySessionState, key)
}

// PutCookie puts a value into the session
func PutCookie(w http.ResponseWriter, key, val string) {
	putState(w, CTXKeyCookieState, key, val)
}

// DelCookie deletes a key-value from the session.
func DelCookie(w http.ResponseWriter, key string) {
	delState(w, CTXKeyCookieState, key)
}

// GetCookie fetches a value from the session
func GetCookie(r *http.Request, key string) (string, bool) {
	return getState(r, CTXKeyCookieState, key)
}

func putState(w http.ResponseWriter, CTXKey contextKey, key, val string) {
	setState(w, CTXKey, ClientStateAuthEventPut, key, val)
}

func delState(w http.ResponseWriter, CTXKey contextKey, key string) {
	setState(w, CTXKey, ClientStateAuthEventDel, key, "")
}

func delAllState(w http.ResponseWriter, CTXKey contextKey, whitelist []string) {
	setState(w, CTXKey, ClientStateAuthEventDelAll, strings.Join(whitelist, ","), "")
}

func setState(w http.ResponseWriter, ctxKey contextKey, op ClientStateAuthEventKind, key, val string) {
	csrw := MustClientStateResponseWriter(w)
	ev := ClientStateAuthEvent{
		Kind: op,
		Key:  key,
	}

	if op == ClientStateAuthEventPut {
		ev.Value = val
	}

	switch ctxKey {
	case CTXKeySessionState:
		csrw.sessionStateAuthEvents = append(csrw.sessionStateAuthEvents, ev)
	case CTXKeyCookieState:
		csrw.cookieStateAuthEvents = append(csrw.cookieStateAuthEvents, ev)
	}
}

func getState(r *http.Request, ctxKey contextKey, key string) (string, bool) {
	val := r.Context().Value(ctxKey)
	if val == nil {
		return "", false
	}

	state := val.(ClientState)
	return state.Get(key)
}

// FlashSuccess returns FlashSuccessKey from the session and removes it.
func FlashSuccess(w http.ResponseWriter, r *http.Request) string {
	str, ok := GetSession(r, FlashSuccessKey)
	if !ok {
		return ""
	}

	DelSession(w, FlashSuccessKey)
	return str
}

// FlashError returns FlashError from the session and removes it.
func FlashError(w http.ResponseWriter, r *http.Request) string {
	str, ok := GetSession(r, FlashErrorKey)
	if !ok {
		return ""
	}

	DelSession(w, FlashErrorKey)
	return str
}
