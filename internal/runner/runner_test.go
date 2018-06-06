package runner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/frozzare/go/yaml2"
	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/task"
)

func TestRunner(t *testing.T) {
	var buf bytes.Buffer

	runner := New(
		Config(&config.Config{
			Tasks: map[string]*task.Task{
				"hello": {
					Summary:  "Hello task",
					Commands: yaml2.NewList("echo Hello $NAME"),
				},
			},
			Variables: map[string]string{
				"NAME": "Fredrik",
			},
		}),
		Verbose(true),
	)

	runner.Stdout = &buf

	if err := runner.Run("hello"); err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}

	got := strings.TrimSpace(buf.String())

	if got != "Hello Fredrik" {
		t.Errorf("Expected: 'Hello Fredrik', got: %s", got)
	}
}
