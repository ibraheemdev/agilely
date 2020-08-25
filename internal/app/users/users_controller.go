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
	if err = u.Engine.Config.Core.ViewRenderer.Load(PageLogin, PageRecoverStart, PageRecoverEnd, PageRegister); err != nil {
		return err
	}

	if err = u.Engine.Config.Core.MailRenderer.Load(EmailConfirmHTML, EmailConfirmTxt, EmailRecoverHTML, EmailRecoverTxt); err != nil {
		return err
	}

	if _, ok := u.Config.Storage.Server.(engine.CreatingServerStorer); !ok {
		return errors.New("register module activated but storer could not be upgraded to CreatingServerStorer")
	}

	// login
	u.Engine.Config.Core.Router.GET("/login", u.Engine.Core.ErrorHandler.Wrap(u.LoginGet))
	u.Engine.Config.Core.Router.POST("/login", u.Engine.Core.ErrorHandler.Wrap(u.LoginPost))

	// logout
	u.Engine.Config.Core.Router.DELETE("/logout", u.Engine.Core.ErrorHandler.Wrap(u.Logout))

	// confirmation
	u.Engine.Config.Core.Router.GET("/confirm", u.Engine.Config.Core.ErrorHandler.Wrap(u.GetConfirm))

	// account recovery
	u.Engine.Config.Core.Router.GET("/recover", u.Core.ErrorHandler.Wrap(u.StartGetRecover))
	u.Engine.Config.Core.Router.POST("/recover", u.Core.ErrorHandler.Wrap(u.StartPostRecover))
	u.Engine.Config.Core.Router.GET("/recover/end", u.Core.ErrorHandler.Wrap(u.EndGetRecover))
	u.Engine.Config.Core.Router.POST("/recover/end", u.Core.ErrorHandler.Wrap(u.EndPostRecover))

	// registration
	u.Config.Core.Router.GET("/register", u.Config.Core.ErrorHandler.Wrap(u.GetRegister))
	u.Config.Core.Router.POST("/register", u.Config.Core.ErrorHandler.Wrap(u.PostRegister))

	// authentication events
	u.Events.After(engine.EventRegister, u.StartConfirmationWeb)

	u.Events.Before(engine.EventAuth, u.PreventAuth)
	u.Events.Before(engine.EventAuth, u.EnsureNotLocked)
	u.Events.Before(engine.EventOAuth2, u.EnsureNotLocked)

	u.Events.After(engine.EventAuth, u.ResetLoginAttempts)
	u.Events.After(engine.EventAuth, u.CreateRememberToken)
	u.Events.After(engine.EventOAuth2, u.CreateRememberToken)
	u.Events.After(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
		refreshExpiry(w)
		return false, nil
	})

	u.Events.After(engine.EventAuthFail, u.UpdateLockAttempts)

	u.Events.After(engine.EventPasswordReset, u.ResetAllTokens)

	return nil
}
