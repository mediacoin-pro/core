package state

import (
	"bytes"
	"encoding/json"

	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/common/hex"
	"github.com/mediacoin-pro/core/crypto"
)

type Value struct {
	ChainID uint64
	Asset   []byte
	Address []byte
	Balance bignum.Int
	Memo    uint64
}

func (v *Value) String() string {
	return enc.JSON(v)
}

func (v *Value) Equal(b *Value) bool {
	return v.ChainID == b.ChainID &&
		bytes.Equal(v.Asset, b.Asset) &&
		bytes.Equal(v.Address, b.Address) &&
		v.Balance.Equal(b.Balance) &&
		v.Memo == b.Memo
}

func (v *Value) IsMDC() bool {
	return assets.IsMDC(v.Asset)
}

func (v *Value) StateKey() []byte {
	b := make([]byte, 0, 26)
	b = append(b, v.Address[:]...)
	b = append(b, v.Asset...)
	return b
}

func (v *Value) Hash() []byte {
	return bin.Hash256(
		v.ChainID,
		v.Asset,
		v.Address,
		v.Memo,
		v.Balance,
	)
}

func (v *Value) MarshalJSON() ([]byte, error) {
	var j = struct {
		ChainID uint64     `json:"chain"`
		Asset   string     `json:"asset"`
		Address string     `json:"address"`
		Raw     hex.Bytes  `json:"data"`
		Balance bignum.Int `json:"balance"`
		Memo    hex.Uint64 `json:"memo"`
	}{
		ChainID: v.ChainID,
		Asset:   assets.Encode(v.Asset),
		Address: crypto.EncodeAddress(v.Address),
		Raw:     v.Balance.Bytes(),
		Balance: v.Balance,
		Memo:    hex.Uint64(v.Memo),
	}
	if len(j.Raw) >= 64 {
		j.Balance = bignum.NewInt(0)
	}
	return json.Marshal(j)
}

func (v *Value) Encode() []byte {
	return bin.Encode(
		v.ChainID,
		v.Asset,
		v.Address,
		v.Memo,
		v.Balance,
	)
}

func (v *Value) Decode(data []byte) error {
	return bin.Decode(data,
		&v.ChainID,
		&v.Asset,
		&v.Address,
		&v.Memo,
		&v.Balance,
	)
}
