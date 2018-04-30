package strings2

// Truncate will truncate the string and add dots or given character at the end.
// The additional values are length of the dots or given character and the character
// that should be used.
func Truncate(s string, w int, args ...interface{}) string {
	c := "."
	l := 3

	if len(args) > 0 {
		l = args[0].(int)
	}

	if len(args) > 1 {
		c = args[1].(string)
	}

	if len(s) < w {
		return s
	}

	dots := ""

	for i := 0; i < l; i++ {
		dots = dots + c
	}

	return s[0:w] + dots
}
