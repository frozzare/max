package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Set("test", []byte("test")); err != nil {
		t.Fatal(err)
	}

	v, err := c.Get("test")
	if err != nil {
		t.Fatal(err)
	}

	if string(v) != "test" {
		t.Fatal("Expected test value to be test")
	}

	if err := c.Delete("test"); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Get("test"); err == nil {
		t.Fatal(err)
	}
}
