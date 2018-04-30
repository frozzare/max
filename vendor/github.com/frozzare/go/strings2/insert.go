package strings2

// Insert inserts a new string at a given index and returns the string.
func Insert(s string, i int, a string) string {
	return s[0:i] + a + s[i:]
}
