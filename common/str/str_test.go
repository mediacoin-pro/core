package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWords(t *testing.T) {
	s := " Hello, world !... "

	words := Words(s)

	assert.Equal(t, []string{"Hello", "world"}, words)
}

func TestToTitle(t *testing.T) {

	s := ToTitle("hello, world!")

	assert.Equal(t, "Hello, world!", s)
}
