package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestReverse(t *testing.T) {
	assert.Equal(t, "raboof", Reverse("foobar"))
}
