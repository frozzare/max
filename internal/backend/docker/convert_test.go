package docker

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

func TestToVolumes(t *testing.T) {
	vols := toVolumes([]string{"/var/test:/var/test2"})
	exp := map[string]struct{}{"/var/test2": {}}

	if !reflect.DeepEqual(vols, exp) {
		t.Fatalf("Expected maps to be the same")
	}
}
