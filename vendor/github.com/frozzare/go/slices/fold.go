package slices

import (
	"errors"
	"reflect"
)

// Fold slice to single value using the predicate function.
func Fold(input, predicate interface{}) (interface{}, error) {
	var output reflect.Value

	l, item := source(input)

	if l == 0 {
		return nil, errors.New("Not a valid input slice")
	}

	iv := reflect.ValueOf(predicate)

	for i := 0; i < l; i++ {
		_, v := item(i)

		if !output.IsValid() {
			output = reflect.Zero(v.Type())
		}

		current := iv.Call([]reflect.Value{output, v})

		if len(current) > 0 {
			output = current[0]
		}
	}

	if output.IsValid() {
		return output.Interface(), nil
	}

	return nil, errors.New("Not a valid output value")
}
