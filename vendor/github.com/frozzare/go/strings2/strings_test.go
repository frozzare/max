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

func TestChars(t *testing.T) {
	assert.Equal(t, len(Chars("fredrik")), 7)
	assert.Equal(t, len(Chars("åäö")), 3)
	assert.Equal(t, len(Chars("Hello, 世界")), 9)
}

func TestSwapCase(t *testing.T) {
	assert.Equal(t, "FredriK", SwapCase("fREDRIk"))
	assert.Equal(t, "Åääö", SwapCase("åÄÄÖ"))
	assert.Equal(t, "hELLO, 世界", SwapCase("Hello, 世界"))
}

func TestInsert(t *testing.T) {
	assert.Equal(t, "Helworldlo", Insert("Hello", 3, "world"))
	assert.Equal(t, "Hello world", Insert("Hello ", 6, "world"))
	assert.Equal(t, "Hello, 世 world 界", Insert("Hello, 世界", 10, " world "))
}

func TestLines(t *testing.T) {
	assert.Equal(t, len(Lines("Hello\nWorld")), 2)
}

func TestReverse(t *testing.T) {
	assert.Equal(t, "raboof", Reverse("foobar"))
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "Hello...", Truncate("Hello world", 5))
	assert.Equal(t, "Hello.....", Truncate("Hello world", 5, 5))
	assert.Equal(t, "Hello-----", Truncate("Hello world", 5, 5, "-"))
	assert.Equal(t, "Hello", Truncate("Hello", 10))
}
