package defaults

import (
	"testing"

	"github.com/ibraheemdev/poller/pkg/authboss/authboss"
)

func TestSetCore(t *testing.T) {
	t.Parallel()

	config := &authboss.Config{}
	SetCore(config, false, false, "/auth", "/templates", "/layouts")

	if config.Core.Logger == nil {
		t.Error("logger should be set")
	}
	if config.Core.Router == nil {
		t.Error("router should be set")
	}
	if config.Core.ErrorHandler == nil {
		t.Error("error handler should be set")
	}
	if config.Core.Responder == nil {
		t.Error("responder should be set")
	}
	if config.Core.Redirector == nil {
		t.Error("redirector should be set")
	}
	if config.Core.BodyReader == nil {
		t.Error("bodyreader should be set")
	}
	if config.Core.Mailer == nil {
		t.Error("mailer should be set")
	}
	if config.Core.Logger == nil {
		t.Error("logger should be set")
	}
	if config.Core.ViewRenderer == nil {
		t.Error("view renderer should be set")
	}
	if config.Core.MailRenderer == nil {
		t.Error("mail renderer should be set")
	}
}
