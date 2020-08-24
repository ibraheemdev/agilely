package users

import (
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

func init() {
	engine.RegisterModule("logout", &Logout{})
}

// Logout module
type Logout struct {
	*engine.Engine
}

// Init the module
func (l *Logout) Init(ab *engine.Engine) error {
	l.Engine = ab

	l.Engine.Config.Core.Router.DELETE("/logout", l.Engine.Core.ErrorHandler.Wrap(l.Logout))

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
	handled, err = l.Events.FireBefore(engine.EventLogout, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	engine.DelAllSession(w, l.Config.Storage.SessionStateWhitelistKeys)
	engine.DelKnownCookie(w)

	handled, err = l.Engine.Events.FireAfter(engine.EventLogout, w, r)
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
	return l.Engine.Core.Redirector.Redirect(w, r, ro)
}
