package enc

import (
	"encoding/base64"
	"strings"
)

func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

var base64StdToURLReplacer = strings.NewReplacer(
	"+", "-",
	"/", "_",
)

func Base64Decode(str64 string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(base64StdToURLReplacer.Replace(str64))
}
