package slices

// First returns the first item in a slice or nil.
func First(input interface{}) interface{} {
	l, item := source(input)

	if l == 0 {
		return nil
	}

	_, v := item(0)
	return v.Interface()
}
