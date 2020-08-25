package engine

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

	e := New()

	m := &testMailer{}
	renderer := &mockEmailRenderer{}
	e.Config.Core.Mailer = m
	e.Config.Core.MailRenderer = renderer

	email := mailer.Email{
		To:      []string{"support@engine.com"},
		Subject: "Send help",
	}

	ro := EmailResponseOptions{
		Data:         nil,
		HTMLTemplate: "html",
		TextTemplate: "text",
	}

	if err := e.Email(context.Background(), email, ro); err != nil {
		t.Error(err)
	}

	if !m.sent {
		t.Error("the e-mail should have been sent")
	}
}
