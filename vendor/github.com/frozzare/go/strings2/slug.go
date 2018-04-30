package strings2

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/rainycape/unidecode"
)

var (
	regCleanString        = regexp.MustCompile(`[^a-zA-Z0-9\s\-\_]+`)
	regNonAuthorizedChars = regexp.MustCompile(`[^a-z\d]+`)
	regMultipleDashes     = regexp.MustCompile(`-+`)
)

func decamelize(s string) string {
	var w []string
	l := 0

	for i := s; i != ""; i = i[l:] {
		l = strings.IndexFunc(i[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(i)
		}
		w = append(w, i[:l])
	}

	return strings.Join(w, " ")
}

// Slug generates a slug from a input string.
func Slug(s string) string {
	s = strings.TrimSpace(s)
	s = decamelize(s)
	s = unidecode.Unidecode(s)
	s = strings.ToLower(s)

	// Clean string and remove bad chars.
	s = regCleanString.ReplaceAllString(s, "")

	// Replace all non authorized chars with a dash.
	s = regNonAuthorizedChars.ReplaceAllString(s, "-")

	// Replace all multiple dahses with a single dash.
	s = regMultipleDashes.ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}
