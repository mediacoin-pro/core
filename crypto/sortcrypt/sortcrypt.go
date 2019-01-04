package sortcrypt

import (
	"bytes"
	"crypto/sha256"
	"sort"

	"golang.org/x/crypto/sha3"
)

func GenerateKey(pass []byte, rounds, size int) []byte {
	hash := pass
	ss := make(BytesSlice, size)
	for r := 0; r < rounds; r++ {
		for i := 0; i < size; i++ {
			h := sha256.Sum224(hash)
			hash = h[:]
			ss[i] = hash[:4]
		}
		sort.Sort(ss)
		h := sha3.Sum512(ss.Join())
		hash = h[:]
	}
	return hash
}

type BytesSlice [][]byte

func (bb BytesSlice) Join() []byte {
	res := []byte{}
	for _, b := range bb {
		res = append(res, b...)
	}
	return res
}

func (bb BytesSlice) Len() int {
	return len(bb)
}

func (bb BytesSlice) Less(i, j int) bool {
	return bytes.Compare(bb[i], bb[j]) < 0
}

func (bb BytesSlice) Swap(i, j int) {
	bb[i], bb[j] = bb[j], bb[i]
}
