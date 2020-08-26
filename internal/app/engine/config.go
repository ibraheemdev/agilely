package engine

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Config holds all the configuration for both engine and it's modules.
type Config struct {

	// Mount is the path to mount engine's routes at (eg /auth).
	Mount string

	// RootURL is the scheme+host+port of the web application
	// (eg https://www.happiness.com:8080) for url generation.
	// No trailing slash.
	RootURL string

	// SessionStateWhitelistKeys are set to preserve keys in the session
	// when engine.DelAllSession is called. A correct implementation
	// of ClientStateReadWriter will delete ALL session key-value pairs
	// unless that key is whitelisted here.
	SessionStateWhitelistKeys []string

	Authboss struct {
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
		// Typically using a combination of .RootURL and .Mount
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
}

// Defaults sets the configuration's default values.
func (c *Config) Defaults() {
	c.Mount = "/auth"
	c.RootURL = "http://localhost:8080"

	c.Authboss.BCryptCost = bcrypt.DefaultCost
	c.Authboss.ExpireAfter = time.Hour
	c.Authboss.LockAfter = 3
	c.Authboss.LockWindow = 5 * time.Minute
	c.Authboss.LockDuration = 12 * time.Hour
	c.Authboss.RecoverTokenDuration = 24 * time.Hour
	c.Authboss.RegisterPreserveFields = []string{"email"}
}
