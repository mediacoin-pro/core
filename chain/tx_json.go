package chain

import (
	"encoding/json"

	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/hex"
	"github.com/mediacoin-pro/core/crypto"
)

type transactionJSON struct {
	TxID         hex.Uint64        `json:"id"`             //
	TxHash       hex.Bytes         `json:"hash"`           //
	BlockNum     uint64            `json:"block_num"`      //
	BlockIdx     int               `json:"block_idx"`      //
	BlockTs      int64             `json:"block_ts"`       //
	TxSeq        string            `json:"seq"`            //
	Type         int               `json:"type"`           // tx type
	TypeStr      string            `json:"stype"`          // tx type as string
	Version      int               `json:"version"`        // tx version
	Network      int               `json:"network"`        //
	ChainID      uint64            `json:"chain"`          //
	Nonce        uint64            `json:"nonce"`          //
	Sender       *crypto.PublicKey `json:"sender"`         // tx sender
	SenderAddr   string            `json:"sender_address"` // tx sender address
	SenderNick   string            `json:"sender_nick"`    // tx sender nickname (can be empty)
	ObjRaw       hex.Bytes         `json:"raw_data"`       // encoded tx-data
	Obj          ITransaction      `json:"obj"`            // unserialized data
	Sig          hex.Bytes         `json:"sig"`            //
	StateUpdates state.Values      `json:"state"`          //
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	if tx == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(&transactionJSON{
		Type:         tx.Type,
		TypeStr:      tx.StrType(),
		Version:      tx.Version,
		Network:      tx.Network,
		ChainID:      tx.ChainID,
		Nonce:        tx.Nonce,
		Sender:       tx.Sender,
		SenderAddr:   tx.SenderAddressStr(),
		SenderNick:   tx.SenderNick(),
		ObjRaw:       tx.Data,
		Obj:          tx.TxObject(),
		TxID:         hex.Uint64(tx.ID()),
		TxHash:       hex.Bytes(tx.Hash()),
		Sig:          hex.Bytes(tx.Sig),
		BlockNum:     tx.BlockNum(),
		BlockIdx:     tx.BlockIdx(),
		BlockTs:      tx.BlockTs(),
		TxSeq:        "0x" + hex.EncodeUint(tx.Seq()),
		StateUpdates: tx.StateUpdates,
	})
}
