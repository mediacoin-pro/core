package X15

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	key := GenerateKeyByPassword([]byte("secret-string"), 256)

	assert.Equal(t, 256/8, len(key))
	assert.Equal(t, "1905fb44a8e19bba89afba7f8a13a5218e5af4b93ee81d895c5cdab93a0e8bd9", hex.EncodeToString(key))
}
