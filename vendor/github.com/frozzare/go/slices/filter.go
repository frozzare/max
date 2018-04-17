package slices

import (
	"errors"
	"reflect"
)

func filter(input, predicate interface{}, compare bool) (interface{}, error) {
	var slice reflect.Value

	each(input, predicate, func(current, key, value reflect.Value) bool {
		if !slice.IsValid() {
			typ := reflect.SliceOf(value.Type())
			slice = reflect.MakeSlice(typ, 0, 0)
		}

		if current.Bool() == compare {
			slice = reflect.Append(slice, value)
		}

		return false
	})

	if slice.IsValid() {
		return slice.Interface(), nil
	}

	return nil, errors.New("Not a valid slice")
}

// Filter returns a new slice containing all values
// in the slice that satisfy the predicate function.
func Filter(input, predicate interface{}) (interface{}, error) {
	return filter(input, predicate, true)
}
