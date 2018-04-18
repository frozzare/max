package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const httptask = `args:
  name:
summary: Hello task
commands:
- echo Hello {{ .name }}`

const httpinclude = `tasks:
  hello: !include %s`

func TestReadContentIncludeHttp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(httptask))
	}))

	defer server.Close()

	s := fmt.Sprintf(httpinclude, server.URL)

	c, err := ReadContent(s)

	if err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	if c == nil {
		t.Errorf("Expected: struct, got: %v", c)
	}

	if c.Tasks["hello"].Summary != "Hello task" {
		t.Errorf("Expected: 'Hello task', got: %v", c.Tasks["Hello"])
	}
}

func TestReadContent(t *testing.T) {
	buf, err := ioutil.ReadFile("./config.yml")

	if err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	c, err := ReadContent(string(buf))

	if err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	if c == nil {
		t.Errorf("Expected: struct, got: %v", c)
	}

	if c.Tasks["hello"].Summary != "Hello task" {
		t.Errorf("Expected: 'Hello task', got: %v", c.Tasks["Hello"])
	}
}

func TestReadFile(t *testing.T) {
	c, err := ReadFile("./config.yml")

	if err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	if c == nil {
		t.Errorf("Expected: struct, got: %v", c)
	}

	if c.Tasks["hello"].Summary != "Hello task" {
		t.Errorf("Expected: 'Hello task', got: %v", c.Tasks["Hello"])
	}

	if c.Tasks["hello2"] == nil {
		t.Errorf("Expected: task, got: %v", c.Tasks["hello2"])
	}

	if c.Tasks["hello2"].Summary != "Hello task" {
		t.Errorf("Expected: 'Hello task', got: %v", c.Tasks["Hello"])
	}
}
