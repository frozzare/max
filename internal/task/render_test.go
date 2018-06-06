package task

import (
	"testing"

	"github.com/frozzare/go/env"
	"github.com/frozzare/go/yaml2"
	"github.com/frozzare/max/internal/backend/config"
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

	c = renderEnvVariables("echo Hello $1", map[string]string{
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

func TestRenderStruct(t *testing.T) {
	task := &Task{
		Args: map[string]interface{}{
			"name": "Fredrik",
		},
		Summary:  "Hello task",
		Commands: yaml2.NewList("echo Hello $NAME"),
		Docker: &config.Docker{
			Auth: &config.Auth{
				Username: "{{ .name }}",
			},
			Image: "$NAME",
		},
		Variables: map[string]string{
			"NAME": "Fredrik",
		},
	}

	v, err := renderStruct(task, task.Args, task.Variables)

	if err != nil {
		t.Fatal("Expected error to be nil")
	}

	task = v.(*Task)

	if task.Commands.Values[0] != "echo Hello Fredrik" {
		t.Fatal("Expected command to be 'echo Hello Fredrik'")
	}

	if task.Docker.Auth.Username != "Fredrik" {
		t.Fatal("Expected username value to be 'Fredrik'")
	}

	if task.Docker.Image != "Fredrik" {
		t.Fatal("Expected image value to be 'Fredrik'")
	}
}
