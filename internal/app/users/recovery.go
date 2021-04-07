package users

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/pkg/mailer"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Constants for templates etc.
const (
	DataRecoverToken = "recover_token"
	DataRecoverURL   = "recover_url"

	FormValueToken = "token"

	EmailRecoverHTML = "recover.html.tpl"
	EmailRecoverTxt  = "recover.text.tpl"

	PageRecoverStart  = "recover_start.html.tpl"
	PageRecoverMiddle = "recover_middle.html.tpl"
	PageRecoverEnd    = "recover_end.html.tpl"

	recoverInitiateSuccessFlash = "An email has been sent to you with further instructions on how to reset your password."

	recoverTokenSize  = 64
	recoverTokenSplit = recoverTokenSize / 2
)

// StartGetRecover starts the recover procedure by rendering a form for the user.
func (u *Users) StartGetRecover(w http.ResponseWriter, req *http.Request) error {
	return u.Engine.Core.Responder.Respond(w, req, http.StatusOK, PageRecoverStart, nil)
}

// StartPostRecover starts the recover procedure using values provided from the user
// usually from the StartGet's form.
func (u *Users) StartPostRecover(w http.ResponseWriter, req *http.Request) error {
	logger := u.RequestLogger(req)

	validatable, err := u.Engine.Core.BodyReader.Read(PageRecoverStart, req)
	if err != nil {
		return err
	}

	if errs := validatable.Validate(); errs != nil {
		logger.Info("recover validation failed")
		data := engine.HTMLData{engine.DataValidation: engine.ErrorMap(errs)}
		return u.Engine.Core.Responder.Respond(w, req, http.StatusOK, PageRecoverStart, data)
	}

	recoverVals := engine.MustHaveRecoverStartValues(validatable)

	user, err := GetUser(req.Context(), u.Core.Database, recoverVals.GetPID())
	if err == engine.ErrNoDocuments {
		logger.Infof("user %s was attempted to be recovered, user does not exist, faking successful response", recoverVals.GetPID())
		ro := engine.RedirectOptions{
			Code:         http.StatusTemporaryRedirect,
			RedirectPath: "/",
			Success:      recoverInitiateSuccessFlash,
		}
		return u.Engine.Core.Redirector.Redirect(w, req, ro)
	}

	selector, verifier, token, err := GenerateRecoverCreds()
	if err != nil {
		return err
	}

	user.RecoverSelector = selector
	user.RecoverVerifier = verifier
	user.RecoverTokenExpiry = time.Now().UTC().Add(u.Config.Authboss.RecoverTokenDuration)

	_, err = u.Engine.Core.Database.Collection(Collection).InsertOne(req.Context(), user)
	if err != nil {
		return err
	}

	go u.SendRecoverEmail(req.Context(), user.Email, token)

	logger.Infof("user %s password recovery initiated", user.Email)
	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: "/",
		Success:      recoverInitiateSuccessFlash,
	}
	return u.Engine.Core.Redirector.Redirect(w, req, ro)
}

// SendRecoverEmail to a specific e-mail address passing along the encodedToken
// in an escaped URL to the templates.
func (u *Users) SendRecoverEmail(ctx context.Context, to, encodedToken string) {
	logger := u.Engine.Logger(ctx)

	mailRecoverURL := u.mailRecoverURL(encodedToken)

	email := mailer.Email{
		To:       []string{to},
		From:     u.Engine.Config.Mail.From,
		FromName: u.Engine.Config.Mail.FromName,
		Subject:  u.Engine.Config.Mail.SubjectPrefix + "Password Reset",
	}

	ro := engine.EmailResponseOptions{
		HTMLTemplate: EmailRecoverHTML,
		TextTemplate: EmailRecoverTxt,
		Data: engine.HTMLData{
			DataRecoverURL: mailRecoverURL,
		},
	}

	logger.Infof("sending recover e-mail to: %s", to)
	if err := u.Engine.Email(ctx, email, ro); err != nil {
		logger.Errorf("failed to recover send e-mail to %s: %+v", to, err)
	}
}

// EndGetRecover shows a password recovery form, and it should have the token that
// the user brought in the query parameters in it on submission.
func (u *Users) EndGetRecover(w http.ResponseWriter, req *http.Request) error {
	validatable, err := u.Engine.Core.BodyReader.Read(PageRecoverMiddle, req)
	if err != nil {
		return err
	}

	values := engine.MustHaveRecoverMiddleValues(validatable)
	token := values.GetToken()

	data := engine.HTMLData{
		DataRecoverToken: token,
	}

	return u.Engine.Core.Responder.Respond(w, req, http.StatusOK, PageRecoverEnd, data)
}

