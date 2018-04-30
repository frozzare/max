package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestLines(t *testing.T) {
	assert.Equal(t, len(Lines("Hello\nWorld")), 2)
}
