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
