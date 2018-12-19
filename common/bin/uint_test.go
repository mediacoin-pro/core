package bin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64ToBytes(t *testing.T) {
	val := uint64(0xfedcba9876543210)

	b := Uint64ToBytes(val)

	assert.Equal(t, "fedcba9876543210", fmt.Sprintf("%x", b))
}

func TestUint64ToBytes_WithNulls(t *testing.T) {
	val := uint64(0xba9876543210)

	b := Uint64ToBytes(val)

	assert.Equal(t, "0000ba9876543210", fmt.Sprintf("%x", b))
}

func TestBytesToUint64(t *testing.T) {
	b, _ := hex.DecodeString("fedcba9876543210")

	val := BytesToUint64(b)

	assert.Equal(t, "fedcba9876543210", fmt.Sprintf("%x", val))
}

func TestBytesToUint64_WithNulls(t *testing.T) {
	b, _ := hex.DecodeString("ba9876543210")

	val := BytesToUint64(b)

	assert.Equal(t, "ba9876543210", fmt.Sprintf("%x", val))
}

func TestBytesToFloat32(t *testing.T) {
	org32 := float32(1. / 3)

	b := Float32ToBytes(org32)
	dec32 := BytesToFloat32(b)

	assert.Equal(t, 4, len(b))
	assert.Equal(t, org32, dec32)
}

func TestBytesToFloat64(t *testing.T) {
	org64 := float64(-1. / 13)

	b := Float64ToBytes(org64)
	dec64 := BytesToFloat64(b)

	assert.Equal(t, 8, len(b))
	assert.Equal(t, org64, dec64)
}
