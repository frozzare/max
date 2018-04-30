package strings2

import (
	"testing"

	assert "github.com/frozzare/go-assert"
)

func TestInsert(t *testing.T) {
	tests := map[string]struct {
		Start string
		Len   int
		End   string
	}{
		"HelWorldlo": {
			Start: "Hello",
			Len:   3,
			End:   "World",
		},
		"Hello world": {
			Start: "Hello ",
			Len:   6,
			End:   "world",
		},
		"Hello, 世 world 界": {
			Start: "Hello, 世界",
			Len:   10,
			End:   " world ",
		},
	}

	for k, v := range tests {
		assert.Equal(t, k, Insert(v.Start, v.Len, v.End))
	}
}
