package defaults

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/ibraheemdev/poller/pkg/authboss/authboss"
)

// CookieState is an authboss.ClientState implementation to hold
// cookie state for the duration of the request
type CookieState map[string]string

// Get a cookie's value
func (c CookieState) Get(key string) (string, bool) {
	cookie, ok := c[key]
	return cookie, ok
}

// CookieStorer writes and reads cookies to an underlying
// gorilla secure cookie storage.
//
// Because it embeds the SecureCookie piece this can be used
// as the cookie storage for your entire application (rather than
// only as a stub for authboss).
type CookieStorer struct {
	Cookies []string
	*securecookie.SecureCookie

	// Defaults empty
	Domain string
	// Defaults to /
	Path string
	// Defaults to 1 month
	MaxAge int
	// Defaults to true
	HTTPOnly bool
	// Defaults to true
	Secure bool
	// Samesite defaults to 0 or "off"
	SameSite http.SameSite
}

// NewCookieStorer constructor simply wraps the constructor for
// securecookie.New. The parameters are the hash key and the block key.
//
// hashKey is required, used to authenticate values using HMAC. Create
// it using GenerateRandomKey(). It is recommended to use a key with 32 or 64 bytes.
//
// blockKey is optional, used to encrypt values. Create it using GenerateRandomKey().
// The key length must correspond to the key size of the encryption algorithm.
// For AES, used by default, valid lengths are 16, 24, or 32 bytes to select AES-128,
// AES-192, or AES-256. The default encoder used for cookie serialization is
// encoding/gob.
//
// Note that keys created using GenerateRandomKey() are not automatically persisted.
// New keys will be created when the application is restarted, and previously issued
// cookies will not be able to be decoded.
func NewCookieStorer(hashKey, blockKey []byte) CookieStorer {
	return CookieStorer{
		Cookies:      []string{authboss.CookieRemember},
		SecureCookie: securecookie.New(hashKey, blockKey),
		Path:         "/",
		MaxAge:       int((time.Hour * 730) / time.Second), // 1 month
		HTTPOnly:     true,
		Secure:       true,
	}
}

// ReadState from the request
func (c CookieStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	cs := make(CookieState)

	for _, cookie := range r.Cookies() {
		for _, n := range c.Cookies {
			if n == cookie.Name {
				var str string
				if err := c.SecureCookie.Decode(n, cookie.Value, &str); err != nil {
					if e, ok := err.(securecookie.Error); ok {
						// Ignore bad cookies, this means that the client
						// may have bad cookies for a long time, but they should
						// eventually be overwritten by the application.
						if e.IsDecode() {
							continue
						}
					}
					return nil, err
				}

				cs[n] = str
			}
		}
	}

	return cs, nil
}

// WriteState to the responsewriter
func (c CookieStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, ev []authboss.ClientStateEvent) error {
	for _, ev := range ev {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			encoded, err := c.SecureCookie.Encode(ev.Key, ev.Value)
			if err != nil {
				return fmt.Errorf("failed to encode cookie: %w", err)
			}

			cookie := &http.Cookie{
				Expires: time.Now().UTC().AddDate(1, 0, 0),
				Name:    ev.Key,
				Value:   encoded,

				Domain:   c.Domain,
				Path:     c.Path,
				MaxAge:   c.MaxAge,
				HttpOnly: c.HTTPOnly,
				Secure:   c.Secure,
				SameSite: c.SameSite,
			}
			http.SetCookie(w, cookie)
		case authboss.ClientStateEventDel:
			cookie := &http.Cookie{
				MaxAge: -1,
				Name:   ev.Key,
				Domain: c.Domain,
				Path:   c.Path,
			}
			http.SetCookie(w, cookie)
		}
	}

	return nil
}
