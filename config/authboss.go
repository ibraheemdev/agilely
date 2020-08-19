package config

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/users"
	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
	"github.com/ibraheemdev/agilely/pkg/authboss/authboss/defaults"
	"github.com/ibraheemdev/agilely/pkg/mailer"
	"github.com/ibraheemdev/agilely/pkg/renderer"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"golang.org/x/crypto/bcrypt"

	// Blank import all the modules for init functions
	_ "github.com/ibraheemdev/agilely/pkg/authboss/authenticatable"
	_ "github.com/ibraheemdev/agilely/pkg/authboss/confirmable"
	_ "github.com/ibraheemdev/agilely/pkg/authboss/lockable"
	_ "github.com/ibraheemdev/agilely/pkg/authboss/logoutable"
	_ "github.com/ibraheemdev/agilely/pkg/authboss/oauthable"
	_ "github.com/ibraheemdev/agilely/pkg/authboss/recoverable"
	_ "github.com/ibraheemdev/agilely/pkg/authboss/registerable"
	"github.com/ibraheemdev/agilely/pkg/authboss/rememberable"
)

// SetupAuthboss :
func SetupAuthboss(r *httprouter.Router) {
	ab := authboss.New()

	rt := defaults.NewRouter(r)
	rt.Use(alice.New(ab.LoadClientStateMiddleware, rememberable.Middleware(ab)))
	ab.Config.Core.Router = rt

	ab.Config.Core.ErrorHandler = defaults.NewErrorHandler(defaults.NewLogger(os.Stdout))
	ab.Config.Core.ViewRenderer = renderer.NewHTMLRenderer("/", "web/templates/authboss/*.tpl", "web/templates/layouts/*.tpl")
	ab.Config.Core.MailRenderer = renderer.NewHTMLRenderer("/", "web/templates/authboss/mailer/*.tpl", "web/templates/layouts/mailer/*.tpl")
	ab.Config.Core.Responder = defaults.NewResponder(ab.Config.Core.ViewRenderer)
	ab.Config.Core.Redirector = defaults.NewRedirector(ab.Config.Core.ViewRenderer, authboss.FormValueRedirect)
	ab.Config.Core.BodyReader = defaults.NewHTTPBodyReader(false, false)
	ab.Config.Core.Mailer = mailer.NewLogMailer(os.Stdout)
	ab.Config.Core.Logger = defaults.NewLogger(os.Stdout)

	ab.Config.Storage.Server = users.DB
	ab.Config.Storage.SessionState = defaults.NewSessionStorer("agilely_session", []byte("TODO"), nil)
	ab.Config.Storage.CookieState = defaults.NewCookieStorer([]byte("TODO"), nil)
	ab.Config.Storage.SessionStateWhitelistKeys = []string{}

	ab.Config.Paths.Mount = "/"
	ab.Config.Paths.NotAuthorized = "/"
	ab.Config.Paths.AuthLoginOK = "/"
	ab.Config.Paths.ConfirmOK = "/"
	ab.Config.Paths.ConfirmNotOK = "/auth/login"
	ab.Config.Paths.LockNotOK = "/auth/login"
	ab.Config.Paths.LogoutOK = "/"
	ab.Config.Paths.LogoutOK = "/"
	ab.Config.Paths.OAuth2LoginOK = "/"
	ab.Config.Paths.OAuth2LoginNotOK = "/"
	ab.Config.Paths.OAuth2LoginNotOK = "/"
	ab.Config.Paths.RecoverOK = "/"
	ab.Config.Paths.RegisterOK = "/"
	ab.Config.Paths.RootURL = "http://localhost:8080"
	ab.Config.Paths.TwoFactorEmailAuthNotOK = "/"

	ab.Config.Modules.BCryptCost = bcrypt.DefaultCost
	ab.Config.Modules.ExpireAfter = time.Hour
	ab.Config.Modules.LockAfter = 3
	ab.Config.Modules.LockWindow = 5 * time.Minute
	ab.Config.Modules.LockDuration = 12 * time.Hour
	ab.Config.Modules.LogoutMethod = "DELETE"
	ab.Config.Modules.MailRouteMethod = http.MethodGet
	ab.Config.Modules.MailNoGoroutine = false
	ab.Config.Modules.RegisterPreserveFields = []string{"email"}
	ab.Config.Modules.RecoverTokenDuration = 24 * time.Hour
	ab.Config.Modules.RecoverLoginAfterRecovery = false
	ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{}
	ab.Config.Modules.TwoFactorEmailAuthRequired = true

	ab.Config.Mail.RootURL = ""
	ab.Config.Mail.From = "agilely@agilely.com"
	ab.Config.Mail.FromName = "agilely"
	ab.Config.Mail.SubjectPrefix = ""

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}
