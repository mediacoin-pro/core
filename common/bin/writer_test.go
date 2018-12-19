package bin

import (
	"bytes"
	"testing"

	"math/big"

	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {

	w := NewWriter(bytes.NewBuffer(nil))

	assert.True(t, w != nil)
}

func TestWriter_WriteVar(t *testing.T) {
	w := NewBuffer(nil)

	w.WriteVar(0)
	w.WriteVar(13)
	w.WriteVar(255)
	w.WriteVar(256)
	w.WriteVar(-13)
	w.WriteVar(0x01020304050607)
	var max64 uint64 = 0xffffffffffffffff
	w.WriteVar(max64)
	w.WriteVar(0.3)
	w.WriteVar([]int{77, 88, 99})
	w.WriteVar("abc")
	w.WriteVar(map[int]int{4: 55})

	assert.Equal(t, []byte{
		0,          // 0
		13,         // 13
		0x81, 0xff, // 255
		0x82, 1, 0, // 256
		0xc1, 13, // -13
		0x87, 1, 2, 3, 4, 5, 6, 7, // 0x01020304050607
		0xc1, 1, // -1
		0x3f, 0xd3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, // 0.3
		0x3, 77, 88, 99, // []int{77, 88, 99}
		0x3, 0x61, 0x62, 0x63, // "abc"
		0x1, 4, 55,
	}, w.Bytes())
}

func TestWriter_WriteBigInt(t *testing.T) {
	w := NewBuffer(nil)

	w.WriteVar(big.NewInt(0))
	w.WriteVar(big.NewInt(13))
	w.WriteVar(big.NewInt(255))
	w.WriteVar(big.NewInt(256))
	w.WriteVar(big.NewInt(-13))
	w.WriteVar(big.NewInt(0x01020304050607))
	w.WriteVar(big.NewInt(-0x01020304050607))
	w.WriteVar(newBigInt("ffeeddccbbaa99887766554433221100ffeeddccbbaa99887766554433221100")) // 256 bit

	assert.Equal(t, []byte{
		0,          // 0
		13,         // 13
		0x81, 0xff, // 255
		0x82, 1, 0, // 256
		0xc1, 13, // -13
		0x87, 1, 2, 3, 4, 5, 6, 7, // 0x01020304050607
		0xc7, 1, 2, 3, 4, 5, 6, 7, // -0x01020304050607
		0xa0, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
	}, w.Bytes())
}

func TestWriter_WriteString(t *testing.T) {
	w := NewBuffer(nil)

	w.WriteVar("")
	w.WriteVar("Abc")

	assert.Equal(t, []byte{0, 3, 'A', 'b', 'c'}, w.Bytes())
}

func newBigInt(hex string) *big.Int {
	i, _ := big.NewInt(0).SetString(hex, 16)
	return i
}
