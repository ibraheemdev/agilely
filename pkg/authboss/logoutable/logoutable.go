// Package logoutable allows users to log out (from auth or oauth2 logins)
package logoutable

import (
	"net/http"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
)

func init() {
	authboss.RegisterModule("logout", &Logout{})
}

// Logout module
type Logout struct {
	*authboss.Authboss
}

// Init the module
func (l *Logout) Init(ab *authboss.Authboss) error {
	l.Authboss = ab

	l.Authboss.Config.Core.Router.DELETE("/logout", l.Authboss.Core.ErrorHandler.Wrap(l.Logout))

	return nil
}

// Logout the user
func (l *Logout) Logout(w http.ResponseWriter, r *http.Request) error {
	logger := l.RequestLogger(r)

	user, err := l.CurrentUser(r)
	if err == nil && user != nil {
		logger.Infof("user %s logged out", user.GetPID())
	} else {
		logger.Info("user (unknown) logged out")
	}

	var handled bool
	handled, err = l.Events.FireBefore(authboss.EventLogout, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	authboss.DelAllSession(w, l.Config.Storage.SessionStateWhitelistKeys)
	authboss.DelKnownCookie(w)

	handled, err = l.Authboss.Events.FireAfter(authboss.EventLogout, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	ro := authboss.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: "/",
		Success:      "You have been logged out",
	}
	return l.Authboss.Core.Redirector.Redirect(w, r, ro)
}
