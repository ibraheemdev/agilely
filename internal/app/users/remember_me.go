package users

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

const (
	nNonceSize = 32
)

// InitRemember :
func (u *Users) InitRemember() error {

	u.Events.After(engine.EventAuth, u.CreateRememberToken)
	u.Events.After(engine.EventOAuth2, u.CreateRememberToken)
	u.Events.After(engine.EventPasswordReset, u.ResetAllTokens)

	return nil
}

// CreateRememberToken creates a remember token and saves it in the user's cookies.
func (u *Users) CreateRememberToken(w http.ResponseWriter, req *http.Request, handled bool) (bool, error) {
	rmIntf := req.Context().Value(engine.CTXKeyValues)
	if rmIntf == nil {
		return false, nil
	} else if rm, ok := rmIntf.(engine.RememberValuer); !ok || !rm.GetShouldRemember() {
		return false, nil
	}

	user := u.Engine.CurrentUserP(req)
	hash, token, err := GenerateToken(user.GetPID())
	if err != nil {
		return false, err
	}

	storer := engine.EnsureCanRemember(u.Engine.Config.Storage.Server)
	if err = storer.AddRememberToken(req.Context(), user.GetPID(), hash); err != nil {
		return false, err
	}

	engine.PutCookie(w, engine.CookieRemember, token)

	return false, nil
}

// RememberMiddleware automatically authenticates users if they have remember me tokens
// If the user has been loaded already, it returns early
func RememberMiddleware(e *engine.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Safely can ignore error here
			if id, _ := e.CurrentUserID(r); len(id) == 0 {
				if err := Authenticate(e, w, &r); err != nil {
					logger := e.RequestLogger(r)
					logger.Errorf("failed to authenticate user via remember me: %+v", err)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Authenticate the user using their remember cookie.
// If the cookie proves unusable it will be deleted. A cookie
// may be unusable for the following reasons:
// - Can't decode the base64
// - Invalid token format
// - Can't find token in DB
//
// In order to authenticate it adds to the request context as well as to the
// cookie and session states.
func Authenticate(e *engine.Engine, w http.ResponseWriter, req **http.Request) error {
	logger := e.RequestLogger(*req)
	cookie, ok := engine.GetCookie(*req, engine.CookieRemember)
	if !ok {
		return nil
	}

	rawToken, err := base64.URLEncoding.DecodeString(cookie)
	if err != nil {
		engine.DelCookie(w, engine.CookieRemember)
		logger.Infof("failed to decode remember me cookie, deleting cookie")
		return nil
	}

	index := bytes.IndexByte(rawToken, ';')
	if index < 0 {
		engine.DelCookie(w, engine.CookieRemember)
		logger.Infof("failed to decode remember me token, deleting cookie")
		return nil
	}

	pid := string(rawToken[:index])
	sum := sha512.Sum512(rawToken)
	hash := base64.StdEncoding.EncodeToString(sum[:])

	storer := engine.EnsureCanRemember(e.Config.Storage.Server)
	err = storer.UseRememberToken((*req).Context(), pid, hash)
	switch {
	case err == engine.ErrTokenNotFound:
		logger.Infof("remember me cookie had a token that was not in storage, deleting cookie")
		engine.DelCookie(w, engine.CookieRemember)
		return nil
	case err != nil:
		return err
	}

	hash, token, err := GenerateToken(pid)
	if err != nil {
		return err
	}

	if err = storer.AddRememberToken((*req).Context(), pid, hash); err != nil {
		return fmt.Errorf("failed to save remember me token %w", err)
	}

	*req = (*req).WithContext(context.WithValue((*req).Context(), engine.CTXKeyPID, pid))
	engine.PutSession(w, engine.SessionKey, pid)
	engine.PutSession(w, engine.SessionHalfAuthKey, "true")
	engine.DelCookie(w, engine.CookieRemember)
	engine.PutCookie(w, engine.CookieRemember, token)

	return nil
}

// ResetAllTokens is called after the password has been reset, since
// it should invalidate all tokens associated to that user.
func (u *Users) ResetAllTokens(w http.ResponseWriter, req *http.Request, handled bool) (bool, error) {
	user, err := u.Engine.CurrentUser(req)
	if err != nil {
		return false, err
	}

	logger := u.Engine.RequestLogger(req)
	storer := engine.EnsureCanRemember(u.Engine.Config.Storage.Server)

	pid := user.GetPID()
	engine.DelCookie(w, engine.CookieRemember)

	logger.Infof("deleting tokens and rm cookies for user %s due to password reset", pid)

	return false, storer.DelRememberTokens(req.Context(), pid)
}

// GenerateToken creates a remember me token
func GenerateToken(pid string) (hash string, token string, err error) {
	rawToken := make([]byte, nNonceSize+len(pid)+1)
	copy(rawToken, pid)
	rawToken[len(pid)] = ';'

	if _, err := io.ReadFull(rand.Reader, rawToken[len(pid)+1:]); err != nil {
		return "", "", fmt.Errorf("%wfailed to create remember me nonce", err)
	}

	sum := sha512.Sum512(rawToken)
	return base64.StdEncoding.EncodeToString(sum[:]), base64.URLEncoding.EncodeToString(rawToken), nil
}
