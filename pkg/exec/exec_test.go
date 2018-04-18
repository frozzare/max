package exec

import (
	"os"
	"testing"
)

func TestCmd(t *testing.T) {
	path, _ := os.Getwd()

	if len(path) == 0 {
		t.Error("Expected: path, got: none")
	}

	if err := Exec("ls", path); err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
