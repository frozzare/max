package map2

import (
	"errors"
	"fmt"
	"reflect"
)

func mp(s interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Map {
		return v, fmt.Errorf("%s is not a map", v.Kind())
	}

	return v, nil
}

// Keys return the given map keys or error.
func Keys(m interface{}) (interface{}, error) {
	var slice reflect.Value

	s, err := mp(m)
	if err != nil {
		return nil, err
	}

	for _, k := range s.MapKeys() {
		if !slice.IsValid() {
			typ := reflect.SliceOf(k.Type())
			slice = reflect.MakeSlice(typ, 0, 0)
		}

		slice = reflect.Append(slice, k)
	}

	if slice.IsValid() {
		return slice.Interface(), nil
	}

	return nil, errors.New("No keys found")
}

// Values return the given map values or error.
func Values(m interface{}) (interface{}, error) {
	var slice reflect.Value

	s, err := mp(m)
	if err != nil {
		return nil, err
	}

	for _, k := range s.MapKeys() {
		v := s.MapIndex(k)

		if !slice.IsValid() {
			typ := reflect.SliceOf(v.Type())
			slice = reflect.MakeSlice(typ, 0, 0)
		}

		slice = reflect.Append(slice, v)
	}

	if slice.IsValid() {
		return slice.Interface(), nil
	}

	return nil, errors.New("No values found")
}
