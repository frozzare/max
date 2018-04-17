package task

import (
	"testing"

	"github.com/frozzare/go/env"
)

func TestTask(t *testing.T) {
	s := &Task{
		Summary: "Hello task",
		Commands: []string{
			"echo $SAY",
		},
	}

	env.Set("SAY", "Hello")

	err := s.Run([]interface{}{})

	if err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
