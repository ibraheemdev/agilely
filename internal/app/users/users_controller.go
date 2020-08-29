package users

import (
	"context"
	"errors"
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

// Users controller
type Users struct {
	*engine.Engine
}

// NewController : Returns a new users controller
func NewController(e *engine.Engine) *Users {
	return &Users{Engine: e}
}

// Init :
func (u *Users) Init() (err error) {
	if err = u.Engine.Core.ViewRenderer.Load(PageLogin, PageRecoverStart, PageRecoverEnd, PageRegister); err != nil {
		return err
	}

	if err = u.Engine.Core.MailRenderer.Load(EmailConfirmHTML, EmailConfirmTxt, EmailRecoverHTML, EmailRecoverTxt); err != nil {
		return err
	}

	if _, ok := u.Core.Server.(engine.CreatingServerStorer); !ok {
		return errors.New("register module activated but storer could not be upgraded to CreatingServerStorer")
	}

	// authentication events
	u.AuthEvents.After(engine.EventRegister, u.StartConfirmationWeb)

	u.AuthEvents.Before(engine.EventAuth, u.PreventAuth)
	u.AuthEvents.Before(engine.EventAuth, u.EnsureNotLocked)
	u.AuthEvents.Before(engine.EventOAuth2, u.EnsureNotLocked)

	u.AuthEvents.After(engine.EventAuth, u.ResetLoginAttempts)
	u.AuthEvents.After(engine.EventAuth, u.CreateRememberToken)
	u.AuthEvents.After(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
		refreshExpiry(w)
		return false, nil
	})
	u.AuthEvents.After(engine.EventOAuth2, u.CreateRememberToken)

	u.AuthEvents.After(engine.EventAuthFail, u.UpdateLockAttempts)

	u.AuthEvents.After(engine.EventPasswordReset, u.ResetAllTokens)

	return nil
}

type contextKey string

const (
	// CTXKeyPID ...
	CTXKeyPID contextKey = "pid"

	// CTXKeyUser ...
	CTXKeyUser contextKey = "user"
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
)

func (c contextKey) String() string {
	return "users ctx key " + string(c)
}

// CurrentUserID retrieves the current user from the session.
// TODO(aarondl): This method never returns an error, one day we'll change
// the function signature.
func (u *Users) CurrentUserID(r *http.Request) (string, error) {
	if pid := r.Context().Value(CTXKeyPID); pid != nil {
		return pid.(string), nil
	}

	pid, _ := engine.GetSession(r, engine.SessionKey)
	return pid, nil
}

// CurrentUserIDP retrieves the current user but panics if it's not available for
// any reason.
func (u *Users) CurrentUserIDP(r *http.Request) string {
	i, err := u.CurrentUserID(r)
	if err != nil {
		panic(err)
	} else if len(i) == 0 {
		panic(engine.ErrUserNotFound)
	}

	return i
}

// CurrentUser retrieves the current user from the session and the database.
// Before the user is loaded from the database the context key is checked.
// If the session doesn't have the user ID ErrUserNotFound will be returned.
func (u *Users) CurrentUser(r *http.Request) (engine.User, error) {
	if user := r.Context().Value(CTXKeyUser); user != nil {
		return user.(engine.User), nil
	}

	pid, err := u.CurrentUserID(r)
	if err != nil {
		return nil, err
	} else if len(pid) == 0 {
		return nil, engine.ErrUserNotFound
	}

	return u.currentUser(r.Context(), pid)
}

// CurrentUserP retrieves the current user but panics if it's not available for
// any reason.
func (u *Users) CurrentUserP(r *http.Request) engine.User {
	i, err := u.CurrentUser(r)
	if err != nil {
		panic(err)
	} else if i == nil {
		panic(engine.ErrUserNotFound)
	}
	return i
}

func (u *Users) currentUser(ctx context.Context, pid string) (engine.User, error) {
	user, err := u.Core.Server.Load(ctx, pid)
	return user.(engine.User), err
}

// LoadCurrentUserID takes a pointer to a pointer to the request in order to
// change the current method's request pointer itself to the new request that
// contains the new context that has the pid in it.
func (u *Users) LoadCurrentUserID(r **http.Request) (string, error) {
	pid, err := u.CurrentUserID(*r)
	if err != nil {
		return "", err
	}

	if len(pid) == 0 {
		return "", nil
	}

	ctx := context.WithValue((**r).Context(), CTXKeyPID, pid)
	*r = (**r).WithContext(ctx)

	return pid, nil
}

// LoadCurrentUserIDP loads the current user id and panics if it's not found
func (u *Users) LoadCurrentUserIDP(r **http.Request) string {
	pid, err := u.LoadCurrentUserID(r)
	if err != nil {
		panic(err)
	} else if len(pid) == 0 {
		panic(engine.ErrUserNotFound)
	}

	return pid
}

// LoadCurrentUser takes a pointer to a pointer to the request in order to
// change the current method's request pointer itself to the new request that
// contains the new context that has the user in it. Calls LoadCurrentUserID
// so the primary id is also put in the context.
func (u *Users) LoadCurrentUser(r **http.Request) (engine.User, error) {
	if user := (*r).Context().Value(CTXKeyUser); user != nil {
		return user.(engine.User), nil
	}

	pid, err := u.LoadCurrentUserID(r)
	if err != nil {
		return nil, err
	} else if len(pid) == 0 {
		return nil, engine.ErrUserNotFound
	}

	ctx := (**r).Context()
	user, err := u.currentUser(ctx, pid)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, CTXKeyUser, user)
	*r = (**r).WithContext(ctx)
	return user, nil
}

// LoadCurrentUserP does the same as LoadCurrentUser but panics if
// the current user is not found.
func (u *Users) LoadCurrentUserP(r **http.Request) engine.User {
	user, err := u.LoadCurrentUser(r)
	if err != nil {
		panic(err)
	} else if user == nil {
		panic(engine.ErrUserNotFound)
	}

	return user
}

// IsFullyAuthed returns false if the user has a SessionHalfAuth
// in his session.
func IsFullyAuthed(r *http.Request) bool {
	_, hasHalfAuth := engine.GetSession(r, SessionHalfAuthKey)
	return !hasHalfAuth
}
