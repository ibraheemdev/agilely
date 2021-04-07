// Package oauth allows users to be created and authenticated
// via oauth services like facebook, google etc. Currently
// only the web server flow is supported.
//
// The general flow looks like this:
//   1. User goes to Start handler and has his session packed with goodies
//      then redirects to the OAuth service.
//   2. OAuth service returns to OAuthCallback which extracts state and
//      parameters and generally checks that everything is ok. It uses the
//      token received to get an access token from the oauth2 library
//   3. Calls the OAuth2Provider.FindUserDetails which should return the user's
//      details in a generic form.
//   4. Passes the user details into the OAuth2ServerStorer.NewFromOAuth2 in
//      order to create a user object we can work with.
//   5. Saves the user in the database, logs them in, redirects.
//
// In order to do this there are a number of parts:
//   1. The configuration of a provider
//      (handled by engine.Config.Authboss.OAuth2Providers).
//   2. The flow of redirection of client, parameter passing etc
//      (handled by this package)
//   3. The HTTP call to the service once a token has been retrieved to
//      get user details (handled by OAuth2Provider.FindUserDetails)
//   4. The creation of a user from the user details returned from the
//      FindUserDetails (engine.OAuth2ServerStorer)
//   5. The special casing of the ServerStorer implementation's Load()
//      function to deal properly with incoming OAuth2 pids. See
//      engine.ParseOAuth2PID as a way to do this.
//
// Of these parts, the responsibility of the engine library consumer
// is on 1, 3, 4, and 5. Configuration of providers that should be used is
// totally up to the consumer. The FindUserDetails function is typically up to
// the user, but we have some basic ones included in this package too.
// The creation of users from the FindUserDetail's map[string]string return
// is handled as part of the implementation of the OAuth2ServerStorer.
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/oauth2"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/internal/app/users"
)

// FormValue constants
const (
	FormValueOAuth2State = "state"
	FormValueOAuth2Redir = "redir"
)

var (
	errOAuthStateValidation = errors.New("could not validate oauth2 state param")
)

// OAuth2 module
type OAuth2 struct {
	*engine.Engine
}

// Init module
func (o *OAuth2) Init(e *engine.Engine) error {
	o.Engine = e

	// Do annoying sorting on keys so we can have predictable
	// route registration (both for consistency inside the router but
	// also for tests -_-)
	var keys []string
	for k := range o.Engine.Config.Authboss.OAuth2Providers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, provider := range keys {
		cfg := o.Engine.Config.Authboss.OAuth2Providers[provider]
		provider = strings.ToLower(provider)

		init := fmt.Sprintf("/oauth2/%s", provider)
		callback := fmt.Sprintf("/oauth2/callback/%s", provider)

		o.Engine.Core.Router.GET(init, o.Engine.Core.ErrorHandler.Wrap(o.Start))
		o.Engine.Core.Router.GET(callback, o.Engine.Core.ErrorHandler.Wrap(o.End))

		if mount := o.Engine.Config.Mount; len(mount) > 0 {
			callback = path.Join(mount, callback)
		}

		cfg.OAuth2Config.RedirectURL = o.Engine.Config.RootURL + callback
	}

	return nil
}

// Start the oauth2 process
func (o *OAuth2) Start(w http.ResponseWriter, r *http.Request) error {
	logger := o.Engine.RequestLogger(r)

	provider := strings.ToLower(filepath.Base(r.URL.Path))
	logger.Infof("started oauth2 flow for provider: %s", provider)
	cfg, ok := o.Engine.Config.Authboss.OAuth2Providers[provider]
	if !ok {
		return fmt.Errorf("oauth2 provider %q not found", provider)
	}

	// Create nonce
	nonce := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("%w failed to create nonce", err)
	}

	state := base64.URLEncoding.EncodeToString(nonce)
	engine.PutSession(w, users.SessionOAuth2State, state)

	// This clearly ignores the fact that query parameters can have multiple
	// values but I guess we're ignoring that
	passAlongs := make(map[string]string)
	for k, vals := range r.URL.Query() {
		for _, val := range vals {
			passAlongs[k] = val
		}
	}

	if len(passAlongs) > 0 {
		byt, err := json.Marshal(passAlongs)
		if err != nil {
			return err
		}
		engine.PutSession(w, users.SessionOAuth2Params, string(byt))
	} else {
		engine.DelSession(w, users.SessionOAuth2Params)
	}

	authCodeURL := cfg.OAuth2Config.AuthCodeURL(state)

	extraParams := cfg.AdditionalParams.Encode()
	if len(extraParams) > 0 {
		authCodeURL = fmt.Sprintf("%s&%s", authCodeURL, extraParams)
	}

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: authCodeURL,
	}
	return o.Engine.Core.Redirector.Redirect(w, r, ro)
}

