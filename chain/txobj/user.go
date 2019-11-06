package txobj

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/common/json"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/crypto/base58"
	"github.com/mediacoin-pro/core/model"
)

type User struct {
	Object
	Nick       string // login
	ReferrerID uint64 // referrer address 160
}

var _ = model.RegisterModel(model.TxUser, &User{})

var reNick = regexp.MustCompile(`^[a-z][a-z0-9_\-]{2,20}$`)

func NewUser(
	bc chain.BCContext,
	sender *crypto.PrivateKey,
	nick string,
	referrerID uint64,
) *chain.Transaction {
	return chain.NewTx(bc, nil, sender, 0, &User{
		Nick:       nick,
		ReferrerID: referrerID,
	})
}

func (u *User) String() string {
	return fmt.Sprintf("{USER#%016x @%s}", u.UserID(), u.Nick)
}

func (u *User) UserID() uint64 {
	return u.SenderID()
}

func (u *User) InviteCode() string {
	return base58.Itoa(u.UserID())
}

func (u *User) Registered() time.Time {
	return u.Tx().Timestamp()
}

func (u *User) PublicKey() *crypto.PublicKey {
	return u.Sender()
}

func (u *User) Address() []byte {
	return u.SenderAddress()
}

//func (u *User) NewDoc(sCID string, data json.Object) *dsobj.DocumentObject {
//	cid, _ := dsobj.ParseChannelID(sCID)
//	return dsobj.NewDocumentObject(u.UserID(), cid, data)
//}

func (u *User) Encode() []byte {
	return bin.Encode(
		0, // version

		u.Nick,
		u.ReferrerID,
	)
}

func (u *User) Decode(data []byte) error {
	return bin.Decode(data,
		new(int), // version

		&u.Nick,
		&u.ReferrerID,
	)
}

func (u *User) ChannelID() []byte {
	return ChannelID(u.Sender())
}

func (u *User) StrChannelID() string {
	return EncodeChannelID(u.ChannelID())
}

func ChannelID(pub *crypto.PublicKey) []byte {
	if pub != nil {
		bb := pub.Bytes()
		return bb[:16]
	}
	return nil
}

func EncodeChannelID(cid []byte) (s string) {
	if len(cid) > 16 {
		return base58.Encode(cid[:16]) + "-" + base58.Encode(cid[16:])
	} else {
		return base58.Encode(cid)
	}
}

func (u *User) Verify() error {
	if !reNick.MatchString(u.Nick) {
		return ErrTxIncorrectNick
	}
	return nil
}

func (u *User) Execute(st *state.State) {
	// do noting
}

func (u *User) MarshalJSON() (data []byte, err error) {
	return json.Object{
		"id":       enc.UintToHex(u.UserID()),
		"nick":     u.Nick,
		"referrer": enc.UintToHex(u.ReferrerID),
		"address":  u.SenderAddressStr(),
		"cid":      u.StrChannelID(), // channelID of user`s PLs
	}.Bytes(), nil
}
