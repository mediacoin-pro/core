package txobj

import (
	"bytes"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/hex"
	"github.com/mediacoin-pro/core/common/json"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/model"

	"github.com/mediacoin-pro/core/common/bin"

	"github.com/mediacoin-pro/core/common/bignum"
)

// Transfer
type Transfer struct {
	Object
	Outs    []*TransferOutput //
	Comment []byte            // sender encrypted comment
}

var _ = model.RegisterModel(model.TxTransfer, &Transfer{})

var subscriptionRewardAddress = crypto.MustParseAddress("MDCB4B4kLg8f2W31a27QrbPLrc7c8nPMkyQ")

type TransferOutput struct {
	Asset     []byte     //
	Amount    bignum.Int //
	Tag       uint64     // sender memo
	To        []byte     //
	ToMemo    uint64     //
	ToChainID uint64     //
	Comment   []byte     // recipient encrypted comment
	Reserved1 []byte     //
	Reserved2 []byte     //

	// not imported
	decryptedComment string
	tr               *Transfer
}

func NewTransfer(
	bc chain.BCContext,
	sender *crypto.PublicKey,
	prv *crypto.PrivateKey,
	outs []*TransferOutput,
	comment string,
	nonce uint64,
) *chain.Transaction {

	//encComment := sender.EncryptRaw(bin.Uint64ToBytes(nonce), []byte(comment), nil)
	encComment := []byte(comment)

	tr := &Transfer{
		Outs:    outs,
		Comment: encComment,
	}
	defer tr.initOutputsContext()

	return chain.NewTx(bc, sender, prv, nonce, tr)
}

func NewSimpleTransfer(
	bc chain.BCContext,
	sender *crypto.PublicKey,
	prv *crypto.PrivateKey,
	asset []byte,
	amount bignum.Int,
	tag uint64, // sender memo
	toAddress []byte, //
	toMemo uint64, //
	comment string,
	nonce uint64,
) *chain.Transaction {
	var toChainID = chain.DefaultConfig.ChainID
	if bc != nil {
		toChainID = bc.Config().ChainID
	}
	if asset == nil {
		asset = assets.Default
	}
	return NewTransfer(bc, sender, prv, []*TransferOutput{{
		Asset:     asset,
		Amount:    amount,
		Tag:       tag,
		To:        toAddress,
		ToMemo:    toMemo,
		ToChainID: toChainID,
	}}, comment, nonce)
}

func (tr *Transfer) TotalAmount() (s bignum.Int) {
	for _, out := range tr.Outs {
		s.Increment(out.Amount)
	}
	return
}

func (tr *Transfer) Encode() []byte {
	return bin.Encode(
		0, // ver
		tr.Comment,
		tr.Outs,
	)
}

func (tr *Transfer) Decode(data []byte) error {
	defer tr.initOutputsContext()

	return bin.Decode(data,
		new(int),
		&tr.Comment,
		&tr.Outs,
	)
}

func (tr *Transfer) IsDonation() bool {
	return bytes.HasPrefix(tr.Comment, []byte("Donation for"))
}

func (tr *Transfer) IsSubscriptionRewards() bool {
	return bytes.Equal(tr.SenderAddress(), subscriptionRewardAddress)
}

func (tr *Transfer) initOutputsContext() {
	for _, out := range tr.Outs {
		out.tr = tr
	}
}

func (out *TransferOutput) Encode() []byte {
	return bin.Encode(
		out.Asset,
		out.Amount,
		out.Tag,
		out.To,
		out.ToMemo,
		out.ToChainID,
		out.Comment,
		out.Reserved1,
		out.Reserved2,
	)
}

func (out *TransferOutput) Decode(data []byte) error {
	return bin.Decode(data,
		&out.Asset,
		&out.Amount,
		&out.Tag,
		&out.To,
		&out.ToMemo,
		&out.ToChainID,
		&out.Comment,
		&out.Reserved1,
		&out.Reserved2,
	)
}

func (tr *Transfer) Verify() error {

	if len(tr.Outs) == 0 {
		return ErrTxEmptyOuts
	}
	if len([]byte(tr.Comment)) > 200 {
		return ErrTxLongComment
	}

	// check values; check sum of In and Out
	for _, out := range tr.Outs {
		if out.Amount.Sign() <= 0 {
			return ErrTxIncorrectAmount
		}
		if !crypto.IsValidAddress(out.To) {
			return ErrTxIncorrectAddress
		}
	}

	return nil
}

func (tr *Transfer) Execute(st *state.State) {

	senderAddr := tr.SenderAddress()

	//st.Decrement(assets.MDC, senderAddr, t.Fee(), 0) // todo: ??? fee

	for _, out := range tr.Outs {
		st.Decrement(out.Asset, senderAddr, out.Amount, out.Tag)

		// increment coins on address
		if out.ToChainID == tr.tx.ChainID {
			st.Increment(out.Asset, out.To, out.Amount, out.ToMemo)
		} else {
			st.CrossChainSet(out.ToChainID, out.Asset, out.To, out.Amount, out.ToMemo)
		}
	}
}

func (tr *Transfer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Outs       []*TransferOutput `json:"outs"`        //
		RawComment []byte            `json:"raw_comment"` //
		Comment    string            `json:"comment"`     //
	}{
		Outs:       tr.Outs,
		RawComment: tr.Comment,
		Comment:    string(tr.Comment), // todo: decoded comment
	})
}

func (out *TransferOutput) ToNick() string {
	if out == nil || len(out.To) == 0 {
		return ""
	}
	return chain.UserNameByID(crypto.AddressToUserID(out.To))
}

func (out *TransferOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Asset      string     `json:"asset"`       //
		Amount     bignum.Int `json:"amount"`      //
		Tag        string     `json:"tag"`         //
		To         string     `json:"to"`          //
		ToMemo     string     `json:"to_memo"`     //
		ToChainID  uint64     `json:"to_chain_id"` //
		ToNick     string     `json:"to_nick"`     //
		RawComment []byte     `json:"raw_comment"` //
		//Comment string  `json:"comment"` // todo: decoded comment
	}{
		Asset:      hex.Encode(out.Asset),
		Amount:     out.Amount,
		Tag:        hex.Encode(out.Tag),
		To:         crypto.EncodeAddress(out.To),
		ToMemo:     crypto.EncodeAddress(out.To, out.ToMemo),
		ToChainID:  out.ToChainID,
		ToNick:     out.ToNick(),
		RawComment: out.Comment,
	})
}

//func (t *Transfer) DecryptData(prv *crypto.PrivateKey) {
//	if prv == nil {
//		return
//	}
//	pub := prv.PublicKey()
//	if !pub.Equal(t.Sender()) {
//		return
//	}
//
//	tx := t.Tx()
//	if tx == nil {
//		return
//	}
//	for _, out := range t.Outs {
//		prv.DecryptRaw(bin.Uint64ToBytes(tx.Nonce))
//	}
//
//}
