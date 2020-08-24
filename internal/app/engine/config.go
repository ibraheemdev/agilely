package engine

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Config holds all the configuration for both engine and it's modules.
type Config struct {
	Paths struct {
		// Mount is the path to mount engine's routes at (eg /auth).
		Mount string

		// RootURL is the scheme+host+port of the web application
		// (eg https://www.happiness.com:8080) for url generation.
		// No trailing slash.
		RootURL string
	}

	Modules struct {
		// BCryptCost is the cost of the bcrypt password hashing function.
		BCryptCost int

		// ExpireAfter controls the time an account is idle before being
		// logged out by the ExpireMiddleware.
		ExpireAfter time.Duration

		// LockAfter this many tries.
		LockAfter int
		// LockWindow is the waiting time before the number of attempts are reset.
		LockWindow time.Duration
		// LockDuration is how long an account is locked for.
		LockDuration time.Duration

		// RegisterPreserveFields are fields used with registration that are
		// to be rendered when post fails in a normal way
		// (for example validation errors), they will be passed back in the
		// data of the response under the key DataPreserve which
		// will be a map[string]string. This way the user does not have to
		// retype these whitelisted fields.
		//
		// All fields that are to be preserved must be able to be returned by
		// the ArbitraryValuer.GetValues()
		//
		// This means in order to have a field named "address" you would need
		// to have that returned by the ArbitraryValuer.GetValues() method and
		// then it would be available to be whitelisted by this
		// configuration variable.
		RegisterPreserveFields []string

		// RecoverTokenDuration controls how long a token sent via
		// email for password recovery is valid for.
		RecoverTokenDuration time.Duration

		// OAuth2Providers lists all providers that can be used. See
		// OAuthProvider documentation for more details.
		OAuth2Providers map[string]OAuth2Provider
	}

	Mail struct {
		// RootURL is a full path to an application that is hosting a front-end
		// Typically using a combination of Paths.RootURL and Paths.Mount
		// MailRoot will be assembled if not set.
		// Typically looks like: https://our-front-end.com/authenication
		// No trailing slash.
		RootURL string

		// From is the email address engine e-mails come from.
		From string
		// FromName is the name engine e-mails come from.
		FromName string
		// SubjectPrefix is used to add something to the front of the engine
		// email subjects.
		SubjectPrefix string
	}

	Storage struct {
		// Storer is the interface through which Engine accesses the web apps
		// database for user operations.
		Server ServerStorer

		// CookieState must be defined to provide an interface capapable of
		// storing cookies for the given response, and reading them from the
		// request.
		CookieState ClientStateReadWriter
		// SessionState must be defined to provide an interface capable of
		// storing session-only values for the given response, and reading them
		// from the request.
		SessionState ClientStateReadWriter

		// SessionStateWhitelistKeys are set to preserve keys in the session
		// when engine.DelAllSession is called. A correct implementation
		// of ClientStateReadWriter will delete ALL session key-value pairs
		// unless that key is whitelisted here.
		SessionStateWhitelistKeys []string
	}

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
	}
}

// Defaults sets the configuration's default values.
func (c *Config) Defaults() {
	c.Paths.Mount = "/auth"
	c.Paths.RootURL = "http://localhost:8080"

	c.Modules.BCryptCost = bcrypt.DefaultCost
	c.Modules.ExpireAfter = time.Hour
	c.Modules.LockAfter = 3
	c.Modules.LockWindow = 5 * time.Minute
	c.Modules.LockDuration = 12 * time.Hour
	c.Modules.RecoverTokenDuration = 24 * time.Hour
	c.Modules.RegisterPreserveFields = []string{"email"}
}