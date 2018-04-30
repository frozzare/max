package strings2

import (
	"regexp"
	"strconv"
)

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
