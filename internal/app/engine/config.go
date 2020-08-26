package engine

import (
	"time"
)

// Config holds all the configuration for both engine and it's modules.
type Config struct {
	// Mount is the path to mount engine's routes at
	Mount string `yaml:"mount"`

	// No trailing slash.
	RootURL string `yaml:"root_url"`

	// SessionStateWhitelistKeys are set to preserve keys in the session
	// when engine.DelAllSession is called. A correct implementation
	// of ClientStateReadWriter will delete ALL session key-value pairs
	// unless that key is whitelisted here.
	SessionStateWhitelistKeys []string `yaml:"session_state_whitelist_keys"`

	Database struct {
		// database host
		Host string `yaml:"host"`

		// database port
		Port int `yaml:"port"`

		// database name
		Name string `yaml:"name"`
	} `yaml:"database"`

	Server struct {
		// server host
		Host string `yaml:"host"`

		// server port
		Port int `yaml:"port"`

		// config for static file server
		Static struct {
			// path to webpack manifest.json
			ManifestPath string `yaml:"manifestpath"`
		} `yaml:"static"`
	} `yaml:"server"`

	Authboss struct {
		// BCryptCost is the cost of the bcrypt password hashing function.
		BCryptCost int `yaml:"bcrypt_cost"`

		// ExpireAfter controls the time an account is idle before being
		// logged out by the ExpireMiddleware.
		ExpireAfter time.Duration `yaml:"expire_after"`

		// LockAfter this many tries.
		LockAfter int `yaml:"lock_after"`

		// LockWindow is the waiting time before the number of attempts are reset.
		LockWindow time.Duration `yaml:"lock_window"`

		// LockDuration is how long an account is locked for.
		LockDuration time.Duration `yaml:"lock_duration"`

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
		RegisterPreserveFields []string `yaml:"register_preserve_fields"`

		// RecoverTokenDuration controls how long a token sent via
		// email for password recovery is valid for.
		RecoverTokenDuration time.Duration `yaml:"recover_token_duration"`

		// OAuth2Providers lists all providers that can be used. See
		// OAuthProvider documentation for more details.
		OAuth2Providers map[string]OAuth2Provider
	} `yaml:"authboss"`

	Mail struct {
		// No trailing slash.
		RootURL string `yaml:"root_url"`

		// From is the email address engine e-mails come from.
		From string `yaml:"from"`

		// FromName is the name engine e-mails come from.
		FromName string `yaml:"from_name"`

		// SubjectPrefix is used to add something to the front of the engine
		// email subjects.
		SubjectPrefix string `yaml:"subject_prefix"`
	} `yaml:"mail"`
}
