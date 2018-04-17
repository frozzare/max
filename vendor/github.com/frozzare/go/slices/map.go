package slices

import (
	"errors"
	"reflect"
)

// Map creates a new slice with the results of calling a
// provided predicate on every item in the calling slice.
func Map(input, predicate interface{}) (interface{}, error) {
	var slice reflect.Value

	each(input, predicate, func(current, key, value reflect.Value) bool {
		if !slice.IsValid() {
			typ := reflect.SliceOf(current.Type())
			slice = reflect.MakeSlice(typ, 0, 0)
		}

		slice = reflect.Append(slice, current)

		return false
	})

	if slice.IsValid() {
		return slice.Interface(), nil
	}

	return nil, errors.New("Not a valid slice")
}
