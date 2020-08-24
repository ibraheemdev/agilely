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
	ab := engine.New()

	rt := router.NewRouter(r)
	rt.Use(alice.New(ab.LoadClientStateMiddleware, users.RememberMiddleware(ab)))
	ab.Config.Core.Router = rt

	ab.Config.Core.ErrorHandler = defaults.NewErrorHandler(logger.New(os.Stdout))
	ab.Config.Core.ViewRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/*.tpl", "web/templates/layouts/*.tpl")
	ab.Config.Core.MailRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/mailer/*.tpl", "web/templates/layouts/mailer/*.tpl")
	ab.Config.Core.Responder = responder.New(ab.Config.Core.ViewRenderer)
	ab.Config.Core.Redirector = responder.NewRedirector(ab.Config.Core.ViewRenderer, engine.FormValueRedirect)
	ab.Config.Core.BodyReader = defaults.NewHTTPBodyReader(false, false)
	ab.Config.Core.Mailer = mailer.NewLogMailer(os.Stdout)
	ab.Config.Core.Logger = logger.New(os.Stdout)

	ab.Config.Storage.Server = users.DB
	ab.Config.Storage.SessionState = clientstate.NewSessionStorer("agilely_session", []byte("TODO"), nil)
	ab.Config.Storage.CookieState = clientstate.NewCookieStorer([]byte("TODO"), nil)
	ab.Config.Storage.SessionStateWhitelistKeys = []string{}

	ab.Config.Paths.Mount = "/"
	ab.Config.Paths.RootURL = "http://localhost:8080"

	ab.Config.Modules.BCryptCost = bcrypt.DefaultCost
	ab.Config.Modules.ExpireAfter = time.Hour
	ab.Config.Modules.LockAfter = 3
	ab.Config.Modules.LockWindow = 5 * time.Minute
	ab.Config.Modules.LockDuration = 12 * time.Hour
	ab.Config.Modules.RegisterPreserveFields = []string{"email"}
	ab.Config.Modules.RecoverTokenDuration = 24 * time.Hour
	ab.Config.Modules.OAuth2Providers = map[string]engine.OAuth2Provider{}

	ab.Config.Mail.RootURL = ""
	ab.Config.Mail.From = "agilely@agilely.com"
	ab.Config.Mail.FromName = "agilely"
	ab.Config.Mail.SubjectPrefix = ""

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}
