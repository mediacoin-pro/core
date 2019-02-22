package chain

import (
	"encoding/json"

	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/hex"
	"github.com/mediacoin-pro/core/crypto"
)

type AddressInfo struct {
	Address  []byte
	Memo     uint64
	Balance  bignum.Int
	Asset    []byte
	LastTxID uint64
	LastTxTs int64 // time in µsec
	UserID   uint64
	UserNick string
}

func (i *AddressInfo) Encode() []byte {
	return bin.Encode(
		i.Address,
		i.Memo,
		i.Balance,
		i.Asset,
		i.LastTxID,
		i.LastTxTs,
		i.UserID,
		i.UserNick,
		0, 0, 0, // reserved
	)
}

func (i *AddressInfo) Decode(data []byte) error {
	return bin.Decode(data,
		&i.Address,
		&i.Memo,
		&i.Balance,
		&i.Asset,
		&i.LastTxID,
		&i.LastTxTs,
		&i.UserID,
		&i.UserNick,
	)
}

func (i *AddressInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address  string     `json:"address"`
		AddrMemo string     `json:"address_memo"`
		Memo     string     `json:"memo"`
		Balance  bignum.Int `json:"balance"`
		Asset    string     `json:"asset"`
		LastTxID string     `json:"last_tx_id"`
		LastTxTs int64      `json:"last_tx_ts"`
		UserID   string     `json:"user_id"`
		UserNick string     `json:"user_nick"`
	}{
		Address:  crypto.EncodeAddress(i.Address, 0),
		AddrMemo: crypto.EncodeAddress(i.Address, i.Memo),
		Memo:     encodeUint64(i.Memo),
		Balance:  i.Balance,
		Asset:    assets.Encode(i.Asset),
		LastTxID: encodeUint64(i.LastTxID),
		LastTxTs: i.LastTxTs,
		UserID:   encodeUint64(i.UserID),
		UserNick: i.UserNick,
	})
}

func encodeUint64(i uint64) string {
	if i == 0 {
		return ""
	}
	return "0x" + hex.EncodeUint(i)
}
