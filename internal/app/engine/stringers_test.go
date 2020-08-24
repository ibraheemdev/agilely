package engine

import (
	"testing"
)

func TestEventString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ev  Event
		str string
	}{
		{EventRegister, "EventRegister"},
		{EventAuth, "EventAuth"},
		{EventOAuth2, "EventOAuth2"},
		{EventAuthFail, "EventAuthFail"},
		{EventOAuth2Fail, "EventOAuth2Fail"},
		{EventRecoverStart, "EventRecoverStart"},
		{EventRecoverEnd, "EventRecoverEnd"},
		{EventGetUser, "EventGetUser"},
		{EventGetUserSession, "EventGetUserSession"},
		{EventPasswordReset, "EventPasswordReset"},
	}

	for i, test := range tests {
		if got := test.ev.String(); got != test.str {
			t.Errorf("%d) Wrong string for Event(%d) expected: %v got: %s", i, test.ev, test.str, got)
		}
	}

	// This test is only for 100% test coverage of stringers.go
	var testEvent Event = -1
	if got := testEvent.String(); got != "Event(-1)" {
		t.Errorf("Wrong string for Event(%d) expected: \"\", got: %s", testEvent, got)
	}
}
