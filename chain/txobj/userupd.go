package txobj

import (
	"fmt"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/common/json"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/model"
)

type UserUpd struct {
	Object
	Pubkey *crypto.PublicKey // new pub key

	reserved1 []byte
	reserved2 []byte
	reserved3 []byte
}

var _ = model.RegisterModel(model.TxUserUpd, &UserUpd{})

func NewUserUpd(
	bc chain.BCContext,
	sender *crypto.PublicKey,
	prv *crypto.PrivateKey,
	newKey *crypto.PublicKey,
) *chain.Transaction {
	return chain.NewTx(bc, sender, prv, 0, &UserUpd{
		Pubkey: newKey,
	})
}

func (u *UserUpd) String() string {
	return fmt.Sprintf("{UserUPD#%016x newPubKey:%s}", u.UserID(), u.Pubkey)
}

func (u *UserUpd) UserID() uint64 {
	return u.SenderID()
}

func (u *UserUpd) Encode() []byte {
	return bin.Encode(
		0, // version

		u.Pubkey,

		u.reserved1,
		u.reserved2,
		u.reserved3,
	)
}

func (u *UserUpd) Decode(data []byte) error {
	return bin.Decode(data,
		new(int), // version

		&u.Pubkey,

		&u.reserved1,
		&u.reserved2,
		&u.reserved3,
	)
}

func (u *UserUpd) Verify() error {
	if u.Pubkey.Empty() {
		return ErrTxEmptyParam
	}
	return nil
}

func (u *UserUpd) Execute(st *state.State) {
	st.SetBytes(assets.AUTH, u.SenderAddress(), u.Pubkey.Encode())
}

func (u *UserUpd) MarshalJSON() (data []byte, err error) {
	return json.Object{
		"id":     enc.UintToHex(u.UserID()),
		"pubKey": u.Pubkey,
	}.Bytes(), nil
}
