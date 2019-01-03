package merkle

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	hashes := newHashes(15)

	root := Root(hashes...)

	assert.Equal(t, "5b1dc1a8872f60d0d8e80eadbbd05d95217835419c3f75e7a8366f6556fd3784", Hex(root))
}

func TestProof(t *testing.T) {
	hashes := newHashes(15)

	proof, root := Proof(hashes, 4)

	assert.Equal(t, ""+
		"004f5aa6aec3fc78c6aae081ac8120c720efcd6cea84b6925e607be063716f96dd"+
		"00b1ade06004f311750384e54a2d00dc3935f08434ad35461a59bba762cdf4f7a1"+
		"017839fa3a98157c652dd8b67fd61d3956600b30df17a4c5e1f811bc07dd1f63c1"+
		"00a2eebf45b258ac5b60516b3316d469cb095e7e41a73342877818b79e9595885e",
		Hex(proof),
	)
	assert.Equal(t, "5b1dc1a8872f60d0d8e80eadbbd05d95217835419c3f75e7a8366f6556fd3784", Hex(root))
}

func TestVerify(t *testing.T) {
	for i := 0; i < 100; i++ {
		hashes := newHashes(100)
		hash := hashes[i]
		proof, root := Proof(hashes, i)

		ok := Verify(hash, proof, root)

		assert.True(t, ok)
	}
}

func TestVerify_fail(t *testing.T) {
	hashes := newHashes(100)
	proof, root := Proof(hashes, 13)

	proof = proof[:len(proof)-1] // corrupt proof (cut last byte)
	ok := Verify(hashes[13], proof, root)

	assert.False(t, ok)
}

//-------------------------------------------------------------------
func newHashes(n int) (data [][]byte) {
	r := rand.New(rand.NewSource(0))
	for i := 0; i < n; i++ {
		buf := make([]byte, HashSize)
		r.Read(buf)
		data = append(data, buf)
	}
	return
}

func Hex(b []byte) string {
	return hex.EncodeToString(b)
}
