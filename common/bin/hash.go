package bin

import (
	"crypto/sha256"
	"hash/fnv"
)

func Hash32(values ...interface{}) uint32 {
	h256 := Hash256(values...)
	return BytesToUint32(h256[:4])
}

func Hash64(values ...interface{}) uint64 {
	h256 := Hash256(values...)
	return BytesToUint64(h256)
}

func Hash128(values ...interface{}) []byte {
	h256 := Hash256(values...)
	return h256[:16]
}

func Hash160(values ...interface{}) []byte {
	h256 := Hash256(values...)
	return h256[:20]
}

func Hash256(values ...interface{}) []byte {
	hash := sha256.New()
	w := NewWriter(hash)
	for _, val := range values {
		w.WriteVar(val)
	}
	return hash.Sum(nil)
}

func FastHash64(values ...interface{}) uint64 {
	hash := fnv.New64()
	w := NewWriter(hash)
	for _, val := range values {
		w.WriteVar(val)
	}
	return hash.Sum64()
}
