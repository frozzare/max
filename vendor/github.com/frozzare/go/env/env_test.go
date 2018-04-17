package env

import (
	"os"
	"testing"

	"github.com/frozzare/go-assert"
)

func TestGet(t *testing.T) {
	v := Get("NAME")

	assert.Empty(t, v)

	os.Setenv("NAME", "")

	v = Get("NAME", "fredrik")

	assert.Equal(t, "fredrik", v)

	Set("NAME", "go")

	v = Get("NAME", "fredrik")

	assert.Equal(t, "go", v)
}

func TestLoad(t *testing.T) {
	test := map[string]string{
		"FOO": "BAR",
		"URL": "https://golang.org",
	}
	files := []string{"testdata/a.env", "testdata/b.env"}

	for k := range test {
		assert.Empty(t, Get(k))
	}

	err := Load(files...)
	assert.Nil(t, err)

	for k, v := range test {
		assert.Equal(t, v, Get(k))
	}

	err = Load("testdata/c.env")
	assert.NotNil(t, err)
}

func TestSet(t *testing.T) {
	os.Setenv("NAME", "")

	v := Get("NAME", "fredrik")

	assert.Equal(t, "fredrik", v)

	Set("NAME", "go")

	v = Get("NAME", "fredrik")

	assert.Equal(t, "go", v)
}
