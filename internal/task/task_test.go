package task

import (
	"context"
	"testing"

	"github.com/frozzare/go/yaml2"
)

func TestPrepare(t *testing.T) {
	task := &Task{
		Args: map[string]interface{}{
			"name": "Fredrik",
		},
		Dir:      "./{{ .name }}",
		Commands: yaml2.NewList("Hello {{ .name }}"),
		Status:   yaml2.NewList("Hello {{ .name }}"),
	}

	task.Prepare()

	if task.Dir != "./Fredrik" {
		t.Fatalf("Expected dir value to './Fredrik', got: %s", task.Dir)
	}

	if task.Commands.Values[0] != "Hello Fredrik" {
		t.Fatalf("Expected command value to be 'Hello Fredrik', got: %s", task.Commands.Values[0])
	}

	if task.Status.Values[0] != "Hello Fredrik" {
		t.Fatalf("Expected status value to be 'Hello Fredrik', got: %s", task.Status.Values[0])
	}
}

func TestUpToDate(t *testing.T) {
	task := &Task{}

	if task.UpToDate(context.Background()) {
		t.Fatal("Expected task to not be up to date")
	}

	task = &Task{
		Status: yaml2.NewList("test -z \"\""),
	}

	if !task.UpToDate(context.Background()) {
		t.Fatal("Expected task to be up to date")
	}
}
