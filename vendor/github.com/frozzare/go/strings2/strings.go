package strings2

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Chars returns all characters a string slice.
func Chars(s string) []string {
	res := []string{}

	for i := 0; i < utf8.RuneCountInString(s); i++ {
		res = append(res, string(s[i]))
	}

	return res
}

// Chop chops the string with the given length and returns a string slice.
func Chop(s string, l int) []string {
	r, _ := regexp.Compile(".{1," + strconv.Itoa(l) + "}")
	matches := r.FindAllStringSubmatch(s, -1)
	res := []string{}

	for i := 0; i < len(matches); i++ {
		res = append(res, string(matches[i][0]))
	}

	return res
}

// Insert inserts a new string at a given index and returns the string.
func Insert(s string, i int, a string) string {
	return s[0:i] + a + s[i:]
}

// Lines returns the lines in the string as a string slice.
func Lines(s string) []string {
	return strings.Split(s, "\n")
}

// Reverse returns the string reversed.
func Reverse(s string) string {
	runes := []rune(s)

	for i, l := 0, len(runes)-1; i < l; i, l = i+1, l-1 {
		runes[i], runes[l] = runes[l], runes[i]
	}

	return string(runes)
}

// SwapCase returns a string where uppercase is swapped with lowercase and lowercase to uppercase and vice versa.
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
