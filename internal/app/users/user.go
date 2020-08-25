package users

import (
	"context"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"
)

// InMemDB : A generic in memory database
type InMemDB struct {
	Users map[string]User
}

// User database model
type User struct {
	ID   string
	Name string

	// Authable
	Email    string
	Password string

	// Recoverable
	RecoverSelector    string
	RecoverVerifier    string
	RecoverTokenExpiry time.Time

	// Confirmable
	ConfirmSelector string
	ConfirmVerifier string
	Confirmed       bool

	// Lockable
	AttemptCount int
	LastAttempt  time.Time
	Locked       time.Time

	// OAuthable
	OAuth2UID          string
	OAuth2Provider     string
	OAuth2AccessToken  string
	OAuth2RefreshToken string
	OAuth2Expiry       time.Time

	// Rememberable
	RememberTokens []string
}

// DB : The in memory database instance
var DB *InMemDB

// The init function ensures that the user model
// satisfies the interface for the included modules
func init() {
	DB = &InMemDB{Users: map[string]User{
		"2a9a-bcb8-4901": {
			ID:        "2a9a-bcb8-4901",
			Name:      "John",
			Password:  "Hashed-Password",
			Email:     "John@example.org",
			Confirmed: true,
		},
	},
	}
	// This rest of the init function is optional,
	// albiet a good practice
	assertUser := new(User)

	var _ engine.User = assertUser
	var _ engine.AuthableUser = assertUser
	var _ engine.ConfirmableUser = assertUser
	var _ engine.LockableUser = assertUser
	var _ engine.RecoverableUser = assertUser
	var _ engine.ArbitraryUser = assertUser
	var _ engine.OAuth2User = assertUser

	var _ engine.CreatingServerStorer = DB
	var _ engine.ConfirmingServerStorer = DB
	var _ engine.RecoveringServerStorer = DB
	var _ engine.RememberingServerStorer = DB
}

// ************** Engine User **************

// GetPID from user
func (u User) GetPID() string {
	return u.ID
}

// PutPID into user
func (u *User) PutPID(pid string) {
	u.ID = pid
}

// ************** Authable **************

// GetPassword from user
func (u User) GetPassword() string {
	return u.Password
}

// PutPassword into user
func (u *User) PutPassword(password string) {
	u.Password = password
}

// ************** Confirmable **************

// GetConfirmSelector from user
func (u User) GetConfirmSelector() string {
	return u.ConfirmSelector
}

// PutConfirmSelector into user
func (u *User) PutConfirmSelector(confirmSelector string) {
	u.ConfirmSelector = confirmSelector
}

// GetConfirmVerifier from user
func (u User) GetConfirmVerifier() string {
	return u.ConfirmVerifier
}

// PutConfirmVerifier into user
func (u *User) PutConfirmVerifier(confirmVerifier string) {
	u.ConfirmVerifier = confirmVerifier
}

// GetConfirmed from user
func (u User) GetConfirmed() bool {
	return u.Confirmed
}

// PutConfirmed into user
func (u *User) PutConfirmed(confirmed bool) {
	u.Confirmed = confirmed
}

// ************** Lockable **************

// GetLastAttempt from user
func (u User) GetLastAttempt() time.Time {
	return u.LastAttempt
}

// PutLastAttempt into user
func (u *User) PutLastAttempt(last time.Time) {
	u.LastAttempt = last
}

// GetLocked from user
func (u User) GetLocked() time.Time {
	return u.Locked
}

// PutLocked into user
func (u *User) PutLocked(locked time.Time) {
	u.Locked = locked
}

// GetAttemptCount from user
func (u User) GetAttemptCount() int {
	return u.AttemptCount
}

// PutAttemptCount into user
func (u *User) PutAttemptCount(attempts int) {
	u.AttemptCount = attempts
}

// ************** Recoverable **************

// GetEmail from user
func (u User) GetEmail() string {
	return u.Email
}

// PutEmail into user
func (u *User) PutEmail(email string) {
	u.Email = email
}

// GetRecoverVerifier from user
func (u User) GetRecoverVerifier() string {
	return u.RecoverVerifier
}

// PutRecoverVerifier into user
func (u *User) PutRecoverVerifier(token string) {
	u.RecoverVerifier = token
}

// GetRecoverExpiry from user
func (u User) GetRecoverExpiry() time.Time {
	return u.RecoverTokenExpiry
}

// PutRecoverExpiry into user
func (u *User) PutRecoverExpiry(expiry time.Time) {
	u.RecoverTokenExpiry = expiry
}

// GetRecoverSelector from user
func (u User) GetRecoverSelector() string {
	return u.RecoverSelector
}

// PutRecoverSelector into user
func (u *User) PutRecoverSelector(token string) {
	u.RecoverSelector = token
}

// ************** Arbitrary **************

