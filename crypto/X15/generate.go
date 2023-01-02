package X15

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/mediacoin-pro/core/crypto/bcrypt"
	"github.com/mediacoin-pro/core/crypto/sortcrypt"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

// redefined trace-function is using for benchmarks
var trace = func(string, []byte) {}

func GenerateKeyByPassword(hash []byte, keyLen int) []byte {
	key := generateKeyByPassword(hash)
	switch keyLen {
	case 256:
		hash256 := sha256.Sum256(key)
		return hash256[:]

	case 384:
		hash256 := sha3.Sum384(key)
		return hash256[:]

	case 512:
		hash256 := sha3.Sum512(key)
		return hash256[:]

	default:
		panic("GenerateKeyByPassword: Invalid key length")
	}
}

func generateKeyByPassword(hash []byte) []byte {

	// sha2-512
	h512 := sha512.Sum512(hash)
	hash = append(h512[:], hash...)
	trace("sha2-512", hash)

	// md5
	hashMD5 := md5.Sum(hash)
	hash = append(hashMD5[:], hash...)
	trace("md5", hash)

	// sha3-512
	h512 = sha3.Sum512(hash)
	hash = append(hash, h512[:]...)
	trace("sha3-512", hash)

	// scrypt
	if h, err := scrypt.Key(hash, h512[:], 8<<10, 8, 32, 256); err == nil {
		hash = append(h, hash...)
	} else {
		panic(err)
	}
	trace("scrypt", hash)

	// sha3-224
	h224 := sha3.Sum224(hash)
	hash = append(hash, h224[:]...)
	trace("sha3-224", hash)

	// elliptic curve 256
	x, y := elliptic.P256().ScalarBaseMult(hash[:256])
	hash = append(x.Bytes(), append(hash, y.Bytes()...)...)
	trace("curve-256", hash)

	// sha3-384
	h384 := sha3.Sum384(hash)
	hash = append(h384[:], hash...)
	trace("sha3-384", hash)

	// sortcrypt
	hSort := sortcrypt.GenerateKey(hash, 127, 4219)
	hash = append(hash, hSort...)
	trace("sortcrypt", hash)

	// bcrypt
	if bcrpt, err := bcrypt.GenerateFromPassword(bytes.NewBuffer(hash), hash, 12); err == nil {
		hash = append(bcrpt, hash...)
	} else {
		panic(err)
	}
	trace("bcrypt", hash)

	// elliptic curve 384
	x, y = elliptic.P384().ScalarBaseMult(hash[:384])
	hash = append(x.Bytes(), append(hash, y.Bytes()...)...)
	trace("curve-384", hash)

	// sha1
	h1 := sha1.Sum(hash)
	hash = append(h1[:], hash...)
	trace("sha1", hash)

	// ed25519
	if pub, _, err := ed25519.GenerateKey(bytes.NewBuffer(hash)); err == nil {
		hash = append(hash, []byte(pub)...)
	} else {
		panic(err)
	}
	trace("ed25519", hash)

	// sha2-256
	h256 := sha256.Sum256(hash)
	hash = append(h256[:], hash...)
	trace("sha2-256", hash)

	// rsa public key as hash
	rsaPub := rsaGeneratePublicKey(bytes.NewBuffer(hash))
	hash = append(rsaPub, hash...)
	trace("rsa", hash)

	// aes
	if bl, err := aes.NewCipher(h256[:32]); err == nil {
		cipher.NewCBCEncrypter(bl, h512[:16]).CryptBlocks(hash, hash[:len(hash)>>4<<4])
	} else {
		panic(err)
	}
	trace("aes", hash)

	return hash
}
