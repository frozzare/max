package strings2

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestChop(t *testing.T) {
	chop := Chop("fredrik", 3)
	res := []string{"fre", "dri", "k"}

	assert.Equal(t, chop[0], res[0])
	assert.Equal(t, chop[1], res[1])
	assert.Equal(t, chop[2], res[2])
}
