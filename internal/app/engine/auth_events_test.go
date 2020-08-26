package engine

import (
	"errors"
	"net/http"
	"testing"
)

func TestEvents(t *testing.T) {
	t.Parallel()

	e := New()
	afterCalled := false
	beforeCalled := false

	e.AuthEvents.Before(EventRegister, func(http.ResponseWriter, *http.Request, bool) (bool, error) {
		beforeCalled = true
		return false, nil
	})
	e.AuthEvents.After(EventRegister, func(http.ResponseWriter, *http.Request, bool) (bool, error) {
		afterCalled = true
		return false, nil
	})

	if beforeCalled || afterCalled {
		t.Error("Neither should be called.")
	}

	handled, err := e.AuthEvents.FireBefore(EventRegister, nil, nil)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if handled {
		t.Error("It should not have been handled.")
	}

	if !beforeCalled {
		t.Error("Expected before to have been called.")
	}
	if afterCalled {
		t.Error("Expected after not to be called.")
	}

	e.AuthEvents.FireAfter(EventRegister, nil, nil)
	if !afterCalled {
		t.Error("Expected after to be called.")
	}
}

func TestEventsHandled(t *testing.T) {
	t.Parallel()

	e := New()
	firstCalled := false
	secondCalled := false

	firstHandled := false
	secondHandled := false

	e.AuthEvents.Before(EventRegister, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
		firstCalled = true
		firstHandled = handled
		return true, nil
	})
	e.AuthEvents.Before(EventRegister, func(w http.ResponseWriter, r *http.Request, handled bool) (bool, error) {
		secondCalled = true
		secondHandled = handled
		return false, nil
	})

	handled, err := e.AuthEvents.FireBefore(EventRegister, nil, nil)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !handled {
		t.Error("it should have been handled")
	}

	if !firstCalled {
		t.Error("expected first to have been called")
	}
	if !secondCalled {
		t.Error("expected second to have been called")
	}

	if firstHandled {
		t.Error("first should not see the event as being handled")
	}
	if !secondHandled {
		t.Error("second should see the event as being handled")
	}
}

func TestEventsErrors(t *testing.T) {
	t.Parallel()

	e := New()
	firstCalled := false
	secondCalled := false

	expect := errors.New("error")

	e.AuthEvents.Before(EventRegister, func(http.ResponseWriter, *http.Request, bool) (bool, error) {
		firstCalled = true
		return false, expect
	})
	e.AuthEvents.Before(EventRegister, func(http.ResponseWriter, *http.Request, bool) (bool, error) {
		secondCalled = true
		return false, nil
	})

	_, err := e.AuthEvents.FireBefore(EventRegister, nil, nil)
	if err != expect {
		t.Error("got the wrong error back:", err)
	}

	if !firstCalled {
		t.Error("expected first to have been called")
	}
	if secondCalled {
		t.Error("expected second to not have been called")
	}
}
