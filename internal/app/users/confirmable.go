package users

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/pkg/mailer"
)

const (
	// PageConfirm is only really used for the BodyReader
	PageConfirm = "confirm.html.tpl"

	// EmailConfirmHTML is the name of the html template for e-mails
	EmailConfirmHTML = "confirm.html.tpl"
	// EmailConfirmTxt is the name of the text template for e-mails
	EmailConfirmTxt = "confirm.text.tpl"

	// FormValueConfirm is the name of the form value for
	FormValueConfirm = "cnf"

	// DataConfirmURL is the name of the e-mail template variable
	// that gives the url to send to the user for confirmation.
	DataConfirmURL = "url"

	confirmTokenSize  = 64
	confirmTokenSplit = confirmTokenSize / 2
)

// InitConfirm module
func (u *Users) InitConfirm(e *engine.Engine) (err error) {
	if err = u.Engine.Config.Core.MailRenderer.Load(EmailConfirmHTML, EmailConfirmTxt); err != nil {
		return err
	}

	u.Engine.Config.Core.Router.GET("/confirm", u.Engine.Config.Core.ErrorHandler.Wrap(u.GetConfirm))

	u.Events.Before(engine.EventAuth, u.PreventAuth)
	u.Events.After(engine.EventRegister, u.StartConfirmationWeb)

	return nil
}

// PreventAuth stops the EventAuth from succeeding when a user is not confirmed
// This relies on the fact that the context holds the user at this point in time
// loaded by the auth module (or something else).
func (u *Users) PreventAuth(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
	logger := u.Engine.RequestLogger(r)

	user, err := u.Engine.CurrentUser(r)
	if err != nil {
		return false, err
	}

	cuser := engine.MustBeConfirmable(user)
	if cuser.GetConfirmed() {
		logger.Infof("user %s is confirmed, allowing auth", user.GetPID())
		return false, nil
	}

	logger.Infof("user %s was not confirmed, preventing auth", user.GetPID())
	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: "/login",
		Failure:      "Your account has not been confirmed, please check your e-mail.",
	}
	return true, u.Engine.Config.Core.Redirector.Redirect(w, r, ro)
}

// StartConfirmationWeb hijacks a request and forces a user to be confirmed
// first it's assumed that the current user is loaded into the request context.
func (u *Users) StartConfirmationWeb(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
	user, err := u.Engine.CurrentUser(r)
	if err != nil {
		return false, err
	}

	cuser := engine.MustBeConfirmable(user)
	if err = u.StartConfirmation(r.Context(), cuser, true); err != nil {
		return false, err
	}

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: "/login",
		Success:      "Please verify your account, an e-mail has been sent to you.",
	}
	return true, u.Engine.Config.Core.Redirector.Redirect(w, r, ro)
}

// StartConfirmation begins confirmation on a user by setting them to require
// confirmation via a created token, and optionally sending them an e-mail.
func (u *Users) StartConfirmation(ctx context.Context, user engine.ConfirmableUser, sendEmail bool) error {
	logger := u.Engine.Logger(ctx)

	selector, verifier, token, err := GenerateConfirmCreds()
	if err != nil {
		return err
	}

	user.PutConfirmed(false)
	user.PutConfirmSelector(selector)
	user.PutConfirmVerifier(verifier)

	logger.Infof("generated new confirm token for user: %s", user.GetPID())
	if err := u.Engine.Config.Storage.Server.Save(ctx, user); err != nil {
		return fmt.Errorf("%w failed to save user during StartConfirmation, user data may be in weird state", err)
	}

	go u.SendConfirmEmail(ctx, user.GetEmail(), token)

	return nil
}

// SendConfirmEmail sends a confirmation e-mail to a user
func (u *Users) SendConfirmEmail(ctx context.Context, to, token string) {
	logger := u.Engine.Logger(ctx)

	mailConfirmURL := u.mailConfirmURL(token)

	email := mailer.Email{
		To:       []string{to},
		From:     u.Config.Mail.From,
		FromName: u.Config.Mail.FromName,
		Subject:  u.Config.Mail.SubjectPrefix + "Confirm New Account",
	}

	logger.Infof("sending confirm e-mail to: %s", to)

	ro := engine.EmailResponseOptions{
		Data:         engine.NewHTMLData(DataConfirmURL, mailConfirmURL),
		HTMLTemplate: EmailConfirmHTML,
		TextTemplate: EmailConfirmTxt,
	}
	if err := u.Engine.Email(ctx, email, ro); err != nil {
		logger.Errorf("failed to send confirm e-mail to %s: %+v", to, err)
	}
}

