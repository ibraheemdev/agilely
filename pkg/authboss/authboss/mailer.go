package authboss

import (
	"context"

	"github.com/ibraheemdev/agilely/pkg/mailer"
)

// Mailer is a type that is capable of sending an e-mail.
type Mailer interface {
	Send(context.Context, mailer.Email) error
}
