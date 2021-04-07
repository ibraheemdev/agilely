package engine

// Engine contains a configuration and other details for running.
type Engine struct {
	// Engine configuration
	Config Config

	// Authentication Events, used as controller level middleware
	AuthEvents *AuthEvents

	Core struct {
		// Router is the entity that controls all routing to engine routes
		// modules will register their routes with it.
		Router Router

		// ErrorHandler wraps http requests with centralized error handling.
		ErrorHandler ErrorHandler

		// Responder takes a generic response from a controller and prepares
		// the response, uses a renderer to create the body, and replies to the
		// http request.
		Responder HTTPResponder

		// Redirector can redirect a response, similar to Responder but
		// responsible only for redirection.
		Redirector HTTPRedirector

		// BodyReader reads validatable data from the body of a request to
		// be able to get data from the user's client.
		BodyReader BodyReader

		// ViewRenderer loads the templates for the application.
		ViewRenderer Renderer

		// MailRenderer loads the templates for mail
		MailRenderer Renderer

		// Mailer is the mailer being used to send e-mails out via smtp
		Mailer Mailer

		// Logger implies just a few log levels for use, can optionally
		// also implement the ContextLogger to be able to upgrade to a
		// request specific logger.
		Logger Logger

		// Database
		Database Database

		// CookieState must be defined to provide an interface capapable of
		// storing cookies for the given response, and reading them from the
		// request.
		CookieState ClientStateReadWriter

		// SessionState must be defined to provide an interface capable of
		// storing session-only values for the given response, and reading them
		// from the request.
		SessionState ClientStateReadWriter
	}
}

// New makes a new instance of engine with a default
// configuration.
func New() *Engine {
	e := &Engine{}
	e.AuthEvents = NewAuthEvents()
	e.Config.Authboss.OAuth2Providers = map[string]OAuth2Provider{}

	return e
}
