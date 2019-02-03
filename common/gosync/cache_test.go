package gosync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheSet(t *testing.T) {
	c := NewCache(5)

	for i := 0; i < 100; i++ {
		c.Set(i%13, 1)
		assert.True(t, c.Size() <= 5)
	}
	assert.Equal(t, 5, c.Size())
}

func TestCacheMultiSet(t *testing.T) {
	c := NewCache(100)

	c.Set(1, 0)
	c.Set("2", 0)
	c.Set([]byte{3}, 0)
	c.Set(1, 1)
	c.Set("2", 2)
	c.Set([]byte{3}, 3)

	assert.Equal(t, 3, c.Size())
	assert.Equal(t, 1, c.Get(1))
	assert.Equal(t, 2, c.Get("2"))
	assert.Equal(t, 3, c.Get([]byte{3}))
}

func TestCacheSetGet(t *testing.T) {
	c := NewCache(2)

	c.Set(1, 1)
	c.Set(2, 2)
	c.Set(3, 3)

	assert.Equal(t, 2, c.Size())
	assert.Equal(t, 3, c.Get(3))
	assert.False(t, c.Exists(1) && c.Exists(2))
	assert.True(t, c.Exists(1) || c.Exists(2))
}
