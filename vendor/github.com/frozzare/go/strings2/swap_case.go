package strings2

import (
	"strings"
)

// SwapCase returns a string where uppercase is swapped with
// lowercase and lowercase to uppercase and vice versa.
func SwapCase(s string) string {
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		c := string(runes[i])
		if c == strings.ToUpper(c) {
			runes[i] = []rune(strings.ToLower(c))[0]
		} else {
			runes[i] = []rune(strings.ToUpper(c))[0]
		}
	}

	return string(runes)
}
