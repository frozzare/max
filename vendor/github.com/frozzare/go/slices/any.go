package slices

import (
	"reflect"
)

// Any returns true if one of the values
// in the slice satisfies the predicate function.
func Any(input, predicate interface{}) bool {
	var result bool

	each(input, predicate, func(current, key, value reflect.Value) bool {
		result = current.Bool()

		return result
	})

	return result
}
