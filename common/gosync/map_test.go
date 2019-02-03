package gosync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapInit(t *testing.T) {
	var m Map

	assert.Nil(t, m.Get(123))
	assert.Equal(t, 0, m.Size())
}

func TestMapSetGet(t *testing.T) {
	var m Map

	m.Set("abc", 123)
	m.Set(456, "def")
	m.Set([]byte("XYZ"), 789)

	assert.Nil(t, m.Get(1))
	assert.False(t, m.Exists(1))
	assert.True(t, m.Exists(456))
	assert.True(t, m.Exists("abc"))
	assert.Equal(t, 123, m.Get("abc"))
	assert.NotEqual(t, 123., m.Get("abc"))
	assert.Equal(t, "def", m.Get(456))
	assert.Equal(t, 789, m.Get([]byte("XYZ")))
	assert.Equal(t, 3, m.Size())
}

func TestMapDelete(t *testing.T) {
	var m Map
	m.Set(0, 1)
	m.Set("abc", 123)
	m.Set(456, "def")

	m.Delete("abc")

	assert.False(t, m.Exists("abc"))
	assert.Nil(t, m.Get("abc"))
	assert.Equal(t, 2, m.Size())
}

func TestMapValues(t *testing.T) {
	var m Map
	m.Set("abc", 123)
	m.Set("def", 456)

	vv := m.Values()

	assert.Equal(t, 2, m.Size())
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, 123+456, vv[0].(int)+vv[1].(int))
}

func TestMapString(t *testing.T) {
	var m Map
	m.Set("abc", 123)
	m.Set(456, "def")

	s := m.String()

	assert.Equal(t, `{"456":"def","abc":"123"}`, s)
}
