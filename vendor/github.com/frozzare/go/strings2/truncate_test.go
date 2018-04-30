package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestTruncate(t *testing.T) {
	assert.Equal(t, "Hello...", Truncate("Hello world", 5))
	assert.Equal(t, "Hello.....", Truncate("Hello world", 5, 5))
	assert.Equal(t, "Hello-----", Truncate("Hello world", 5, 5, "-"))
	assert.Equal(t, "Hello", Truncate("Hello", 10))
}
