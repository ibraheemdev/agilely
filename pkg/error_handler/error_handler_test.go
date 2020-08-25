package errorhandler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ibraheemdev/agilely/pkg/logger"
)

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	b := &bytes.Buffer{}

	eh := ErrorHandler{LogWriter: logger.New(b)}

	handler := eh.Wrap(func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("error occurred")
	})
	// Assert that it's the right type
	var _ http.Handler = handler

	handler.ServeHTTP(nil, httptest.NewRequest("GET", "/target", nil))

	if !strings.Contains(b.String(), "error from (192.0.2.1:1234) /target: error occurred") {
		t.Error("output was wrong:", b.String())
	}
}
