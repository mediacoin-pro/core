package enc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataSizeToString(t *testing.T) {
	assert.Equal(t, "0 B", BinarySizeToString(0))
	assert.Equal(t, "1 B", BinarySizeToString(1))
	assert.Equal(t, "123 B", BinarySizeToString(123))
	assert.Equal(t, "999 B", BinarySizeToString(999))
	assert.Equal(t, "0.97 KB", BinarySizeToString(1000))
	assert.Equal(t, "0.99 KB", BinarySizeToString(1023))
	assert.Equal(t, "2 MB", BinarySizeToString(2*1024*1024))
	assert.Equal(t, "2 MB", BinarySizeToString(2*1024*1024+100))
	assert.Equal(t, "2.01 MB", BinarySizeToString(2*1024*1024+10000))
	assert.Equal(t, "20.9 MB", BinarySizeToString(20*1024*1024+1000000))
	assert.Equal(t, "-20.9 MB", BinarySizeToString(-(20*1024*1024 + 1000000)))
}

func TestIP2Uint(t *testing.T) {
	assert.Equal(t, uint32(0x7f000001), IP2Uint("127.0.0.1"))
	assert.Equal(t, uint32(0xffffffff), IP2Uint("255.255.255.255"))
	assert.Equal(t, uint32(0x010203ff), IP2Uint("1.2.3.255"))

	assert.Equal(t, uint32(0), IP2Uint(""))
	assert.Equal(t, uint32(0), IP2Uint("256.0.0.1"))
	assert.Equal(t, uint32(0), IP2Uint("12.5.6.7.7"))
	assert.Equal(t, uint32(0), IP2Uint("127,0,0,1"))
	assert.Equal(t, uint32(0), IP2Uint("a.a.a.a"))
}
