package strings2

import "strings"

// Lines returns the lines in the string as a string slice.
func Lines(s string) []string {
	return strings.Split(s, "\n")
}