// GetArbitrary data from user
func (u User) GetArbitrary() map[string]string {
	return map[string]string{
		"name": u.Name,
		// ...
	}
}

// PutArbitrary data from user
func (u *User) PutArbitrary(values map[string]string) {
	if n, ok := values["name"]; ok {
		u.Name = n
		// ...
	}
}

// ************** OAuthable **************

// IsOAuth2User returns true if the user was created with oauth2
func (u User) IsOAuth2User() bool {
	return len(u.OAuth2UID) != 0
}

// GetOAuth2UID from user
func (u User) GetOAuth2UID() (uid string) {
	return u.OAuth2UID
}

// PutOAuth2UID into user
func (u *User) PutOAuth2UID(uid string) {
	u.OAuth2UID = uid
}

// GetOAuth2Provider from user
func (u User) GetOAuth2Provider() (provider string) {
	return u.OAuth2Provider
}

// PutOAuth2Provider into user
func (u *User) PutOAuth2Provider(provider string) {
	u.OAuth2Provider = provider
}

// GetOAuth2AccessToken from user
func (u User) GetOAuth2AccessToken() (token string) {
	return u.OAuth2AccessToken
}

// PutOAuth2AccessToken into user
func (u *User) PutOAuth2AccessToken(token string) {
	u.OAuth2AccessToken = token
}

// GetOAuth2RefreshToken from user
func (u User) GetOAuth2RefreshToken() (refreshToken string) {
	return u.OAuth2RefreshToken
}

// PutOAuth2RefreshToken into user
func (u *User) PutOAuth2RefreshToken(refreshToken string) {
	u.OAuth2RefreshToken = refreshToken
}

// GetOAuth2Expiry from user
func (u User) GetOAuth2Expiry() (expiry time.Time) {
	return u.OAuth2Expiry
}

// PutOAuth2Expiry into user
func (u *User) PutOAuth2Expiry(expiry time.Time) {
	u.OAuth2Expiry = expiry
}

// ************** CreatingServerStorer **************

// Create and save the user
func (db *InMemDB) Create(_ context.Context, u engine.User) error {
	user := u.(*User)
	db.Users[user.ID] = *user
	return nil
}

// New user creation; the user not saved
func (db *InMemDB) New(_ context.Context) engine.User {
	return &User{}
}

// Save the user
func (db *InMemDB) Save(_ context.Context, u engine.User) error {
	user := u.(*User)
	db.Users[user.ID] = *user
	return nil
}

// Load the user
func (db *InMemDB) Load(_ context.Context, id string) (engine.User, error) {
	// Check to see if our key is actually an oauth2 id
	provider, uid, err := engine.ParseOAuth2PID(id)
	if err == nil {
		for _, u := range db.Users {
			if u.OAuth2Provider == provider && u.OAuth2UID == uid {
				return &u, nil
			}
		}
		return nil, engine.ErrUserNotFound
	}

	u, ok := db.Users[id]
	if !ok {
		return nil, engine.ErrUserNotFound
	}
	return &u, nil
}

// ************** ConfirmingServerStorer **************

// LoadByConfirmSelector looks a user up by confirmation token
func (db *InMemDB) LoadByConfirmSelector(_ context.Context, selector string) (user engine.ConfirmableUser, err error) {
	for _, v := range db.Users {
		if v.ConfirmSelector == selector {
			return &v, nil
		}
	}

	return nil, engine.ErrUserNotFound
}

// ************** RecoveringServerStorer **************

// LoadByRecoverSelector looks a user up by confirmation selector
func (db *InMemDB) LoadByRecoverSelector(_ context.Context, selector string) (user engine.RecoverableUser, err error) {
	for _, v := range db.Users {
		if v.RecoverSelector == selector {
			return &v, nil
		}
	}

	return nil, engine.ErrUserNotFound
}

// ************** RememberingServerStorer **************

// AddRememberToken to a user
func (db *InMemDB) AddRememberToken(_ context.Context, pid, t string) error {
	user := db.Users[pid]
	user.RememberTokens = append(user.RememberTokens, t)
	return nil
}

// DelRememberTokens removes all tokens for the given pid
func (db *InMemDB) DelRememberTokens(_ context.Context, pid string) error {
	user := db.Users[pid]
	user.RememberTokens = nil
	return nil
}

// UseRememberToken finds the pid-token pair and deletes it.
// ie: the user uses the remember token, so it is now invalid
// If the token could not be found return ErrTokenNotFound
func (db *InMemDB) UseRememberToken(_ context.Context, pid, token string) error {
	user := db.Users[pid]
	tokens := user.RememberTokens
	for i, t := range tokens {
		if t == token {
			tokens[i] = tokens[len(tokens)-1]
			user.RememberTokens = tokens[:len(tokens)-1]
			return nil
		}
	}
	return engine.ErrTokenNotFound
}
