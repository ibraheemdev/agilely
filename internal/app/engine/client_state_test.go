package engine

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStateGet(t *testing.T) {
	t.Parallel()

	e := New()
	e.Core.SessionState = newMockClientStateRW("one", "two")
	e.Core.CookieState = newMockClientStateRW("three", "four")

	r := httptest.NewRequest("GET", "/", nil)
	w := e.NewResponse(httptest.NewRecorder())

	var err error
	r, err = e.LoadClientState(w, r)
	if err != nil {
		t.Error(err)
	}

	if got, _ := GetSession(r, "one"); got != "two" {
		t.Error("session value was wrong:", got)
	}
	if got, _ := GetCookie(r, "three"); got != "four" {
		t.Error("cookie value was wrong:", got)
	}
}

func TestStateResponseWriterDoubleWritePanic(t *testing.T) {
	t.Parallel()

	e := New()
	e.Core.SessionState = newMockClientStateRW("one", "two")

	w := e.NewResponse(httptest.NewRecorder())

	w.WriteHeader(200)
	// Check this doesn't panic
	w.WriteHeader(200)

	defer func() {
		if recover() == nil {
			t.Error("expected a panic")
		}
	}()

	_ = w.putClientState()
}

func TestStateResponseWriterLastSecondWriteHeader(t *testing.T) {
	t.Parallel()

	e := New()
	e.Core.SessionState = newMockClientStateRW()

	w := e.NewResponse(httptest.NewRecorder())

	PutSession(w, "one", "two")

	w.WriteHeader(200)
	got := strings.TrimSpace(w.Header().Get("test_session"))
	if got != `{"one":"two"}` {
		t.Error("got:", got)
	}
}

func TestStateResponseWriterLastSecondWriteWrite(t *testing.T) {
	t.Parallel()

	e := New()
	e.Core.SessionState = newMockClientStateRW()

	w := e.NewResponse(httptest.NewRecorder())

	PutSession(w, "one", "two")

	io.WriteString(w, "Hello world!")

	got := strings.TrimSpace(w.Header().Get("test_session"))
	if got != `{"one":"two"}` {
		t.Error("got:", got)
	}
}

func TestStateResponseWriterAuthEvents(t *testing.T) {
	t.Parallel()

	e := New()
	w := e.NewResponse(httptest.NewRecorder())

	PutSession(w, "one", "two")
	DelSession(w, "one")
	DelCookie(w, "one")
	PutCookie(w, "two", "one")

	want := ClientStateAuthEvent{Kind: ClientStateAuthEventPut, Key: "one", Value: "two"}
	if got := w.sessionStateAuthEvents[0]; got != want {
		t.Error("event was wrong", got)
	}

	want = ClientStateAuthEvent{Kind: ClientStateAuthEventDel, Key: "one"}
	if got := w.sessionStateAuthEvents[1]; got != want {
		t.Error("event was wrong", got)
	}

	want = ClientStateAuthEvent{Kind: ClientStateAuthEventDel, Key: "one"}
	if got := w.cookieStateAuthEvents[0]; got != want {
		t.Error("event was wrong", got)
	}

	want = ClientStateAuthEvent{Kind: ClientStateAuthEventPut, Key: "two", Value: "one"}
	if got := w.cookieStateAuthEvents[1]; got != want {
		t.Error("event was wrong", got)
	}
}

func TestFlashClearer(t *testing.T) {
	t.Parallel()

	e := New()
	e.Core.SessionState = newMockClientStateRW(FlashSuccessKey, "a", FlashErrorKey, "b")

	r := httptest.NewRequest("GET", "/", nil)
	w := e.NewResponse(httptest.NewRecorder())

	if msg := FlashSuccess(w, r); msg != "" {
		t.Error("unexpected flash success:", msg)
	}

	if msg := FlashError(w, r); msg != "" {
		t.Error("unexpected flash error:", msg)
	}

	var err error
	r, err = e.LoadClientState(w, r)
	if err != nil {
		t.Error(err)
	}

	if msg := FlashSuccess(w, r); msg != "a" {
		t.Error("Unexpected flash success:", msg)
	}

	if msg := FlashError(w, r); msg != "b" {
		t.Error("Unexpected flash error:", msg)
	}

	want := ClientStateAuthEvent{Kind: ClientStateAuthEventDel, Key: FlashSuccessKey}
	if got := w.sessionStateAuthEvents[0]; got != want {
		t.Error("event was wrong", got)
	}
	want = ClientStateAuthEvent{Kind: ClientStateAuthEventDel, Key: FlashErrorKey}
	if got := w.sessionStateAuthEvents[1]; got != want {
		t.Error("event was wrong", got)
	}
}

func TestDelAllSession(t *testing.T) {
	t.Parallel()

	csrw := &ClientStateResponseWriter{}

	DelAllSession(csrw, []string{"notthisone", "orthis"})

	if len(csrw.sessionStateAuthEvents) != 1 {
		t.Error("should have one delete all")
	}
	if ev := csrw.sessionStateAuthEvents[0]; ev.Kind != ClientStateAuthEventDelAll {
		t.Error("it should be a delete all event:", ev.Kind)
	} else if ev.Key != "notthisone,orthis" {
		t.Error("the whitelist should be passed through as CSV:", ev.Key)
	}
}
