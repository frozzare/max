package task

import (
	"testing"

	"github.com/frozzare/go/env"
	"github.com/frozzare/max/pkg/yamllist"
)

func TestTask(t *testing.T) {
	s := &Task{
		Summary:  "Hello task",
		Commands: yamllist.NewList("echo $SAY"),
	}

	env.Set("SAY", "Hello")

	err := s.Run(map[string]interface{}{})

	if err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
