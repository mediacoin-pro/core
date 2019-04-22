package crypto

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/common/gosync"
	"github.com/mediacoin-pro/core/crypto/X15"
)

const PrivateKeyVersion = '\x01'

type PrivateKey struct {
	d   *big.Int
	pub *PublicKey
}

var errInvalidPrivateKey = errors.New("Invalid private key")

func NewPrivateKey() *PrivateKey {
	return generateKey(randInt())
}

func MustParsePrivateKey(prvKey64 string) (prv *PrivateKey) {
	prv, err := ParsePrivateKey(prvKey64)
	if err != nil {
		panic(err)
	}
	return
}

func ParsePrivateKey(prvKey64 string) (prv *PrivateKey, err error) {
	if prvKey64 == "" {
		return nil, errInvalidPrivateKey
	}
	bb, err := enc.Base64Decode(prvKey64)
	if err != nil {
		return
	}
	return decodePrivateKey(bb)
}

func NewPrivateKeyBySecret(secret string) *PrivateKey {
	key := X15.GenerateKeyByPassword([]byte(secret), KeySize*8)
	return generateKey(normInt(key))
}

func decodePrivateKey(b []byte) (*PrivateKey, error) {
	if len(b) < 1 || b[0] != PrivateKeyVersion {
		return nil, errInvalidPrivateKey
	}
	return generateKey(new(big.Int).SetBytes(b[1:])), nil
}

func (prv *PrivateKey) SubKey(subKeyName string) *PrivateKey {
	d := intToBytes(prv.d)
	secret := []byte{}
	secret = append(secret, d...)
	secret = append(secret, []byte(subKeyName)...)
	secret = append(secret, d...)
	return generateKey(hashInt(secret))
}

func (prv *PrivateKey) String() string {
	return enc.Base64Encode(prv.Bytes())
}

func (prv *PrivateKey) Bytes() []byte {
	if prv == nil {
		return nil
	}
	buf := []byte{PrivateKeyVersion} // head
	return append(buf, intToBytes(prv.d)...)
}

func (prv *PrivateKey) Secret64(tag string) uint64 {
	return bin.Hash64(prv.Secret(tag))
}

func (prv *PrivateKey) Secret(tag string) []byte {
	bb := prv.d.Bytes()
	bb = hash256(append(bb, []byte(tag)...))
	bb = hash256(append(bb, []byte(tag)...))
	bb = hash256(append(bb, []byte(tag)...))
	return bb
}

func (prv *PrivateKey) UserID() uint64 {
	return prv.PublicKey().ID()
}

func (prv *PrivateKey) PublicKey() *PublicKey {
	return prv.pub
}

// Sign signs a data using the private key, prv. It returns the signature as a
// pair of integers. The security of the private key depends on the entropy of
// rand.
func (prv *PrivateKey) Sign(data []byte) []byte {
	var k, s, r *big.Int
	e := hashInt(data)
	for {
		for {
			k = randInt()
			r, _ = curve.ScalarBaseMult(k.Bytes())
			r.Mod(r, curveParams.N)
			if r.Sign() != 0 {
				break
			}
		}
		s = new(big.Int).Mul(prv.d, r)
		s.Add(s, e)
		s.Mul(s, fermatInverse(k, curveParams.N))
		s.Mod(s, curveParams.N)
		if s.Sign() != 0 {
			break
		}
	}
	return append(intToBytes(r), intToBytes(s)...)
}

var cacheCipherKeys = gosync.NewCache(500)

func (prv *PrivateKey) sharedCipherKey(pub *PublicKey) []byte {
	if pub == nil {
		return hash256(hash256(prv.d.Bytes()))
	}
	pp := append(prv.Bytes(), pub.Bytes()...)
	if v, ok := cacheCipherKeys.Get(pp).([]byte); ok {
		return v
	}
	cipherKey := prv.calcSharedCipherKey(pub)
	cacheCipherKeys.Set(pp, cipherKey)
	return cipherKey
}

// Calculate shared key as hash of private key.
// It can be useful for RawEncryption
func (prv *PrivateKey) getPersonalCipherKey() []byte {
	bb := prv.d.Bytes()
	return hash256(append(bb, hash256(append(bb, hash256(bb)...))...))
}

func (prv *PrivateKey) calcSharedCipherKey(pub *PublicKey) []byte {
	s, _ := curve.ScalarMult(pub.x, pub.y, prv.d.Bytes())
	s.Mod(s, curveParams.N)
	return intToBytes(s)
}

func (prv *PrivateKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(prv.d)
}

func (prv *PrivateKey) UnmarshalJSON(data []byte) (err error) {
	prv.d = new(big.Int)
	if err = json.Unmarshal(data, prv.d); err == nil {
		prv.generatePub()
	}
	return
}

func (prv *PrivateKey) generatePub() {
	prv.pub = new(PublicKey)
	prv.pub.x, prv.pub.y = curve.ScalarBaseMult(prv.d.Bytes())
}

// generateKey generates a public and private key pair.
func generateKey(k *big.Int) *PrivateKey {
	prv := new(PrivateKey)
	prv.d = k
	prv.generatePub()
	return prv
}
