package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestChars(t *testing.T) {
	tests := map[string]int{
		"fredrik":   7,
		"åäö":       3,
		"Hello, 世界": 9,
	}

	for k, v := range tests {
		assert.Equal(t, len(Chars(k)), v)
	}
}
