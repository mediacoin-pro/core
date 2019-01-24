package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseObject(t *testing.T) {
	o, err := ParseObject([]byte(`{
		"i":	123,
		"num":	-12.3,
		"str":	"abc",
		"arr":	[1,2,3],
		"obj":	{"a":1},
		"objs":	[{"a":1},null,{"b":2}]
	}`))

	i := o.GetInt("i")
	num := o.GetNum("num")
	str := o.GetStr("str")
	arr := o.GetArr("arr")
	obj := o.GetObj("obj")
	objs := o.GetArr("objs").Objects()

	assert.NoError(t, err)
	assert.Equal(t, 123, i)
	assert.Equal(t, -12.3, num)
	assert.Equal(t, "abc", str)
	assert.Equal(t, `[1,2,3]`, arr.String())
	assert.Equal(t, `{"a":1}`, obj.String())
	assert.Equal(t, 3, len(objs))
	assert.Equal(t, `{"a":1}`, objs[0].String())
	assert.Equal(t, `null`, objs[1].String())
	assert.Equal(t, `{"b":2}`, objs[2].String())
}

func TestArray(t *testing.T) {
	v, err := Parse([]byte(`[null,1,2,3.3,"4",{"a":1}]`))

	ii := v.Array().Ints()
	ff := v.Array().Nums()
	ss := v.Array().Strings()
	oo := v.Array().Objects()

	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 0}, ii)
	assert.Equal(t, []float64{0, 1, 2, 3.3, 4, 0}, ff)
	assert.Equal(t, []string{"", "1", "2", "3.3", "4", `{"a":1}`}, ss)
	assert.Equal(t, `{"a":1}`, oo[5].String())
}

func TestValueToObject(t *testing.T) {
	var v = struct{ A, B int }{1, 2}

	obj, err := ValueToObject(v)

	assert.NoError(t, err)
	assert.Equal(t, 1, obj.GetInt("A"))
	assert.Equal(t, 2, obj.GetInt("B"))
}
