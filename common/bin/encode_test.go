package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode_Nil(t *testing.T) {
	var v interface{}

	buf := Encode(v)

	assert.Equal(t, []byte{0}, buf)
}

func TestDecode(t *testing.T) {
	buf := Encode(
		uint64(123),
		"abc",
		3.1415,
		[]byte{5, 6, 7},
		Point{88, 99},
		[]string{"a", "b", "c"},
	)

	var (
		i  int
		s  string
		f  float64
		b  []byte
		p  Point
		ss []string
	)
	err := Decode(buf, &i, &s, &f, &b, &p, &ss)

	assert.NoError(t, err)
	assert.Equal(t, 123, i)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 3.1415, f)
	assert.Equal(t, []byte{5, 6, 7}, b)
	assert.Equal(t, Point{88, 99}, p)
	assert.Equal(t, []string{"a", "b", "c"}, ss)
}
