package url2

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/frozzare/go-assert"
)

func TestBoolQuery(t *testing.T) {
	tests := []struct {
		Default  interface{}
		Expected bool
		Key      string
		Request  *http.Request
	}{
		{
			Default:  false,
			Expected: false,
			Key:      "",
			Request:  &http.Request{},
		},
		{
			Default:  false,
			Expected: true,
			Key:      "foo",
			Request: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "www.google.com",
					Path:     "/",
					RawQuery: "foo=true",
				},
			},
		},
		{
			Default:  false,
			Expected: false,
			Key:      "",
			Request: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, Query(test.Request).Bool(test.Key, test.Default))
	}
}

func TestFloatQuery(t *testing.T) {
	tests := []struct {
		Default  interface{}
		Expected float64
		Key      string
		Request  *http.Request
	}{
		{
			Default:  0,
			Expected: 0.0,
			Key:      "",
			Request:  &http.Request{},
		},
		{
			Default:  0,
			Expected: 32.0,
			Key:      "foo",
			Request: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "www.google.com",
					Path:     "/",
					RawQuery: "foo=32",
				},
			},
		},
		{
			Default:  false,
			Expected: 0.0,
			Key:      "",
			Request: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, Query(test.Request).Float(test.Key, test.Default))
	}
}

func TestIntQuery(t *testing.T) {
	tests := []struct {
		Default  interface{}
		Expected int
		Key      string
		Request  *http.Request
	}{
		{
			Default:  0,
			Expected: 0,
			Key:      "",
			Request:  &http.Request{},
		},
		{
			Default:  0,
			Expected: 32,
			Key:      "foo",
			Request: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "www.google.com",
					Path:     "/",
					RawQuery: "foo=32",
				},
			},
		},
		{
			Default:  false,
			Expected: 0,
			Key:      "",
			Request: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, Query(test.Request).Int(test.Key, test.Default))
	}
}

func TestGetQuery(t *testing.T) {
	tests := []struct {
		Default  interface{}
		Expected string
		Key      string
		Request  *http.Request
	}{
		{
			Default:  "",
			Expected: "",
			Key:      "",
			Request:  &http.Request{},
		},
		{
			Default:  "",
			Expected: "32",
			Key:      "foo",
			Request: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "www.google.com",
					Path:     "/",
					RawQuery: "foo=32",
				},
			},
		},
		{
			Default:  false,
			Expected: "",
			Key:      "",
			Request: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, Query(test.Request).Get(test.Key, test.Default))
	}
}
