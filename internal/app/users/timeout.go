package users

import (
	"context"
	"net/http"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

var nowTime = time.Now

// InitTimeout :
//
// This installs a hook into the login process so that the
// LastAction is recorded immediately.
func (u *Users) InitTimeout() error {
	u.Events.After(engine.EventAuth, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
		refreshExpiry(w)
		return false, nil
	})

	return nil
}

// TimeToExpiry returns zero if the user session is expired else the time
// until expiry. Takes in the allowed idle duration.
func TimeToExpiry(r *http.Request, expireAfter time.Duration) time.Duration {
	return timeToExpiry(r, expireAfter)
}

func timeToExpiry(r *http.Request, expireAfter time.Duration) time.Duration {
	dateStr, ok := engine.GetSession(r, engine.SessionLastAction)
	if !ok {
		return expireAfter
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		panic("last_action is not a valid RFC3339 date")
	}

	remaining := date.Add(expireAfter).Sub(nowTime().UTC())
	if remaining > 0 {
		return remaining
	}

	return 0
}

// RefreshExpiry updates the last action for the user, so he doesn't
// become expired.
func RefreshExpiry(w http.ResponseWriter, r *http.Request) {
	refreshExpiry(w)
}

func refreshExpiry(w http.ResponseWriter) {
	engine.PutSession(w, engine.SessionLastAction, nowTime().UTC().Format(time.RFC3339))
}

type expireMiddleware struct {
	expireAfter      time.Duration
	next             http.Handler
	sessionWhitelist []string
}

// TimeoutMiddleware ensures that the user's expiry information is kept up-to-date
// on each request. Deletes the SessionKey from the session if the user is
// expired (a.ExpireAfter duration since SessionLastAction).
// This middleware conflicts with use of the Remember module, don't enable both
// at the same time.
func TimeoutMiddleware(e *engine.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return expireMiddleware{
			expireAfter:      e.Config.Modules.ExpireAfter,
			next:             next,
			sessionWhitelist: e.Config.Storage.SessionStateWhitelistKeys,
		}
	}
}

// ServeHTTP removes the session and hides the loaded user from the handlers
// below it.
func (m expireMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := engine.GetSession(r, engine.SessionKey); ok {
		ttl := timeToExpiry(r, m.expireAfter)

		if ttl == 0 {
			engine.DelAllSession(w, m.sessionWhitelist)
			engine.DelSession(w, engine.SessionKey)
			engine.DelSession(w, engine.SessionLastAction)
			ctx := context.WithValue(r.Context(), engine.CTXKeyPID, nil)
			ctx = context.WithValue(ctx, engine.CTXKeyUser, nil)

			ctxState := r.Context().Value(engine.CTXKeySessionState)
			if ctxState != nil {
				state := ctxState.(engine.ClientState)
				whitelist := make(map[string]struct{})

				for _, w := range m.sessionWhitelist {
					whitelist[w] = struct{}{}
				}

				newState := stateHider{cs: state, whitelist: whitelist}
				ctx = context.WithValue(ctx, engine.CTXKeySessionState, newState)
			}

			r = r.WithContext(ctx)
		} else {
			refreshExpiry(w)
		}
	}

	m.next.ServeHTTP(w, r)
}

type stateHider struct {
	whitelist map[string]struct{}
	cs        engine.ClientState
}

func (k stateHider) Get(s string) (string, bool) {
	_, ok := k.whitelist[s]
	if !ok {
		return "", false
	}

	return k.cs.Get(s)
}
