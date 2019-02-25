package rnd

import (
	cryptoRand "crypto/rand"
	"math/rand"
	"reflect"

	"github.com/mediacoin-pro/core/common/bin"
)

func init() {
	rand.Seed(Int64()) // seed from crypto/rand.Reader
}

func Bytes(n int) []byte {
	buf := make([]byte, n)
	if _, err := cryptoRand.Read(buf); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return buf
}

func Uint64() uint64 {
	return bin.BytesToUint64(Bytes(8))
}

func Int64() int64 {
	return int64(Uint64())
}

func Intn(n int) int {
	return New().Intn(n)
}

func New(seeds ...interface{}) *rand.Rand {
	if len(seeds) == 0 {
		return rand.New(rand.NewSource(Int64()))
	}
	return rand.New(rand.NewSource(int64(bin.Hash64(seeds...))))
}

func Shuffle(slice interface{}, seed ...interface{}) {
	rnd := New(seed...)
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	n := rv.Len()
	for i := 0; i < n; i++ {
		swap(i, rnd.Intn(n))
	}
}
