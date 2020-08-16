package defaults

import (
	"context"
	"testing"

	"github.com/ibraheemdev/poller/pkg/authboss/authboss"
)

func TestJSONRenderer(t *testing.T) {
	t.Parallel()

	r := JSONRenderer{}

	success := authboss.HTMLData{"fun": "times"}
	failure := authboss.HTMLData{authboss.DataErr: "problem"}
	hasAlready := authboss.HTMLData{authboss.DataErr: "problem", "status": "noproblem"}

	b, _, err := r.Render(context.Background(), "", nil)
	if err != nil {
		t.Error(err)
	}
	if string(b) != `{"status":"success"}` {
		t.Errorf("wrong json: %s", b)
	}

	b, _, err = r.Render(context.Background(), "", success)
	if err != nil {
		t.Error(err)
	}
	if string(b) != `{"fun":"times","status":"success"}` {
		t.Errorf("wrong json: %s", b)
	}

	b, _, err = r.Render(context.Background(), "", failure)
	if err != nil {
		t.Error(err)
	}
	if string(b) != `{"error":"problem","status":"failure"}` {
		t.Errorf("wrong json: %s", b)
	}

	b, _, err = r.Render(context.Background(), "", hasAlready)
	if err != nil {
		t.Error(err)
	}
	if string(b) != `{"error":"problem","status":"noproblem"}` {
		t.Errorf("wrong json: %s", b)
	}
}

func TestJSONLoad(t *testing.T) {
	t.Parallel()
	r := JSONRenderer{}
	l := r.Load("blabla.json")

	if l != nil {
		t.Errorf("expected nil, got: %s", l)
	}
}
