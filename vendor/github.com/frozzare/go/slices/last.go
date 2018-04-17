package slices

// Last returns the last item in a slice or nil.
func Last(input interface{}) interface{} {
	l, item := source(input)

	if l == 0 {
		return nil
	}

	_, v := item(l - 1)
	return v.Interface()
}
