package config

import (
	"os"

	"github.com/justinas/alice"

	"github.com/ibraheemdev/agilely/internal/app/engine"
	"github.com/ibraheemdev/agilely/pkg/body_reader"
	"github.com/ibraheemdev/agilely/pkg/client_state"
	"github.com/ibraheemdev/agilely/pkg/error_handler"
	"github.com/ibraheemdev/agilely/pkg/logger"
	"github.com/ibraheemdev/agilely/pkg/mailer"
	"github.com/ibraheemdev/agilely/pkg/renderer"
	"github.com/ibraheemdev/agilely/pkg/responder"
	"github.com/ibraheemdev/agilely/pkg/router"
	"github.com/julienschmidt/httprouter"
)

// SetCore : Set's the core components of an engine instance
func SetCore(e *engine.Engine) {
	rt := router.NewRouter(httprouter.New())
	rt.Use(alice.New(e.LoadClientStateMiddleware))
	e.Core.Router = rt

	e.Core.MailRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/mailer/*.tpl", "web/templates/layouts/mailer/*.tpl")
	e.Core.ViewRenderer = renderer.NewHTMLRenderer("/", "web/templates/users/*.tpl", "web/templates/layouts/*.tpl")
	e.Core.SessionState = clientstate.NewSessionStorer("agilely_session", []byte("TODO"), nil)
	e.Core.Redirector = responder.NewRedirector(e.Core.ViewRenderer, engine.FormValueRedirect)
	e.Core.CookieState = clientstate.NewCookieStorer([]byte("TODO"), nil)
	e.Core.ErrorHandler = errorhandler.New(logger.New(os.Stdout))
	e.Core.Responder = responder.New(e.Core.ViewRenderer)
	e.Core.BodyReader = bodyreader.NewHTTP(false, false)
	e.Core.Mailer = mailer.NewLogMailer(os.Stdout)
	e.Core.Logger = logger.New(os.Stdout)
}

// SetConfig : Set's the engine's configuration variables by reading config files
func SetConfig(e *engine.Engine) error {
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}
	e.Config = *cfg
	return nil
}
