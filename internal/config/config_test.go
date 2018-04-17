package config

import (
	"io/ioutil"
	"testing"
)

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
		t.Errorf("Expected: 'Hello task', got: %s", c.Tasks["Hello"])
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
		t.Errorf("Expected: 'Hello task', got: %s", c.Tasks["Hello"])
	}

	if c.Tasks["hello2"] == nil {
		t.Errorf("Expected: task, got: %v", c.Tasks["hello2"])
	}

	if c.Tasks["hello2"].Summary != "Hello task" {
		t.Errorf("Expected: 'Hello task', got: %s", c.Tasks["Hello"])
	}
}
