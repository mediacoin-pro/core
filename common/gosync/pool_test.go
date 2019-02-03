package gosync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoolInit(t *testing.T) {
	var p Pool

	assert.Equal(t, 0, p.Size())
	assert.Nil(t, p.Pop())
}

func TestPoolPush(t *testing.T) {
	var p Pool

	p.Push(123)
	p.Push("abc")

	assert.Equal(t, 2, p.Size())
}

func TestPoolValues(t *testing.T) {
	var p Pool
	p.Push(123)
	p.Push("abc")

	vv := p.Values()

	assert.Equal(t, 2, p.Size())
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, 123, vv[0].(int))
	assert.Equal(t, "abc", vv[1].(string))
}

func TestPoolPop(t *testing.T) {
	var p Pool
	p.Push("abc")
	p.Push(123)
	p.Push(456)

	v0 := p.Pop()
	v1 := p.Pop()

	assert.Equal(t, "abc", v0)
	assert.Equal(t, 123, v1)
	assert.Equal(t, 1, p.Size())
}

func TestPoolClear(t *testing.T) {
	var p Pool
	p.Push("abc")
	p.Push(123)
	p.Push(456)

	p.Clear()

	assert.Equal(t, 0, p.Size())
	assert.Nil(t, p.Pop())
}

func TestPoolString(t *testing.T) {
	var p Pool
	p.Push("abc")
	p.Push(123)
	p.Push(nil)

	s := p.String()

	assert.Equal(t, `["abc",123,null]`, s)
}
