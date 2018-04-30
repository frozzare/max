package strings2

import "unicode/utf8"

// Chars returns all characters a string slice.
func Chars(s string) []string {
	res := []string{}

	for i := 0; i < utf8.RuneCountInString(s); i++ {
		res = append(res, string(s[i]))
	}

	return res
}
