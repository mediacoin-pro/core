package gosync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackInit(t *testing.T) {

	var s Stack

	assert.Nil(t, s.Pop())
	assert.Equal(t, s.Size(), 0)
}

func TestStackClear(t *testing.T) {
	var s Stack
	s.Push(123)
	s.Push("abc")

	s.Clear()

	assert.Nil(t, s.Pop())
	assert.Equal(t, s.Size(), 0)
}

func TestStackPush(t *testing.T) {
	var s Stack

	s.Push(123)
	s.Push("abc")

	assert.Equal(t, s.Size(), 2)
}

func TestStackValues(t *testing.T) {
	var s Stack
	s.Push(123)
	s.Push("abc")

	vv := s.Values()

	assert.Equal(t, s.Size(), 2)
	assert.Equal(t, len(vv), 2)
	assert.Equal(t, vv[0].(int), 123)
	assert.Equal(t, vv[1].(string), "abc")
}

func TestStackPop(t *testing.T) {
	var s Stack
	s.Push(123)
	s.Push("abc")
	s.Push(456)

	v0 := s.Pop()
	v1 := s.Pop()

	assert.Equal(t, v0, 456)
	assert.Equal(t, v1, "abc")
	assert.Equal(t, s.Size(), 1)
}
