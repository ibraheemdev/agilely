package users

import (
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
