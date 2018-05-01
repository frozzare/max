package runner

import (
	"testing"

	"github.com/frozzare/max/internal/backend/local"
	"github.com/frozzare/max/internal/config"
	"github.com/frozzare/max/internal/task"
	"github.com/frozzare/max/pkg/yamllist"
)

func TestRunner(t *testing.T) {
	runner := New(
		Engine(local.New()),
		Config(&config.Config{
			Tasks: map[string]*task.Task{
				"hello": &task.Task{
					Summary:  "Hello task",
					Commands: yamllist.NewList("echo Hello $NAME"),
				},
			},
			Variables: map[string]string{
				"NAME": "Fredrik",
			},
		}),
	)

	if err := runner.Run("hello"); err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
