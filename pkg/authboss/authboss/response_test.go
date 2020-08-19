package authboss

import (
	"context"
	"testing"

	"github.com/ibraheemdev/agilely/pkg/mailer"
)

type testMailer struct{ sent bool }

func (t *testMailer) Send(context.Context, mailer.Email) error {
	t.sent = true
	return nil
}

func TestEmail(t *testing.T) {
	t.Parallel()

	ab := New()

	m := &testMailer{}
	renderer := &mockEmailRenderer{}
	ab.Config.Core.Mailer = m
	ab.Config.Core.MailRenderer = renderer

	email := mailer.Email{
		To:      []string{"support@authboss.com"},
		Subject: "Send help",
	}

	ro := EmailResponseOptions{
		Data:         nil,
		HTMLTemplate: "html",
		TextTemplate: "text",
	}

	if err := ab.Email(context.Background(), email, ro); err != nil {
		t.Error(err)
	}

	if !m.sent {
		t.Error("the e-mail should have been sent")
	}
}
