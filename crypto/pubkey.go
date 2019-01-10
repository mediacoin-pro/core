package crypto

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/mediacoin-pro/core/common/enc"

	"golang.org/x/crypto/ripemd160"
)

type PublicKey struct {
	x *big.Int
	y *big.Int
}

// Verify verifies the signature in r, s of hash using the public key, pub. Its
// return value records whether the signature is valid.
func (pub *PublicKey) Verify(data []byte, sig []byte) bool {
	if pub.Empty() {
		return false
	}
	if len(sig) != PublicKeySize {
		return false
	}
	r := new(big.Int).SetBytes(sig[:KeySize])
	s := new(big.Int).SetBytes(sig[KeySize:])

	if r.Sign() == 0 || r.Cmp(curveParams.N) >= 0 {
		return false
	}
	if s.Sign() == 0 || s.Cmp(curveParams.N) >= 0 {
		return false
	}

	e := hashInt(data)
	w := new(big.Int).ModInverse(s, curveParams.N)

	u1 := e.Mul(e, w)
	u2 := w.Mul(r, w)

	u1.Mod(u1, curveParams.N)
	u2.Mod(u2, curveParams.N)

	x1, y1 := curve.ScalarBaseMult(u1.Bytes())
	x2, y2 := curve.ScalarMult(pub.x, pub.y, u2.Bytes())
	x, y := curve.Add(x1, y1, x2, y2)
	if x.Sign() == 0 && y.Sign() == 0 {
		return false
	}
	x.Mod(x, curveParams.N)
	return x.Cmp(r) == 0
}

func (pub *PublicKey) Empty() bool {
	return pub == nil || pub.x == nil && pub.y == nil
}

func (pub *PublicKey) String() string {
	return "0x04" + hex.EncodeToString(pub.Encode())
}

func (pub *PublicKey) HexID() string {
	return enc.UintToHex(pub.ID())
}

func (pub *PublicKey) ID() uint64 {
	return AddressToUserID(pub.Address())
}

func (pub *PublicKey) Equal(p *PublicKey) bool {
	return pub != nil && p != nil && pub.x.Cmp(p.x) == 0 && pub.y.Cmp(p.y) == 0
}

func (pub *PublicKey) Bytes() []byte {
	return pub.Encode()
}

func (pub *PublicKey) StrAddress() string {
	return EncodeAddress(pub.Address())
}

func (pub *PublicKey) Address() []byte {
	hash256 := newHash256()
	hash256.Write(intToBytes(pub.x))
	hash256.Write(intToBytes(pub.y))

	hash160 := ripemd160.New()
	hash160.Write(hash256.Sum(nil))
	return hash160.Sum(nil)
}

func (pub *PublicKey) Encode() []byte {
	a := make([]byte, 0, KeySize+KeySize)
	a = append(a, intToBytes(pub.x)...) // 32 bytes X
	a = append(a, intToBytes(pub.y)...) // 32 bytes Y
	return a
}

func (pub *PublicKey) Decode(data []byte) error {
	if len(data) == 2*KeySize+1 && data[0] == 04 {
		data = data[1:]
	}
	if len(data) != 2*KeySize {
		return errors.New("crypto.PublicKey.Decode-error")
	}
	pub.x = new(big.Int).SetBytes(data[:KeySize])
	pub.y = new(big.Int).SetBytes(data[KeySize:])
	return nil
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pub.String() + `"`), nil
}

func (pub *PublicKey) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if p, err := ParsePublicKey(str); err != nil {
		return err
	} else {
		pub.x = p.x
		pub.y = p.y
		return nil
	}
}

func MustParsePublicKey(pubkey string) *PublicKey {
	pub, err := ParsePublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	return pub
}

func ParsePublicKey(s string) (pub *PublicKey, err error) {
	var data []byte
	if strings.HasPrefix(s, "0x") {
		data, err = hex.DecodeString(s[2:])
	} else {
		data, err = enc.Base64Decode(s) // deprecated
	}
	if err != nil {
		return
	}
	return decodePublicKey(data)
}

func decodePublicKey(data []byte) (*PublicKey, error) {
	pub := new(PublicKey)
	if err := pub.Decode(data); err != nil {
		return nil, errors.New("crypto: Unknown format of public-key. Err: " + err.Error())
	}
	return pub, nil
}
