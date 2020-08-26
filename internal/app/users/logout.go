package users

import (
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

// Logout the user
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) error {
	logger := u.RequestLogger(r)

	user, err := u.CurrentUser(r)
	if err == nil && user != nil {
		logger.Infof("user %s logged out", user.GetPID())
	} else {
		logger.Info("user (unknown) logged out")
	}

	var handled bool
	handled, err = u.AuthEvents.FireBefore(engine.EventLogout, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	engine.DelAllSession(w, u.Config.SessionStateWhitelistKeys)
	engine.DelKnownCookie(w)

	handled, err = u.Engine.AuthEvents.FireAfter(engine.EventLogout, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: "/",
		Success:      "You have been logged out",
	}
	return u.Engine.Core.Redirector.Redirect(w, r, ro)
}
