package slices

// Reject returns a new slice containing all values
// in the slice that don't satisfy the predicate function.
func Reject(input, predicate interface{}) (interface{}, error) {
	return filter(input, predicate, false)
}
