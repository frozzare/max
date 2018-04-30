package task

import (
	"reflect"
	"testing"
)

func TestToEnv(t *testing.T) {
	vars := map[string]string{
		"TEST": "TEST",
	}

	env := toEnv(vars)

	if !reflect.DeepEqual([]string{"TEST=TEST"}, env) {
		t.Fatalf("Expected slice to be the same")
	}
}
