package exec

import (
	"testing"
)

func TestCmd(t *testing.T) {
	_, err := Cmd("ls")

	if err == nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
