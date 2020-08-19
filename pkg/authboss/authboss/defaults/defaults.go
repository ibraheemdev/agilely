// Package defaults houses default implementations for the very many
// interfaces that authboss has. It's a goal of the defaults package
// to provide the core where authboss implements the shell.
//
// It's simultaneously supposed to be possible to take as many or
// as few of these implementations as you desire, allowing you to
// reimplement where necessary, but reuse where possible.
package defaults

import (
	"os"

	"github.com/ibraheemdev/agilely/pkg/authboss/authboss"
	"github.com/ibraheemdev/agilely/pkg/mailer"
	"github.com/ibraheemdev/agilely/pkg/renderer"
	"github.com/julienschmidt/httprouter"
)

// SetCore creates instances of all the default pieces
// with the exception of ViewRenderer which should be already set
// before calling this method.
func SetCore(config *authboss.Config, readJSON, useUsername bool, mountPath, templatesPath, layoutsDir string) {
	logger := NewLogger(os.Stdout)

	config.Core.Router = NewRouter(httprouter.New())
	config.Core.ErrorHandler = NewErrorHandler(logger)
	config.Core.ViewRenderer = renderer.NewHTMLRenderer(mountPath, templatesPath, layoutsDir)
	config.Core.MailRenderer = renderer.NewHTMLRenderer(mountPath, templatesPath, layoutsDir)
	config.Core.Responder = NewResponder(config.Core.ViewRenderer)
	config.Core.Redirector = NewRedirector(config.Core.ViewRenderer, authboss.FormValueRedirect)
	config.Core.BodyReader = NewHTTPBodyReader(readJSON, useUsername)
	config.Core.Mailer = mailer.NewLogMailer(os.Stdout)
	config.Core.Logger = logger
}
