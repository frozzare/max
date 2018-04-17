package slices

import (
	"errors"
	"reflect"
)

// Uniq returns a new slice containing all values
// in the slice that satisfy the predicate function
// and are uniq.
func Uniq(input, predicate interface{}) (interface{}, error) {
	var mv reflect.Value
	var sv reflect.Value

	if predicate == nil {
		predicate = func(v interface{}) interface{} {
			return v
		}
	}

	each(input, predicate, func(current, key, value reflect.Value) bool {
		var cv reflect.Value
		var add bool

		if current.Kind() == reflect.Bool {
			cv = value
			add = current.Bool()
		} else {
			cv = current
			add = true
		}

		if !mv.IsValid() {
			mt := reflect.MapOf(cv.Type(), reflect.TypeOf(false))
			mv = reflect.MakeMap(mt)

			st := reflect.SliceOf(value.Type())
			sv = reflect.MakeSlice(st, 0, 0)
		}

		mi := mv.MapIndex(cv)

		if !mi.IsValid() && add {
			mv.SetMapIndex(cv, reflect.ValueOf(true))
			sv = reflect.Append(sv, value)
		}

		return false
	})

	if mv.IsValid() {
		return sv.Interface(), nil
	}

	return nil, errors.New("Not a valid slice")
}
