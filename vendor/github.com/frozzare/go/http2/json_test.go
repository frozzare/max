package http2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frozzare/go-assert"
)

type JSONPerson struct {
	Name string
}

func TestGetJSON(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"name\":\"Fredrik\"}")
	}))

	actual := &JSONPerson{}
	err := GetJSON(server.URL, &actual)

	assert.Nil(t, err)
	assert.Equal(t, "Fredrik", actual.Name)
}

func TestGetJSONResponseError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	actual := &JSONPerson{}
	err := GetJSON(server.URL, &actual)

	assert.NotNil(t, err)
	assert.Equal(t, "", actual.Name)
}
