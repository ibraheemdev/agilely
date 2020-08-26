package config

import (
	"os"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/internal/app/users"
	"github.com/ibraheemdev/agilely/pkg/body_reader"
	"github.com/ibraheemdev/agilely/pkg/client_state"
	"github.com/ibraheemdev/agilely/pkg/error_handler"
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
	e.Core.Router = rt

	e.Core.ErrorHandler = errorhandler.New(logger.New(os.Stdout))
	e.Core.ViewRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/*.tpl", "web/templates/layouts/*.tpl")
	e.Core.MailRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/mailer/*.tpl", "web/templates/layouts/mailer/*.tpl")
	e.Core.Responder = responder.New(e.Core.ViewRenderer)
	e.Core.Redirector = responder.NewRedirector(e.Core.ViewRenderer, engine.FormValueRedirect)
	e.Core.BodyReader = bodyreader.NewHTTP(false, false)
	e.Core.Mailer = mailer.NewLogMailer(os.Stdout)
	e.Core.Logger = logger.New(os.Stdout)

	e.Core.Server = users.DB
	e.Core.SessionState = clientstate.NewSessionStorer("agilely_session", []byte("TODO"), nil)
	e.Core.CookieState = clientstate.NewCookieStorer([]byte("TODO"), nil)
	e.Config.SessionStateWhitelistKeys = []string{}

	e.Config.Mount = "/"
	e.Config.RootURL = "http://localhost:8080"

	e.Config.Authboss.BCryptCost = bcrypt.DefaultCost
	e.Config.Authboss.ExpireAfter = time.Hour
	e.Config.Authboss.LockAfter = 3
	e.Config.Authboss.LockWindow = 5 * time.Minute
	e.Config.Authboss.LockDuration = 12 * time.Hour
	e.Config.Authboss.RegisterPreserveFields = []string{"email"}
	e.Config.Authboss.RecoverTokenDuration = 24 * time.Hour
	e.Config.Authboss.OAuth2Providers = map[string]engine.OAuth2Provider{}

	e.Config.Mail.RootURL = ""
	e.Config.Mail.From = "agilely@agilely.com"
	e.Config.Mail.FromName = "agilely"
	e.Config.Mail.SubjectPrefix = ""
}
