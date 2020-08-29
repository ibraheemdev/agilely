package engine

import (
	"testing"
)

func TestCTXKeyString(t *testing.T) {
	t.Parallel()

	if got := CTXKeyData.String(); got != "engine ctx key data" {
		t.Error(got)
	}
}
