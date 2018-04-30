package task

import "testing"

func TestRenderEnvVariables(t *testing.T) {
	c := renderEnvVariables("echo Hello $NAME", map[string]string{
		"NAME": "Fredrik",
	})

	if c != "echo Hello Fredrik" {
		t.Fatal("Expected command to be 'echo Hello Fredrik'")
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
