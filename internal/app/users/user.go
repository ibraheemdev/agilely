package users

import (
	"context"
	"time"

	"github.com/ibraheemdev/agilely/internal/app/engine"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	// Collection ...
	Collection = "users"
)

// User database model
type User struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`

	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`

	RecoverSelector    string    `json:"recover_selector" bson:"recover_selector"`
	RecoverVerifier    string    `json:"recover_verifier" bson:"recover_verifier"`
	RecoverTokenExpiry time.Time `json:"recover_token_expiry" bson:"recover_token_expiry"`

	ConfirmSelector string `json:"confirm_selector" bson:"confirm_selector"`
	ConfirmVerifier string `json:"confirm_verifier" bson:"confirm_verifier"`
	Confirmed       bool   `json:"confirmed" bson:"confirmed"`

	AttemptCount int       `json:"attempt_count" bson:"attempt_count"`
	LastAttempt  time.Time `json:"last_attempt" bson:"last_attempt"`
	Locked       time.Time `json:"locked" bson:"locked"`

	OAuthUID          string    `json:"oauth_uid" bson:"oauth_uid"`
	OAuthProvider     string    `json:"oauth_provider" bson:"oauth_provider"`
	OAuthAccessToken  string    `json:"oauth_access_token" bson:"oauth_access_token"`
	OAuthRefreshToken string    `json:"oauth_refresh_token" bson:"oauth_refresh_token"`
	OAuthExpiry       time.Time `json:"oauth_expiry" bson:"oauth_expiry"`

	RememberTokens []string `json:"remember_tokens" bson:"remember_tokens"`
}

// GetUser ...
func GetUser(ctx context.Context, d engine.Database, key string) (*User, error) {
	// Check to see if our key is actually an oauth2 id
	provider, uid, err := engine.ParseOAuth2PID(key)
	var u User
	if err == nil {
		err := d.Collection(Collection).FindOne(ctx, bson.D{{"oauth_provider", provider}, {"oauth_uid", uid}}).Decode(&u)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}

	err = d.Collection(Collection).FindOne(ctx, bson.D{{"_email", key}}).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UseRememberToken ...
func (u *User) UseRememberToken(token string) ([]string, error) {
	tokens := u.RememberTokens
	for i, t := range tokens {
		if t == token {
			tokens[i] = tokens[len(tokens)-1]
			return tokens[:len(tokens)-1], nil
		}
	}
	// the token did not exist
	return nil, engine.ErrNoDocuments
}
