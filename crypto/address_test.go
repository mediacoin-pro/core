package crypto

import (
	"testing"

	"github.com/mediacoin-pro/core/common/rnd"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAddress(t *testing.T) {
	for i := 0; i < 10000; i++ {

		a := randBytes(20)
		sAddr := EncodeAddress(a)

		assert.Equal(t, 35, len(sAddr))
		assert.Equal(t, "MDC", sAddr[:3])

		addr, memo, err := DecodeAddress(sAddr)
		assert.NoError(t, err)
		assert.Equal(t, a, addr)
		assert.Equal(t, uint64(0), memo)
	}
}

func TestEncodeAddress_withMemo(t *testing.T) {
	for i := 0; i < 10000; i++ {

		a, m := randBytes(20), rnd.Uint64()
		sAddr := EncodeAddress(a, m)

		assert.True(t, len(sAddr) >= 35)
		assert.True(t, len(sAddr) <= 46)
		assert.Equal(t, "MDC", sAddr[:3])

		addr, memo, err := DecodeAddress(sAddr)
		assert.NoError(t, err)
		assert.Equal(t, a, addr)
		assert.Equal(t, m, memo)
	}
}

func TestDecodeAddress(t *testing.T) {
	sAddr := EncodeAddress([]byte("Hello, Qwerty-12345!"), 0)

	addr, memo, err := DecodeAddress(sAddr)

	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello, Qwerty-12345!"), addr)
	assert.Equal(t, uint64(0), memo)
}

func TestDecodeAddress_withMemo(t *testing.T) {
	sAddr := EncodeAddress([]byte("Hello, Qwerty-12345!"), 666)

	addr, memo, err := DecodeAddress(sAddr)

	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello, Qwerty-12345!"), addr)
	assert.Equal(t, uint64(666), memo)
}

func TestDecodeAddress_OldVersions(t *testing.T) {
	addr1, memo1, err1 := DecodeAddress("MTPRzisdxmzBoidNQUsbB1uUoqncuBsdLu") // v1
	addr2, memo2, err2 := DecodeAddress("ZXXXXypHGBtULioy94s9in55gyWvCMbkTR") // v0

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, addr1, addr2)
	assert.Equal(t, 20, len(addr1))
	assert.Equal(t, 20, len(addr2))
	assert.Equal(t, uint64(0), memo1)
	assert.Equal(t, uint64(0), memo2)
}

func TestDecodeAddress_Fail(t *testing.T) {

	addr, _, err := DecodeAddress("MDC7nQNHaA1Zn9FiSSZNbDMihwme9SUAvsy") // change last symbol

	assert.Error(t, err)
	assert.Nil(t, addr)
}

func TestPublicKey_Address(t *testing.T) {
	prv := NewPrivateKey()
	pub := prv.PublicKey()

	addr160 := pub.Address()
	sAddr := pub.StrAddress()

	assert.Equal(t, 160, len(addr160)*8)
	assert.Equal(t, 35, len(sAddr))
}
