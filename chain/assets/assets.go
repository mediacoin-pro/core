package assets

import (
	"bytes"
	"encoding/hex"
)

var (
	MDC  = []byte{0x01}
	AUTH = []byte{0x02} // Users auth-info (User`s public key)

	Default = MDC
)

const (
	NanoCoin  int64 = 1
	MicroCoin int64 = 1000
	MilliCoin int64 = 1000000
	Coin      int64 = 1000000000
	KiloCoin  int64 = 1000000000000
	MegaCoin  int64 = 1000000000000000
	GigaCoin  int64 = 1000000000000000000

	// Synonym of micro coin
	µCoin int64 = 1000
)

func Units(asset []byte) int64 {
	if IsMDC(asset) {
		return Coin
	}
	return 1
}

func IsMDC(typ []byte) bool {
	return len(typ) == 0 || bytes.Equal(typ, MDC)
}

func Encode(asset []byte) string {
	return hex.EncodeToString(asset)
}

func String(asset []byte) string {
	if IsMDC(asset) {
		return "MDC"
	}
	return hex.EncodeToString(asset)
}
