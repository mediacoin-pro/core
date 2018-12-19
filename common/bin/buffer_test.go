package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferWriteRead(t *testing.T) {
	buf := NewBuffer(nil)
	buf.WriteBool(true)
	buf.WriteByte(123)
	buf.WriteString("abc")
	buf.WriteVar("ёпрст")
	buf.WriteVar(uint64(456))

	f, _ := buf.ReadBool()
	b, _ := buf.ReadByte()
	s, _ := buf.ReadString()

	var v string
	buf.ReadVar(&v)

	var num int
	buf.ReadVar(&num)

	assert.NoError(t, buf.Error())
	assert.Equal(t, true, f)
	assert.Equal(t, byte(123), b)
	assert.Equal(t, "abc", s)
	assert.Equal(t, "ёпрст", v)
	assert.Equal(t, 456, num)
}
