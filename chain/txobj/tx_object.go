package txobj

import (
	"errors"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/crypto"
)

var (
	ErrTxIncorrectAmount  = errors.New("tx-Error: Incorrect amount")
	ErrTxIncorrectSender  = errors.New("tx-Error: Incorrect sender")
	ErrTxIncorrectAsset   = errors.New("tx-Error: Incorrect asset")
	ErrTxIncorrectAddress = errors.New("tx-Error: Incorrect address")
	ErrTxIncorrectValue   = errors.New("tx-Error: Incorrect value")
	ErrTxIncorrectNick    = errors.New("tx-Error: Incorrect nick")
	ErrTxLongComment      = errors.New("tx-Error: Comment is too long")
	ErrTxEmptyOuts        = errors.New("tx-Error: Empty outputs")
	ErrTxEmptyParam       = errors.New("tx-Error: Empty param")
)

type Object struct {
	tx *chain.Transaction `json:"-"`
}

func (obj *Object) Tx() *chain.Transaction {
	return obj.tx
}

func (obj *Object) Sender() *crypto.PublicKey {
	if obj.tx != nil {
		return obj.tx.Sender
	}
	return nil
}
func (obj *Object) SenderAddress() []byte {
	if obj != nil && obj.tx != nil {
		return obj.tx.SenderAddress()
	}
	return nil
}

func (obj *Object) SenderAddressStr() string {
	return crypto.EncodeAddress(obj.SenderAddress())
}

func (obj *Object) SenderID() uint64 {
	if sender := obj.Sender(); sender != nil {
		return sender.ID()
	}
	return 0
}

func (obj *Object) ChainID() uint64 {
	return obj.ChainConfig().ChainID
}

func (obj *Object) NetworkID() int {
	return obj.ChainConfig().NetworkID
}

func (obj *Object) ChainConfig() *chain.Config {
	return obj.tx.BCContext().Config()
}

func (obj *Object) SetContext(tx *chain.Transaction) {
	obj.tx = tx
}

func (obj *Object) UserNickByID(userID uint64) (nick string, err error) {
	if obj != nil && obj.tx != nil {
		return obj.tx.UsernameByID(userID)
	}
	return
}

func (obj *Object) UserNickByAddress(addr []byte) (nick string, err error) {
	return obj.UserNickByID(crypto.AddressToUserID(addr))
}
