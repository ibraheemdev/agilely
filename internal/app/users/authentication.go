package users

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/crypto/bcrypt"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

const (
	// PageLogin is for identifying the login page for parsing & validation
	PageLogin = "login.html.tpl"
)

// LoginGet simply displays the login form
func (u *Users) LoginGet(w http.ResponseWriter, r *http.Request) error {
	data := engine.HTMLData{}
	if redir := r.URL.Query().Get(engine.FormValueRedirect); len(redir) != 0 {
		data[engine.FormValueRedirect] = redir
	}
	return u.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
}

// LoginPost attempts to validate the credentials passed in
// to log in a user.
func (u *Users) LoginPost(w http.ResponseWriter, r *http.Request) error {
	logger := u.RequestLogger(r)

	validatable, err := u.Engine.Core.BodyReader.Read(PageLogin, r)
	if err != nil {
		return err
	}

	// Skip validation since all the validation happens during the database lookup and
	// password check.
	creds := engine.MustHaveUserValues(validatable)

	email := creds.GetPID()
	user, err := GetUser(r.Context(), u.Core.Database, email)

	if err == engine.ErrNoDocuments {
		logger.Infof("failed to load user requested by email: %s", email)
		data := engine.HTMLData{engine.DataErr: "Invalid Credentials"}
		return u.Engine.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
	} else if err != nil {
		return err
	}

	r = r.WithContext(context.WithValue(r.Context(), CTXKeyUser, user))

	var handled bool
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.GetPassword()))
	if err != nil {
		handled, err = u.Engine.AuthEvents.FireAfter(engine.EventAuthFail, w, r)
		if err != nil {
			return err
		} else if handled {
			return nil
		}

		logger.Infof("user %s failed to log in", email)
		data := engine.HTMLData{engine.DataErr: "Invalid Credentials"}
		return u.Engine.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
	}

	r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyValues, validatable))

	handled, err = u.AuthEvents.FireBefore(engine.EventAuth, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	handled, err = u.AuthEvents.FireBefore(engine.EventAuthHijack, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	logger.Infof("user %s logged in", email)
	engine.PutSession(w, engine.SessionKey, email)
	engine.DelSession(w, SessionHalfAuthKey)

	handled, err = u.Engine.AuthEvents.FireAfter(engine.EventAuth, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	ro := engine.RedirectOptions{
		Code:             http.StatusTemporaryRedirect,
		RedirectPath:     "/",
		FollowRedirParam: true,
	}
	return u.Engine.Core.Redirector.Redirect(w, r, ro)
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

// AuthenticatedMiddleware prevents someone from accessing a route that should be
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
func (u *Users) AuthenticatedMiddleware(requirements MWRequirements, failureResponse MWRespondOnFailure) func(http.Handler) http.Handler {
	return u.AuthenticatedMountedMiddleware(false, requirements, failureResponse)
}

// AuthenticatedMountedMiddleware hides an option from typical users in "mountPathed".
// Normal routes should never need this only engine routes (since they
// are behind mountPath typically). This method is exported only for use
// by Engine modules, normal users should use Middleware instead.
//
// If mountPathed is true, then before redirecting to a URL it will add
// the mountpath to the front of it.
func (u *Users) AuthenticatedMountedMiddleware(mountPathed bool, reqs MWRequirements, failResponse MWRespondOnFailure) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := u.RequestLogger(r)

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
					if mountPathed && len(u.Config.Mount) != 0 {
						redirURL = path.Join(u.Config.Mount, redirURL)
					}
					vals.Set(engine.FormValueRedirect, redirURL)

					ro := engine.RedirectOptions{
						Code:         http.StatusTemporaryRedirect,
						Failure:      "please re-login",
						RedirectPath: path.Join(u.Config.Mount, fmt.Sprintf("/login?%s", vals.Encode())),
					}

					if err := u.Core.Redirector.Redirect(w, r, ro); err != nil {
						log.Errorf("failed to redirect user during engine.Middleware redirect: %+v", err)
					}
					return
				}
			}

			if hasBit(reqs, RequireFullAuth) && !IsFullyAuthed(r) {
				fail(w, r)
				return
			}

			if _, err := u.LoadCurrentUser(&r); err == engine.ErrNoDocuments {
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
