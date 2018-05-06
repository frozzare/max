package task

import (
	"testing"

	"github.com/frozzare/go/env"
)

func TestRenderEnvVariables(t *testing.T) {
	c := renderEnvVariables("echo Hello $NAME", map[string]string{
		"NAME": "Fredrik",
	})

	if c != "echo Hello Fredrik" {
		t.Fatal("Expected command to be 'echo Hello Fredrik'")
	}

	env.Set("NAME", "Fredrik2")

	c = renderEnvVariables("echo Hello $NAME", map[string]string{
		"NAME": "Fredrik",
	})

	if c != "echo Hello Fredrik2" {
		t.Fatal("Expected command to be 'echo Hello Fredrik2'")
	}

	env.Set("NAME", "")

	c = renderEnvVariables("echo HELLO $1", map[string]string{
		"1": "Fredrik3",
	})

	if c != "echo Hello Fredrik3" {
		t.Fatal("Expected command to be 'echo Hello Fredrik3'")
	}
}

func TestRenderCommand(t *testing.T) {
	c, err := renderCommand("echo Hello {{ .name }}", map[string]interface{}{
		"name": "Fredrik",
	})

	if err != nil {
		t.Fatal("Expected error to be nil")
	}

	if c != "echo Hello Fredrik" {
		t.Fatal("Expected command to be 'echo Hello Fredrik'")
	}
}
