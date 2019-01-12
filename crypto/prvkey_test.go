package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrivateKey(t *testing.T) {
	prv := NewPrivateKey()
	s64 := prv.String()

	prv2, err := ParsePrivateKey(s64)

	assert.NoError(t, err)
	assert.Equal(t, prv.Bytes(), prv2.Bytes())
	assert.Equal(t, prv.String(), prv2.String())
}
