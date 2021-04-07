package users

import (
	"context"
	"net/http"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/bson"
)

// Storage key constants
const (
	StoreAttemptNumber = "attempt_number"
	StoreAttemptTime   = "attempt_time"
	StoreLocked        = "locked"
)

// EnsureNotLocked ensures the account is not locked.
func (u *Users) EnsureNotLocked(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
	return u.updateLockedState(w, r, true)
}

// ResetLoginAttempts resets the attempt number field.
func (u *Users) ResetLoginAttempts(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
	user, err := u.CurrentUser(r)
	if err != nil {
		return false, err
	}

	user.AttemptCount = 0
	user.LastAttempt = time.Now().UTC()

	_, err = u.Engine.Core.Database.Collection(Collection).InsertOne(r.Context(), user)

	return false, err
}

// UpdateLockAttempts adjusts the attempt number and time negatively
// and locks the user if they're beyond limits.
func (u *Users) UpdateLockAttempts(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
	return u.updateLockedState(w, r, false)
}

// updateLockedState exists to minimize any differences between a success and
// a failure path in the case where a correct/incorrect password is entered
func (u *Users) updateLockedState(w http.ResponseWriter, r *http.Request, wasCorrectPassword bool) (bool, error) {
	user, err := u.CurrentUser(r)
	if err != nil {
		return false, err
	}

	// Fetch things
	last := user.LastAttempt
	attempts := user.AttemptCount
	attempts++

	if !wasCorrectPassword {
		if time.Now().UTC().Sub(last) <= u.Config.Authboss.LockWindow {
			if attempts >= u.Config.Authboss.LockAfter {
				user.Locked = time.Now().UTC().Add(u.Config.Authboss.LockDuration)
			}

			user.AttemptCount = attempts
		} else {
			user.AttemptCount = 1
		}
	}
	user.LastAttempt = time.Now().UTC()

	_, err = u.Engine.Core.Database.Collection(Collection).InsertOne(r.Context(), user)
	if err != nil {
		return false, err
	}

	if !IsLocked(user) {
		return false, nil
	}

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		Failure:      "Your account has been locked, please contact the administrator.",
		RedirectPath: "/login",
	}
	return true, u.Engine.Core.Redirector.Redirect(w, r, ro)
}

// Lock a user manually.
func (u *Users) Lock(ctx context.Context, key string) error {
	user, err := GetUser(ctx, u.Core.Database, key)
	if err != nil {
		return err
	}

	_, err = u.Engine.Core.Database.Collection(Collection).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"locked": locked})

	return err
}

// Unlock a user that was locked by this module.
func (u *Users) Unlock(ctx context.Context, key string) error {
	user, err := GetUser(ctx, u.Core.Database, key)
	if err != nil {
		return err
	}

	// Set the last attempt to be -window*2 to avoid immediately
	// giving another login failure. Don't reset Locked to Zero time
	// because some databases may have trouble storing values before
	// unix_time(0): Jan 1st, 1970
	now := time.Now().UTC()
	user.AttemptCount = 0
	user.LastAttempt = now.Add(-u.Engine.Config.Authboss.LockWindow * 2)
	user.Locked = now.Add(-u.Engine.Config.Authboss.LockDuration)

	_, err = u.Engine.Core.Database.Collection(Collection).InsertOne(ctx, user)
	return err
}

// LockMiddleware ensures that a user is not locked, or else it will intercept
// the request and send them to the configured LockNotOK page, this will load
// the user if he's not been loaded yet from the session. And panics if it
// cannot load the user.
func (u *Users) LockMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := u.LoadCurrentUserP(&r)

			if !IsLocked(user) {
				next.ServeHTTP(w, r)
				return
			}

			logger := u.RequestLogger(r)
			logger.Infof("user %s prevented from accessing %s: locked", user.Email, r.URL.Path)
			ro := engine.RedirectOptions{
				Code:         http.StatusTemporaryRedirect,
				Failure:      "Your account has been locked, please contact the administrator.",
				RedirectPath: "/login",
			}
			if err := u.Core.Redirector.Redirect(w, r, ro); err != nil {
				logger.Errorf("error redirecting in lock.Middleware: #%v", err)
			}
		})
	}
}

// IsLocked checks if a user is locked
func IsLocked(u *User) bool {
	return u.Locked.After(time.Now().UTC())
}
