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

// InitAuth :
func (u *Users) InitAuth() (err error) {
	if err = u.Engine.Config.Core.ViewRenderer.Load(PageLogin); err != nil {
		return err
	}

	u.Engine.Config.Core.Router.GET("/login", u.Engine.Core.ErrorHandler.Wrap(u.LoginGet))
	u.Engine.Config.Core.Router.POST("/login", u.Engine.Core.ErrorHandler.Wrap(u.LoginPost))

	return nil
}

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

	pid := creds.GetPID()
	pidUser, err := u.Engine.Storage.Server.Load(r.Context(), pid)
	if err == engine.ErrUserNotFound {
		logger.Infof("failed to load user requested by pid: %s", pid)
		data := engine.HTMLData{engine.DataErr: "Invalid Credentials"}
		return u.Engine.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
	} else if err != nil {
		return err
	}

	authUser := engine.MustBeAuthable(pidUser)
	password := authUser.GetPassword()

	r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyUser, pidUser))

	var handled bool
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(creds.GetPassword()))
	if err != nil {
		handled, err = u.Engine.Events.FireAfter(engine.EventAuthFail, w, r)
		if err != nil {
			return err
		} else if handled {
			return nil
		}

		logger.Infof("user %s failed to log in", pid)
		data := engine.HTMLData{engine.DataErr: "Invalid Credentials"}
		return u.Engine.Core.Responder.Respond(w, r, http.StatusOK, PageLogin, data)
	}

	r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyValues, validatable))

	handled, err = u.Events.FireBefore(engine.EventAuth, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	handled, err = u.Events.FireBefore(engine.EventAuthHijack, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	logger.Infof("user %s logged in", pid)
	engine.PutSession(w, engine.SessionKey, pid)
	engine.DelSession(w, engine.SessionHalfAuthKey)

	handled, err = u.Engine.Events.FireAfter(engine.EventAuth, w, r)
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