// for testing, mocked out at the beginning
var exchanger = (*oauth2.Config).Exchange

// End the oauth2 process, this is the handler for the oauth2 callback
// that the third party will redirect to.
func (o *OAuth2) End(w http.ResponseWriter, r *http.Request) error {
	logger := o.Engine.RequestLogger(r)
	provider := strings.ToLower(filepath.Base(r.URL.Path))
	logger.Infof("finishing oauth2 flow for provider: %s", provider)

	// This shouldn't happen because the router should 404 first, but just in case
	cfg, ok := o.Engine.Config.Authboss.OAuth2Providers[provider]
	if !ok {
		return fmt.Errorf("oauth2 provider %q not found", provider)
	}

	wantState, ok := engine.GetSession(r, users.SessionOAuth2State)
	if !ok {
		return errors.New("oauth2 endpoint hit without session state")
	}

	// Verify we got the same state in the session as was passed to us in the
	// query parameter.
	state := r.FormValue(FormValueOAuth2State)
	if state != wantState {
		return errOAuthStateValidation
	}

	rawParams, ok := engine.GetSession(r, users.SessionOAuth2Params)
	var params map[string]string
	if ok {
		if err := json.Unmarshal([]byte(rawParams), &params); err != nil {
			return fmt.Errorf("%w failed to decode oauth2 params", err)
		}
	}

	engine.DelSession(w, users.SessionOAuth2State)
	engine.DelSession(w, users.SessionOAuth2Params)

	hasErr := r.FormValue("error")
	if len(hasErr) > 0 {
		reason := r.FormValue("error_reason")
		logger.Infof("oauth2 login failed: %s, reason: %s", hasErr, reason)

		handled, err := o.Engine.AuthEvents.FireAfter(engine.EventOAuth2Fail, w, r)
		if err != nil {
			return err
		} else if handled {
			return nil
		}

		ro := engine.RedirectOptions{
			Code:         http.StatusTemporaryRedirect,
			RedirectPath: "/",
			Failure:      fmt.Sprintf("%s login cancelled or failed", strings.Title(provider)),
		}
		return o.Engine.Core.Redirector.Redirect(w, r, ro)
	}

	// Get the code which we can use to make an access token
	code := r.FormValue("code")
	token, err := exchanger(cfg.OAuth2Config, r.Context(), code)
	if err != nil {
		return fmt.Errorf("%w could not validate oauth2 code", err)
	}

	details, err := cfg.FindUserDetails(r.Context(), *cfg.OAuth2Config, token)
	if err != nil {
		return err
	}

	storer := engine.EnsureCanOAuth2(o.Engine.Core.Server)
	user, err := storer.NewFromOAuth2(r.Context(), provider, details)
	if err != nil {
		return fmt.Errorf("%w failed to create oauth2 user from values", err)
	}

	user.PutOAuth2Provider(provider)
	user.PutOAuth2AccessToken(token.AccessToken)
	user.PutOAuth2Expiry(token.Expiry)
	if len(token.RefreshToken) != 0 {
		user.PutOAuth2RefreshToken(token.RefreshToken)
	}

	if err := storer.SaveOAuth2(r.Context(), user); err != nil {
		return err
	}

	r = r.WithContext(context.WithValue(r.Context(), users.CTXKeyUser, user))

	handled, err := o.Engine.AuthEvents.FireBefore(engine.EventOAuth2, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	// Fully log user in
	// TODO : Use email instead of oauth2id
	engine.PutSession(w, engine.SessionKey, engine.MakeOAuth2PID(provider, user.GetOAuth2UID()))
	engine.DelSession(w, users.SessionHalfAuthKey)

	// Create a query string from all the pieces we've received
	// as passthru from the original request.
	redirect := "/"
	query := make(url.Values)
	for k, v := range params {
		switch k {
		case users.CookieRemember:
			if v == "true" {
				r = r.WithContext(context.WithValue(r.Context(), engine.CTXKeyValues, RMTrue{}))
			}
		case FormValueOAuth2Redir:
			redirect = v
		default:
			query.Set(k, v)
		}
	}

	handled, err = o.Engine.AuthEvents.FireAfter(engine.EventOAuth2, w, r)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	if len(query) > 0 {
		redirect = fmt.Sprintf("%s?%s", redirect, query.Encode())
	}

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: redirect,
		Success:      fmt.Sprintf("Logged in successfully with %s.", strings.Title(provider)),
	}
	return o.Engine.Core.Redirector.Redirect(w, r, ro)
}

// RMTrue is a dummy struct implementing engine.RememberValuer
// in order to tell the remember me module to remember them.
type RMTrue struct{}

// GetShouldRemember always returns true
func (RMTrue) GetShouldRemember() bool { return true }