// GetConfirm is a request that confirms a user with a valid token
func (u *Users) GetConfirm(w http.ResponseWriter, r *http.Request) error {
	logger := u.RequestLogger(r)

	validator, err := u.Engine.Config.Core.BodyReader.Read(PageConfirm, r)
	if err != nil {
		return err
	}

	if errs := validator.Validate(); errs != nil {
		logger.Infof("validation failed in Confirm.Get, this typically means a bad token: %+v", errs)
		return u.invalidConfirmToken(w, r)
	}

	values := engine.MustHaveConfirmValues(validator)

	rawToken, err := base64.URLEncoding.DecodeString(values.GetToken())
	if err != nil {
		logger.Infof("error decoding token in Confirm.Get, this typically means a bad token: %s %+v", values.GetToken(), err)
		return u.invalidConfirmToken(w, r)
	}

	if len(rawToken) != confirmTokenSize {
		logger.Infof("invalid confirm token submitted, size was wrong: %d", len(rawToken))
		return u.invalidConfirmToken(w, r)
	}

	selectorBytes := sha512.Sum512(rawToken[:confirmTokenSplit])
	verifierBytes := sha512.Sum512(rawToken[confirmTokenSplit:])
	selector := base64.StdEncoding.EncodeToString(selectorBytes[:])

	storer := engine.EnsureCanConfirm(u.Engine.Config.Storage.Server)
	user, err := storer.LoadByConfirmSelector(r.Context(), selector)
	if err == engine.ErrUserNotFound {
		logger.Infof("confirm selector was not found in database: %s", selector)
		return u.invalidConfirmToken(w, r)
	} else if err != nil {
		return err
	}

	dbVerifierBytes, err := base64.StdEncoding.DecodeString(user.GetConfirmVerifier())
	if err != nil {
		logger.Infof("invalid confirm verifier stored in database: %s", user.GetConfirmVerifier())
		return u.invalidConfirmToken(w, r)
	}

	if subtle.ConstantTimeEq(int32(len(verifierBytes)), int32(len(dbVerifierBytes))) != 1 ||
		subtle.ConstantTimeCompare(verifierBytes[:], dbVerifierBytes) != 1 {
		logger.Info("stored confirm verifier does not match provided one")
		return u.invalidConfirmToken(w, r)
	}

	user.PutConfirmSelector("")
	user.PutConfirmVerifier("")
	user.PutConfirmed(true)

	logger.Infof("user %s confirmed their account", user.GetPID())
	if err = u.Engine.Config.Storage.Server.Save(r.Context(), user); err != nil {
		return err
	}

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		Success:      "You have successfully confirmed your account.",
		RedirectPath: "/login",
	}
	return u.Engine.Config.Core.Redirector.Redirect(w, r, ro)
}

func (u *Users) mailConfirmURL(token string) string {
	query := url.Values{FormValueConfirm: []string{token}}

	if len(u.Config.Mail.RootURL) != 0 {
		return fmt.Sprintf("%s?%s", u.Config.Mail.RootURL+"/confirm", query.Encode())
	}

	p := path.Join(u.Config.Paths.Mount, "confirm")
	return fmt.Sprintf("%s%s?%s", u.Config.Paths.RootURL, p, query.Encode())
}

func (u *Users) invalidConfirmToken(w http.ResponseWriter, r *http.Request) error {
	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		Failure:      "confirm token is invalid",
		RedirectPath: "/login",
	}
	return u.Engine.Config.Core.Redirector.Redirect(w, r, ro)
}

// Middleware ensures that a user is confirmed, or else it will intercept the
// request and send them to the confirm page, this will load the user if he's
// not been loaded yet from the session.
//
// Panics if the user was not able to be loaded in order to allow a panic
// handler to show a nice error page, also panics if it failed to redirect
// for whatever reason.
func Middleware(ab *engine.Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := ab.LoadCurrentUserP(&r)

			cu := engine.MustBeConfirmable(user)
			if cu.GetConfirmed() {
				next.ServeHTTP(w, r)
				return
			}

			logger := ab.RequestLogger(r)
			logger.Infof("user %s prevented from accessing %s: not confirmed", user.GetPID(), r.URL.Path)
			ro := engine.RedirectOptions{
				Code:         http.StatusTemporaryRedirect,
				Failure:      "Your account has not been confirmed, please check your e-mail.",
				RedirectPath: "/login",
			}
			if err := ab.Config.Core.Redirector.Redirect(w, r, ro); err != nil {
				logger.Errorf("error redirecting in confirm.Middleware: #%v", err)
			}
		})
	}
}

// GenerateConfirmCreds generates pieces needed for user confirm
// selector: hash of the first half of a 64 byte value
// (to be stored in the database and used in SELECT query)
// verifier: hash of the second half of a 64 byte value
// (to be stored in database but never used in SELECT query)
// token: the user-facing base64 encoded selector+verifier
func GenerateConfirmCreds() (selector, verifier, token string, err error) {
	rawToken := make([]byte, confirmTokenSize)
	if _, err = io.ReadFull(rand.Reader, rawToken); err != nil {
		return "", "", "", err
	}
	selectorBytes := sha512.Sum512(rawToken[:confirmTokenSplit])
	verifierBytes := sha512.Sum512(rawToken[confirmTokenSplit:])

	return base64.StdEncoding.EncodeToString(selectorBytes[:]),
		base64.StdEncoding.EncodeToString(verifierBytes[:]),
		base64.URLEncoding.EncodeToString(rawToken),
		nil
}
