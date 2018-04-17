package http2

import (
	"testing"

	"net/http"
	"net/http/cookiejar"

	"github.com/frozzare/go-assert"
)

func TestNewClient(t *testing.T) {
	assert.NotNil(t, NewClient(nil))
}

func TestNewClientConfig(t *testing.T) {
	jar, _ := cookiejar.New(nil)

	c := NewClient(&http.Client{
		Jar: jar,
	})

	assert.Equal(t, jar, c.Jar)
}
