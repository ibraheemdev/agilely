package config

import (
	"log"
	"os"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/internal/app/engine/defaults"
	"github.com/ibraheemdev/agilely/internal/app/users"
	"github.com/ibraheemdev/agilely/pkg/client_state"
	"github.com/ibraheemdev/agilely/pkg/logger"
	"github.com/ibraheemdev/agilely/pkg/mailer"
	"github.com/ibraheemdev/agilely/pkg/renderer"
	"github.com/ibraheemdev/agilely/pkg/responder"
	"github.com/ibraheemdev/agilely/pkg/router"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"golang.org/x/crypto/bcrypt"
)

// SetupEngine :
func SetupEngine(r *httprouter.Router) {
	e := engine.New()

	rt := router.NewRouter(r)
	rt.Use(alice.New(e.LoadClientStateMiddleware, users.RememberMiddleware(e)))
	e.Config.Core.Router = rt

	e.Config.Core.ErrorHandler = defaults.NewErrorHandler(logger.New(os.Stdout))
	e.Config.Core.ViewRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/*.tpl", "web/templates/layouts/*.tpl")
	e.Config.Core.MailRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/mailer/*.tpl", "web/templates/layouts/mailer/*.tpl")
	e.Config.Core.Responder = responder.New(e.Config.Core.ViewRenderer)
	e.Config.Core.Redirector = responder.NewRedirector(e.Config.Core.ViewRenderer, engine.FormValueRedirect)
	e.Config.Core.BodyReader = defaults.NewHTTPBodyReader(false, false)
	e.Config.Core.Mailer = mailer.NewLogMailer(os.Stdout)
	e.Config.Core.Logger = logger.New(os.Stdout)

	e.Config.Storage.Server = users.DB
	e.Config.Storage.SessionState = clientstate.NewSessionStorer("agilely_session", []byte("TODO"), nil)
	e.Config.Storage.CookieState = clientstate.NewCookieStorer([]byte("TODO"), nil)
	e.Config.Storage.SessionStateWhitelistKeys = []string{}

	e.Config.Paths.Mount = "/"
	e.Config.Paths.RootURL = "http://localhost:8080"

	e.Config.Modules.BCryptCost = bcrypt.DefaultCost
	e.Config.Modules.ExpireAfter = time.Hour
	e.Config.Modules.LockAfter = 3
	e.Config.Modules.LockWindow = 5 * time.Minute
	e.Config.Modules.LockDuration = 12 * time.Hour
	e.Config.Modules.RegisterPreserveFields = []string{"email"}
	e.Config.Modules.RecoverTokenDuration = 24 * time.Hour
	e.Config.Modules.OAuth2Providers = map[string]engine.OAuth2Provider{}

	e.Config.Mail.RootURL = ""
	e.Config.Mail.From = "agilely@agilely.com"
	e.Config.Mail.FromName = "agilely"
	e.Config.Mail.SubjectPrefix = ""

	if err := e.Init(); err != nil {
		log.Fatal(err)
	}
}
