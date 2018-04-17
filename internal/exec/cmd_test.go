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

	_, err := Cmd("ls", path)

	if err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
