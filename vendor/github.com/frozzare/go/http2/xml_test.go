package http2

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frozzare/go-assert"
)

type XMLPerson struct {
	XMLName xml.Name `xml:"name"`
	Name    string   `xml:",chardata"`
}

func TestGetXML(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<name>Fredrik</name>")
	}))

	actual := &XMLPerson{}
	err := GetXML(server.URL, &actual)

	assert.Nil(t, err)
	assert.Equal(t, "Fredrik", actual.Name)
}

func TestGetXMLResponseError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	actual := &XMLPerson{}
	err := GetXML(server.URL, &actual)

	assert.NotNil(t, err)
	assert.Equal(t, "", actual.Name)
}
