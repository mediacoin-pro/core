package bin

import (
	"github.com/mediacoin-pro/core/common/json"
	"github.com/stretchr/testify/assert"
	"testing"
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
		json.Object(map[string]any{"a": 1, "b": 2}),
		json.Object(nil),
	)

	var (
		i  int
		s  string
		f  float64
		b  []byte
		p  Point
		ss []string
		jo json.Object
		j0 json.Object
	)
	err := Decode(buf, &i, &s, &f, &b, &p, &ss, &jo, &j0)

	assert.NoError(t, err)
	assert.Equal(t, 123, i)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 3.1415, f)
	assert.Equal(t, []byte{5, 6, 7}, b)
	assert.Equal(t, Point{88, 99}, p)
	assert.Equal(t, []string{"a", "b", "c"}, ss)
	assert.JSONEq(t, `{"a":1,"b":2}`, jo.String())
	assert.Equal(t, json.Object(nil), j0)
}
