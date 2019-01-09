package crypto

import (
	"bytes"
	"errors"
	"strings"

	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/crypto/base58"
)

const (
	addressPrefix  = "MDC"
	addressVersion = 1
)

func addressCheckSum(b []byte) []byte {
	return hash256(hash256(append([]byte(addressPrefix), b...)))[:3]
}

func trimLeft0(b []byte) []byte {
	for len(b) > 0 && b[0] == 0 {
		b = b[1:]
	}
	return b
}

func AddressToUserID(addr160 []byte) uint64 {
	if IsValidAddress(addr160) {
		return bin.BytesToUint64(addr160)
	}
	return 0
}

func EncodeAddress(addr160 []byte, memo ...uint64) string {

	var mem uint64
	if len(memo) > 0 {
		mem = memo[0]
	}
	if len(addr160) == 0 && mem == 0 {
		return ""
	}
	if len(addr160) != 20 { // error: invalid address-length
		panic(errEncodeAddrInvalid)
	}
	key := make([]byte, 0, 32)
	key = append(key, addressVersion)                       // 1 byte
	key = append(key, addr160...)                           // 20 bytes
	key = append(key, trimLeft0(bin.Uint64ToBytes(mem))...) // ≤8 bytes
	key = append(key, addressCheckSum(key)...)              // 3 bytes
	return addressPrefix + base58.Encode(key)
}

var (
	errEncodeAddrInvalid    = errors.New("crypto.EncodeAddress: Invalid address")
	errDecodeAddrInvalid    = errors.New("crypto.DecodeAddress: Invalid address")
	errDecodeAddrUnknownVer = errors.New("crypto.DecodeAddress: Unknown address version")
	errDecodeAddrInvalidSum = errors.New("crypto.DecodeAddress: Invalid check-sum")
)

func DecodeAddress(strAddr string) (addr []byte, memo uint64, err error) {

	if len(strAddr) == 34 { // old format
		addr, err = decodeOldVersionAddress(strAddr)
		return
	}
	if len(strAddr) < 35 || !strings.HasPrefix(strAddr, addressPrefix) {
		err = errDecodeAddrInvalid
		return
	}
	bb, err := base58.Decode(strings.TrimPrefix(strAddr, addressPrefix))
	if err != nil {
		return
	}
	if len(bb) < 1+20+3 {
		err = errDecodeAddrInvalid
		return
	}
	if ver := bb[0]; ver != addressVersion {
		err = errDecodeAddrUnknownVer
		return
	}
	// check sum
	var data, sum = bb[:len(bb)-3], bb[len(bb)-3:]
	if !bytes.Equal(sum, addressCheckSum(data)) {
		err = errDecodeAddrInvalidSum
		return
	}

	// ok. valid address
	addr = data[1:21]
	memo = bin.BytesToUint64(data[21:])
	return
}

func decodeOldVersionAddress(strAddr string) (addr []byte, err error) {
	bb, err := base58.Decode(strAddr)
	if err != nil || len(bb) != 25 {
		return nil, errDecodeAddrInvalid
	}
	if ver := bb[0]; ver != '\x50' && ver != '\x32' {
		return nil, errDecodeAddrUnknownVer
	}
	if chSum := hash256(hash256(bb[:21]))[:4]; !bytes.Equal(chSum, bb[21:25]) {
		return nil, errDecodeAddrInvalidSum
	}
	// ok. valid address
	return bb[1:21], nil
}

func MustParseAddress(addr string) []byte {
	a, _, err := DecodeAddress(addr)
	if err != nil {
		panic(err)
	}
	return a
}

func IsValidAddress(addr160 []byte) bool {
	return len(addr160) == 20
}

func DecodeAsset(strAsset string) (asset []byte, err error) {
	if strAsset == "" {
		asset = assets.MDC
		return
	}
	return enc.Base64Decode(strAsset)
}

func DecodeAddressAsset(strAddr, strAsset string) (addr, asset []byte, err error) {
	if addr, _, err = DecodeAddress(strAddr); err == nil {
		asset, err = DecodeAsset(strAsset)
	}
	return
}
