package slices

import (
	"reflect"
)

// All returns true if all of the values
// in the slice satisfies the predicate function.
func All(input, predicate interface{}) bool {
	var result bool

	each(input, predicate, func(current, key, value reflect.Value) bool {
		result = current.Bool()

		return !result
	})

	return result
}
