package bin

import (
	"encoding/hex"
	"encoding/json"
)

type Bytes []byte

func (b Bytes) String() string {
	return hex.EncodeToString(b)
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	return []byte(`"` + b.String() + `"`), nil
}

func (b *Bytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b, err = hex.DecodeString(s)
	return
}
