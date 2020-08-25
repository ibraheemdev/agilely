package users

import (
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

// InitLogout the module
func (u *Users) InitLogout() error {

	u.Engine.Config.Core.Router.DELETE("/logout", u.Engine.Core.ErrorHandler.Wrap(u.Logout))

	return nil
}

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
	handled, err = u.Events.FireBefore(engine.EventLogout, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	engine.DelAllSession(w, u.Config.Storage.SessionStateWhitelistKeys)
	engine.DelKnownCookie(w)

	handled, err = u.Engine.Events.FireAfter(engine.EventLogout, w, r)
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
