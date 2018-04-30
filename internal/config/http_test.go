package config

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const httpTask = `args:
name: default
summary: Hello task
commands:
- echo Hello {{ .name }}
`

func TestIncludeHTTPTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(httpTask))
	}))

	defer server.Close()

	task, err := includeHTTPTask(server.URL, nil)

	if err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}

	if task != nil && task.Summary != "Hello task" {
		t.Errorf("Expected: 'Hello task', got: %s", task.Summary)
	}
}
