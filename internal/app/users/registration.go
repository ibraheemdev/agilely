package users

import (
	"context"
	"errors"
	"net/http"
	"sort"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Pages
const (
	PageRegister = "register.html.tpl"
)

// GetRegister the register page
func (u *Users) GetRegister(w http.ResponseWriter, req *http.Request) error {
	return u.Core.Responder.Respond(w, req, http.StatusOK, PageRegister, nil)
}

// PostRegister to the register page
func (u *Users) PostRegister(w http.ResponseWriter, req *http.Request) error {
	logger := u.RequestLogger(req)
	validatable, err := u.Core.BodyReader.Read(PageRegister, req)
	if err != nil {
		return err
	}

	var arbitrary map[string]string
	var preserve map[string]string
	if arb, ok := validatable.(engine.ArbitraryValuer); ok {
		arbitrary = arb.GetValues()
		preserve = make(map[string]string)

		for k, v := range arbitrary {
			if hasString(u.Config.Authboss.RegisterPreserveFields, k) {
				preserve[k] = v
			}
		}
	}

	errs := validatable.Validate()
	if errs != nil {
		logger.Info("registration validation failed")
		data := engine.HTMLData{
			engine.DataValidation: engine.ErrorMap(errs),
		}
		if preserve != nil {
			data[engine.DataPreserve] = preserve
		}
		return u.Core.Responder.Respond(w, req, http.StatusOK, PageRegister, data)
	}

	// Get values from request
	userVals := engine.MustHaveUserValues(validatable)
	email, password := userVals.GetPID(), userVals.GetPassword()

	var user User

	pass, err := bcrypt.GenerateFromPassword([]byte(password), u.Config.Authboss.BCryptCost)
	if err != nil {
		return err
	}

	user.Email = email
	user.Password = string(pass)

	if n, ok := arbitrary["name"]; ok {
		user.Name = n
	}

	res := u.Core.Database.Collection(Collection).FindOne(req.Context(), bson.M{"email": email})
	// if the email is taken
	if res.Err() != engine.ErrNoDocuments {
		logger.Infof("user %s attempted to re-register", email)
		errs = []error{errors.New("email is taken")}
		data := engine.HTMLData{
			engine.DataValidation: engine.ErrorMap(errs),
		}
		if preserve != nil {
			data[engine.DataPreserve] = preserve
		}
		return u.Core.Responder.Respond(w, req, http.StatusOK, PageRegister, data)
	}

	_, err = u.Core.Database.Collection(Collection).InsertOne(req.Context(), user)
	if err != nil {
		return err
	}

	req = req.WithContext(context.WithValue(req.Context(), CTXKeyUser, user))
	handled, err := u.AuthEvents.FireAfter(engine.EventRegister, w, req)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	// Log the user in, but only if the response wasn't handled previously
	// by a module like confirm.
	engine.PutSession(w, engine.SessionKey, email)

	logger.Infof("registered and logged in user %s", email)
	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		Success:      "Account successfully created, you are now logged in",
		RedirectPath: "/login",
	}
	return u.Core.Redirector.Redirect(w, req, ro)
}

// hasString checks to see if a sorted (ascending) array of
// strings contains a string
func hasString(arr []string, s string) bool {
	index := sort.SearchStrings(arr, s)
	if index < 0 || index >= len(arr) {
		return false
	}

	return arr[index] == s
}
