/*
Package engine is a modular authentication system for the web. It tries to
remove as much boilerplate and "hard things" as possible so that each time you
start a new web project in Go, you can plug it in, configure and be off to the
races without having to think about how to store passwords or remember tokens.
*/
package engine

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/crypto/bcrypt"
)

// Engine contains a configuration and other details for running.
type Engine struct {
	Config
	Events *Events
}

// New makes a new instance of engine with a default
// configuration.
func New() *Engine {
	e := &Engine{}

	e.Events = NewEvents()

	e.Config.Defaults()
	return e
}

// UpdatePassword updates the password field of a user using the same semantics
// that register/auth do to create and verify passwords. It saves this using
// the storer.
//
// In addition to that, it also invalidates any remember me tokens, if the
// storer supports that kind of operation.
//
// If it's also desirable to log the user out, use:
// engine.DelKnown(Session|Cookie)
func (a *Engine) UpdatePassword(ctx context.Context, user AuthableUser, newPassword string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(newPassword), a.Config.Modules.BCryptCost)
	if err != nil {
		return err
	}

	user.PutPassword(string(pass))

	storer := a.Config.Storage.Server
	if err := storer.Save(ctx, user); err != nil {
		return err
	}

	rmStorer, ok := storer.(RememberingServerStorer)
	if !ok {
		return nil
	}

	return rmStorer.DelRememberTokens(ctx, user.GetPID())
}

// VerifyPassword uses engine mechanisms to check that a password is correct.
// Returns nil on success otherwise there will be an error. Simply a helper
// to do the bcrypt comparison.
func VerifyPassword(user AuthableUser, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password))
}

// MWRequirements are user requirements for engine.Middleware
// in order to access the routes in protects. Requirements is a bit-set integer
// to be able to easily combine requirements like so:
//
//   engine.RequireFullAuth
type MWRequirements int

// MWRespondOnFailure tells engine.Middleware how to respond to
// a failure to meet the requirements.
type MWRespondOnFailure int

// Middleware requirements
const (
	RequireNone MWRequirements = 0x00
	// RequireFullAuth means half-authed users will also be rejected
	RequireFullAuth MWRequirements = 0x01
)

// Middleware response types
const (
	// RespondNotFound does not allow users who are not logged in to know a
	// route exists by responding with a 404.
	RespondNotFound MWRespondOnFailure = iota
	// RespondRedirect redirects users to the login page
	RespondRedirect
	// RespondUnauthorized provides a 401, this allows users to know the page
	// exists unlike the 404 option.
	RespondUnauthorized
)

// Middleware prevents someone from accessing a route that should be
// only allowed for users who are logged in.
// It allows the user through if they are logged in (SessionKey is present in
// the session).
//
// requirements are set by logical or'ing together requirements. eg:
//
//   engine.RequireFullAuth
//
// failureResponse is how the middleware rejects the users that don't meet
// the criteria. This should be chosen from the MWRespondOnFailure constants.
func Middleware(e *Engine, requirements MWRequirements, failureResponse MWRespondOnFailure) func(http.Handler) http.Handler {
	return MountedMiddleware(e, false, requirements, failureResponse)
}

// MountedMiddleware hides an option from typical users in "mountPathed".
// Normal routes should never need this only engine routes (since they
// are behind mountPath typically). This method is exported only for use
// by Engine modules, normal users should use Middleware instead.
//
// If mountPathed is true, then before redirecting to a URL it will add
// the mountpath to the front of it.
func MountedMiddleware(e *Engine, mountPathed bool, reqs MWRequirements, failResponse MWRespondOnFailure) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := e.RequestLogger(r)

			fail := func(w http.ResponseWriter, r *http.Request) {
				switch failResponse {
				case RespondNotFound:
					log.Infof("not found for unauthorized user at: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				case RespondUnauthorized:
					log.Infof("unauthorized for unauthorized user at: %s", r.URL.Path)
					w.WriteHeader(http.StatusUnauthorized)
				case RespondRedirect:
					log.Infof("redirecting unauthorized user to login from: %s", r.URL.Path)
					vals := make(url.Values)

					redirURL := r.URL.Path
					if mountPathed && len(e.Config.Paths.Mount) != 0 {
						redirURL = path.Join(e.Config.Paths.Mount, redirURL)
					}
					vals.Set(FormValueRedirect, redirURL)

					ro := RedirectOptions{
						Code:         http.StatusTemporaryRedirect,
						Failure:      "please re-login",
						RedirectPath: path.Join(e.Config.Paths.Mount, fmt.Sprintf("/login?%s", vals.Encode())),
					}

					if err := e.Config.Core.Redirector.Redirect(w, r, ro); err != nil {
						log.Errorf("failed to redirect user during engine.Middleware redirect: %+v", err)
					}
					return
				}
			}

			if hasBit(reqs, RequireFullAuth) && !IsFullyAuthed(r) {
				fail(w, r)
				return
			}

			if _, err := e.LoadCurrentUser(&r); err == ErrUserNotFound {
				fail(w, r)
				return
			} else if err != nil {
				log.Errorf("error fetching current user: %+v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func hasBit(reqs, req MWRequirements) bool {
	return reqs&req == req
}
