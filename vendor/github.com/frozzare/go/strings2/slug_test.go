package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestSlug(t *testing.T) {
	tests := map[string]string{
		"Foo bar":                 "foo-bar",
		"foo bar baz":             "foo-bar-baz",
		"foo bar ":                "foo-bar",
		"   foo bar  ":            "foo-bar",
		"[foo] [bar]":             "foo-bar",
		"Foo ÿ":                   "foo-y",
		"FooBar":                  "foo-bar",
		"fooBar":                  "foo-bar",
		"Foo & Bar":               "foo-bar",
		"Hallå världen":           "halla-varlden",
		"smile ☺":                 "smile",
		"Hellö Wörld хелло ворлд": "hello-world-khello-vorld",
		"\"C'est déjà l’été.\"":   "cest-deja-lete",
		"jaja---lol-méméméoo--a":  "jaja-lol-mememeoo-a",
	}

	for k, v := range tests {
		assert.Equal(t, Slug(k), v)
	}
}
