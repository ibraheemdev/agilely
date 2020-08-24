package users

import (
	"context"
	"errors"
	"net/http"
	"sort"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"golang.org/x/crypto/bcrypt"
)

// Pages
const (
	PageRegister = "register.html.tpl"
)

func init() {
	engine.RegisterModule("register", &Register{})
}

// Register module.
type Register struct {
	*engine.Engine
}

// Init the module.
func (r *Register) Init(ab *engine.Engine) (err error) {
	r.Engine = ab

	if _, ok := ab.Config.Storage.Server.(engine.CreatingServerStorer); !ok {
		return errors.New("register module activated but storer could not be upgraded to CreatingServerStorer")
	}

	if err := ab.Config.Core.ViewRenderer.Load(PageRegister); err != nil {
		return err
	}

	sort.Strings(ab.Config.Modules.RegisterPreserveFields)

	ab.Config.Core.Router.GET("/register", ab.Config.Core.ErrorHandler.Wrap(r.Get))
	ab.Config.Core.Router.POST("/register", ab.Config.Core.ErrorHandler.Wrap(r.Post))

	return nil
}

// Get the register page
func (r *Register) Get(w http.ResponseWriter, req *http.Request) error {
	return r.Config.Core.Responder.Respond(w, req, http.StatusOK, PageRegister, nil)
}

// Post to the register page
func (r *Register) Post(w http.ResponseWriter, req *http.Request) error {
	logger := r.RequestLogger(req)
	validatable, err := r.Core.BodyReader.Read(PageRegister, req)
	if err != nil {
		return err
	}

	var arbitrary map[string]string
	var preserve map[string]string
	if arb, ok := validatable.(engine.ArbitraryValuer); ok {
		arbitrary = arb.GetValues()
		preserve = make(map[string]string)

		for k, v := range arbitrary {
			if hasString(r.Config.Modules.RegisterPreserveFields, k) {
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
		return r.Config.Core.Responder.Respond(w, req, http.StatusOK, PageRegister, data)
	}

	// Get values from request
	userVals := engine.MustHaveUserValues(validatable)
	pid, password := userVals.GetPID(), userVals.GetPassword()

	// Put values into newly created user for storage
	storer := engine.EnsureCanCreate(r.Config.Storage.Server)
	user := engine.MustBeAuthable(storer.New(req.Context()))

	pass, err := bcrypt.GenerateFromPassword([]byte(password), r.Config.Modules.BCryptCost)
	if err != nil {
		return err
	}

	user.PutPID(pid)
	user.PutPassword(string(pass))

	if arbUser, ok := user.(engine.ArbitraryUser); ok && arbitrary != nil {
		arbUser.PutArbitrary(arbitrary)
	}

	err = storer.Create(req.Context(), user)
	switch {
	case err == engine.ErrUserFound:
		logger.Infof("user %s attempted to re-register", pid)
		errs = []error{errors.New("user already exists")}
		data := engine.HTMLData{
			engine.DataValidation: engine.ErrorMap(errs),
		}
		if preserve != nil {
			data[engine.DataPreserve] = preserve
		}
		return r.Config.Core.Responder.Respond(w, req, http.StatusOK, PageRegister, data)
	case err != nil:
		return err
	}

	req = req.WithContext(context.WithValue(req.Context(), engine.CTXKeyUser, user))
	handled, err := r.Events.FireAfter(engine.EventRegister, w, req)
	if err != nil {
		return err
	} else if handled {
		return nil
	}

	// Log the user in, but only if the response wasn't handled previously
	// by a module like confirm.
	engine.PutSession(w, engine.SessionKey, pid)

	logger.Infof("registered and logged in user %s", pid)
	ro := engine.RedirectOptions{
		Code:         http.StatusTemporaryRedirect,
		Success:      "Account successfully created, you are now logged in",
		RedirectPath: "/login",
	}
	return r.Config.Core.Redirector.Redirect(w, req, ro)
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
