package url2

import (
	"net/http"
	"reflect"
	"strconv"
)

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.
type Values map[string][]string

// Query returns query string values from a http request.
func Query(r *http.Request) Values {
	v := Values{}

	if r == nil || r.URL == nil {
		return v
	}

	for key, value := range r.URL.Query() {
		v[key] = value
	}

	return v
}

// Bool gets the first value associated with the given key.
// If there are no values associated with the key, Bool returns
// false. To access multiple values, use the map directly.
func (v Values) Bool(key string, def ...interface{}) bool {
	s := v.Get(key)
	d := value(false, def...).(bool)

	if len(s) == 0 {
		return d
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return d
	}

	return b
}

// Float gets the first value associated with the given key.
// If there are no values associated with the key, Float returns
// zero. To access multiple values, use the map directly.
func (v Values) Float(key string, def ...interface{}) float64 {
	s := v.Get(key)
	d := value(float64(0.0), def...).(float64)

	if len(s) == 0 {
		return d
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return d
	}

	return f
}

// Int gets the first value associated with the given key.
// If there are no values associated with the key, Int returns
// zero. To access multiple values, use the map directly.
func (v Values) Int(key string, def ...interface{}) int {
	s := v.Get(key)
	d := value(0, def...).(int)

	if len(s) == 0 {
		return d
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}

	return i
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v Values) Get(key string, def ...interface{}) string {
	d := value("", def...).(string)

	if v == nil {
		return d
	}

	vs := v[key]
	if len(vs) == 0 {
		return d
	}

	if v := vs[0]; len(v) > 0 {
		return v
	}

	return d
}

// Set sets the key to value. It replaces any existing
// values.
func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v Values) Del(key string) {
	delete(v, key)
}

func value(def interface{}, values ...interface{}) interface{} {
	if len(values) == 0 {
		return def
	}

	v := values[0]
	if reflect.TypeOf(def).Kind() == reflect.TypeOf(v).Kind() {
		return v
	}

	return def
}
