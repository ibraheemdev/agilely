package users

import (
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

const (
	// PageLogin is for identifying the login page for parsing & validation
	PageLogin = "login.html.tpl"
)

func init() {
	engine.RegisterModule("auth", &Auth{})
}

// Auth module
type Auth struct {
	*engine.Engine
}

// Init module
func (a *Auth) Init(ab *engine.Engine) (err error) {
	a.Engine = ab

	if err = a.Engine.Config.Core.ViewRenderer.Load(PageLogin); err != nil {
		return err
	}

	a.Engine.Config.Core.Router.GET("/login", a.Engine.Core.ErrorHandler.Wrap(a.LoginGet))
	a.Engine.Config.Core.Router.POST("/login", a.Engine.Core.ErrorHandler.Wrap(a.LoginPost))

	return nil
}

// LoginGet simply displays the login form
func (a *Auth) LoginGet(w http.ResponseWriter, r *http.Request) error {
	data := engine.HTMLData{}
	if redir := r.URL.Query().Get(engine.FormValueRedirect); len(redir) != 0 {
		data[engine.FormValueRedirect] = redir
	}
	return a.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
}

// LoginPost attempts to validate the credentials passed in
// to log in a user.
func (a *Auth) LoginPost(w http.ResponseWriter, r *http.Request) error {
	logger := a.RequestLogger(r)

	validatable, err := a.Engine.Core.BodyReader.Read(PageLogin, r)
	if err != nil {
		return err
	}

	// Skip validation since all the validation happens during the database lookup and
	// password check.
	creds := engine.MustHaveUserValues(validatable)

	pid := creds.GetPID()
	pidUser, err := a.Engine.Storage.Server.Load(r.Context(), pid)
	if err == engine.ErrUserNotFound {
		logger.Infof("failed to load user requested by pid: %s", pid)
		data := engine.HTMLData{engine.DataErr: "Invalid Credentials"}
		return a.Engine.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
	} else if err != nil {
		return err
	}

	authUser := engine.MustBeAuthable(pidUser)
	password := authUser.GetPassword()

	r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyUser, pidUser))

	var handled bool
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.GetPassword()))
	if err != nil {
		handled, err = a.Engine.Events.FireAfter(engine.EventAuthFail, w, r)
		if err != nil {
			return err
		} else if handled {
			return nil
		}

		logger.Infof("user %s failed to log in", pid)
		data := engine.HTMLData{engine.DataErr: "Invalid Credentials"}
		return a.Engine.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
	}

	r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyValues, validatable))

	handled, err = a.Events.FireBefore(engine.EventAuth, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	handled, err = a.Events.FireBefore(engine.EventAuthHijack, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	logger.Infof("user %s logged in", pid)
	engine.PutSession(w, engine.SessionKey, pid)
	engine.DelSession(w, engine.SessionHalfAuthKey)

	handled, err = a.Engine.Events.FireAfter(engine.EventAuth, w, r)
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
	return a.Engine.Core.Redirector.Redirect(w, r, ro)
}
