package strings2

// Reverse returns the string reversed.
func Reverse(s string) string {
	runes := []rune(s)

	for i, l := 0, len(runes)-1; i < l; i, l = i+1, l-1 {
		runes[i], runes[l] = runes[l], runes[i]
	}

	return string(runes)
}
