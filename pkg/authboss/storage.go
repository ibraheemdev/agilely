package authboss

// A concious decision was made to put all storer
// and user types into this file despite them truly
// belonging to outside modules. The reason for this
// is because documentation-wise, it was previously
// difficult to find what you had to implement or even
// what you could implement.

import (
	"context"
	"errors"
)

var (
	// ErrUserFound should be returned from Create (see ConfirmUser)
	// when the primaryID of the record is found.
	ErrUserFound = errors.New("user found")
	// ErrUserNotFound should be returned from Get when the record is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrTokenNotFound should be returned from UseToken when the
	// record is not found.
	ErrTokenNotFound = errors.New("token not found")
)

// ServerStorer represents the data store that's capable of loading users
// and giving them a context with which to store themselves.
type ServerStorer interface {
	// Load will look up the user based on the passed the PrimaryID. Under
	// normal circumstances this comes from GetPID() of the user.
	Load(ctx context.Context, key string) (User, error)

	// Save persists the user in the database, this should never
	// create a user and instead return ErrUserNotFound if the user
	// does not exist.
	Save(ctx context.Context, user User) error
}

// CreatingServerStorer is used for creating new users
// like when Registration is being done.
type CreatingServerStorer interface {
	ServerStorer

	// New creates a blank user, it is not yet persisted in the database
	// but is just for storing data
	New(ctx context.Context) User
	// Create the user in storage, it should not overwrite a user
	// and should return ErrUserFound if it currently exists.
	Create(ctx context.Context, user User) error
}

// ConfirmingServerStorer can find a user by a confirm token
type ConfirmingServerStorer interface {
	ServerStorer

	// LoadByConfirmSelector finds a user by his confirm selector field
	// and should return ErrUserNotFound if that user cannot be found.
	LoadByConfirmSelector(ctx context.Context, selector string) (ConfirmableUser, error)
}

// RecoveringServerStorer allows users to be recovered by a token
type RecoveringServerStorer interface {
	ServerStorer

	// LoadByRecoverSelector finds a user by his recover selector field
	// and should return ErrUserNotFound if that user cannot be found.
	LoadByRecoverSelector(ctx context.Context, selector string) (RecoverableUser, error)
}

// RememberingServerStorer allows users to be remembered across sessions
type RememberingServerStorer interface {
	ServerStorer

	// AddRememberToken to a user
	AddRememberToken(ctx context.Context, pid, token string) error
	// DelRememberTokens removes all tokens for the given pid
	DelRememberTokens(ctx context.Context, pid string) error
	// UseRememberToken finds the pid-token pair and deletes it.
	// If the token could not be found return ErrTokenNotFound
	UseRememberToken(ctx context.Context, pid, token string) error
}

// EnsureCanCreate makes sure the server storer supports create operations
func EnsureCanCreate(storer ServerStorer) CreatingServerStorer {
	s, ok := storer.(CreatingServerStorer)
	if !ok {
		panic("could not upgrade ServerStorer to CreatingServerStorer, check your struct")
	}

	return s
}

// EnsureCanConfirm makes sure the server storer supports
// confirm-lookup operations
func EnsureCanConfirm(storer ServerStorer) ConfirmingServerStorer {
	s, ok := storer.(ConfirmingServerStorer)
	if !ok {
		panic("could not upgrade ServerStorer to ConfirmingServerStorer, check your struct")
	}

	return s
}

// EnsureCanRecover makes sure the server storer supports
// confirm-lookup operations
func EnsureCanRecover(storer ServerStorer) RecoveringServerStorer {
	s, ok := storer.(RecoveringServerStorer)
	if !ok {
		panic("could not upgrade ServerStorer to RecoveringServerStorer, check your struct")
	}

	return s
}

// EnsureCanRemember makes sure the server storer supports remember operations
func EnsureCanRemember(storer ServerStorer) RememberingServerStorer {
	s, ok := storer.(RememberingServerStorer)
	if !ok {
		panic("could not upgrade ServerStorer to RememberingServerStorer, check your struct")
	}

	return s
}