// EndPostRecover retrieves the token
func (u *Users) EndPostRecover(w http.ResponseWriter, req *http.Request) error {
	logger := u.RequestLogger(req)

	validatable, err := u.Engine.Core.BodyReader.Read(PageRecoverEnd, req)
	if err != nil {
		return err
	}

	values := engine.MustHaveRecoverEndValues(validatable)
	password := values.GetPassword()
	token := values.GetToken()

	if errs := validatable.Validate(); errs != nil {
		logger.Info("recovery validation failed")
		data := engine.HTMLData{
			engine.DataValidation: engine.ErrorMap(errs),
			DataRecoverToken:      token,
		}
		return u.Core.Responder.Respond(w, req, http.StatusOK, PageRecoverEnd, data)
	}

	rawToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		logger.Infof("invalid recover token submitted, base64 decode failed: %+v", err)
		return u.invalidRecoverToken(PageRecoverEnd, w, req)
	}

	if len(rawToken) != recoverTokenSize {
		logger.Infof("invalid recover token submitted, size was wrong: %d", len(rawToken))
		return u.invalidRecoverToken(PageRecoverEnd, w, req)
	}

	selectorBytes := sha512.Sum512(rawToken[:recoverTokenSplit])
	verifierBytes := sha512.Sum512(rawToken[recoverTokenSplit:])
	selector := base64.StdEncoding.EncodeToString(selectorBytes[:])

	var user User
	err = u.Core.Database.Collection(Collection).FindOne(req.Context(), bson.M{"recover_selector": selector}).Decode(&user)
	if err == engine.ErrNoDocuments {
		logger.Info("invalid recover token submitted, user not found")
		return u.invalidRecoverToken(PageRecoverEnd, w, req)
	} else if err != nil {
		return err
	}

	if time.Now().UTC().After(user.RecoverTokenExpiry) {
		logger.Infof("invalid recover token submitted, already expired: %+v", err)
		return u.invalidRecoverToken(PageRecoverEnd, w, req)
	}

	dbVerifierBytes, err := base64.StdEncoding.DecodeString(user.RecoverVerifier)
	if err != nil {
		logger.Infof("invalid recover verifier stored in database: %s", user.RecoverVerifier)
		return u.invalidRecoverToken(PageRecoverEnd, w, req)
	}

	if subtle.ConstantTimeEq(int32(len(verifierBytes)), int32(len(dbVerifierBytes))) != 1 ||
		subtle.ConstantTimeCompare(verifierBytes[:], dbVerifierBytes) != 1 {
		logger.Info("stored recover verifier does not match provided one")
		return u.invalidRecoverToken(PageRecoverEnd, w, req)
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), u.Engine.Config.Authboss.BCryptCost)
	if err != nil {
		return err
	}

	user.Password = string(pass)
	user.RecoverSelector = ""                  // Don't allow another recovery
	user.RecoverVerifier = ""                  // Don't allow another recovery
	user.RecoverTokenExpiry = time.Now().UTC() // Put current time for those DBs that can't handle 0 time

	_, err = u.Engine.Core.Database.Collection(Collection).InsertOne(req.Context(), user)
	if err != nil {
		return err
	}

	// login user
	// engine.PutSession(w, engine.SessionKey, user.GetPID())

	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		RedirectPath: "/",
		Success:      "Successfully updated password",
	}
	return u.Engine.Core.Redirector.Redirect(w, req, ro)
}

func (u *Users) invalidRecoverToken(page string, w http.ResponseWriter, req *http.Request) error {
	errorsAll := []error{errors.New("recovery token is invalid")}
	data := engine.HTMLData{engine.DataValidation: engine.ErrorMap(errorsAll)}
	return u.Engine.Core.Responder.Respond(w, req, http.StatusOK, PageRecoverEnd, data)
}

func (u *Users) mailRecoverURL(token string) string {
	query := url.Values{FormValueToken: []string{token}}

	if len(u.Config.Mail.RootURL) != 0 {
		return fmt.Sprintf("%s?%s", u.Config.Mail.RootURL+"/recover/end", query.Encode())
	}

	p := path.Join(u.Config.Mount, "recover/end")
	return fmt.Sprintf("%s%s?%s", u.Config.RootURL, p, query.Encode())
}

// GenerateRecoverCreds generates pieces needed for user recovery
// selector: hash of the first half of a 64 byte value
// (to be stored in the database and used in SELECT query)
// verifier: hash of the second half of a 64 byte value
// (to be stored in database but never used in SELECT query)
// token: the user-facing base64 encoded selector+verifier
func GenerateRecoverCreds() (selector, verifier, token string, err error) {
	rawToken := make([]byte, recoverTokenSize)
	if _, err = io.ReadFull(rand.Reader, rawToken); err != nil {
		return "", "", "", err
	}
	selectorBytes := sha512.Sum512(rawToken[:recoverTokenSplit])
	verifierBytes := sha512.Sum512(rawToken[recoverTokenSplit:])

	return base64.StdEncoding.EncodeToString(selectorBytes[:]),
		base64.StdEncoding.EncodeToString(verifierBytes[:]),
		base64.URLEncoding.EncodeToString(rawToken),
		nil
}
