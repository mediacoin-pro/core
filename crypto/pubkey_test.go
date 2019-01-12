package crypto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubKeyEncodeParse(t *testing.T) {
	prv := NewPrivateKey()
	pub := prv.PublicKey()

	buf := pub.Bytes()
	pub1, err := decodePublicKey(buf)

	assert.Equal(t, PublicKeySize, len(buf))
	assert.Equal(t, pub, pub1)
	assert.NoError(t, err)
}

func TestPubKeyJSONMarshal(t *testing.T) {
	prv := NewPrivateKey()
	org := prv.PublicKey()
	buf, errEnc := json.Marshal(org)

	var dec = new(PublicKey)
	errDec := json.Unmarshal(buf, dec)

	assert.NoError(t, errEnc)
	assert.NoError(t, errDec)
	assert.True(t, org.Equal(dec))
}

func TestPubKeyJSONMarshalObject(t *testing.T) {
	var org, dec struct {
		Key *PublicKey
	}
	prv := NewPrivateKey()
	org.Key = prv.PublicKey()
	buf, errEnc := json.Marshal(&org)

	errDec := json.Unmarshal(buf, &dec)

	assert.NoError(t, errEnc)
	assert.NoError(t, errDec)
	assert.True(t, org.Key.Equal(dec.Key))
}
