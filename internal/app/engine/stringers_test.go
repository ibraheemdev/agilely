package engine

import (
	"testing"
)

func TestEventString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ev  AuthEvent
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
	var EventTest AuthEvent = -1
	if got := EventTest.String(); got != "AuthEvent(-1)" {
		t.Errorf("Wrong string for Event(%d) expected: 'AuthEvent(-1)', got: %s", EventTest, got)
	}
}
