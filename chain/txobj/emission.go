package txobj

import (
	"bytes"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/consts"
	"github.com/mediacoin-pro/core/common/hex"
	"github.com/mediacoin-pro/core/common/json"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/model"

	"github.com/mediacoin-pro/core/common/bin"

	"github.com/mediacoin-pro/core/common/bignum"
)

type Emission struct {
	Object
	Asset   []byte            `json:"asset"`   // coin
	Comment string            `json:"comment"` //
	Outs    []*EmissionOutput `json:"outs"`    //
}

type EmissionOutput struct {
	Type      int        `json:"type"`    //
	Address   []byte     `json:"address"` // address associated with the media-source
	Value     int64      `json:"value"`   //
	Amount    bignum.Int `json:"amount"`  //
	Reserved1 []byte     `json:"-"`       //
	Reserved2 []byte     `json:"-"`       //
}

var _ = chain.RegisterTxType(model.TxEmission, &Emission{})

const (
	EmissionTypeInit         = 0 //
	EmissionTypeDistribution = 1 //
	EmissionTypeAuthor       = 2 //
	EmissionTypeReferer      = 3 //
)

var emissionTypeStr = map[int]string{
	EmissionTypeInit:         "Initial",
	EmissionTypeDistribution: "Distributors reward",
	EmissionTypeAuthor:       "Authors reward",
	EmissionTypeReferer:      "Referrals reward",
}

func NewEmission(
	bc chain.BCContext,
	emissionKey *crypto.PrivateKey,
	asset []byte,
	comment string,
	vv []*EmissionOutput,
) *chain.Transaction {
	return chain.NewTx(bc, emissionKey, 0, &Emission{
		Asset:   asset,
		Comment: comment,
		Outs:    vv,
	})
}

func (obj *Emission) Encode() []byte {
	return bin.Encode(
		0, // ver
		obj.Asset,
		obj.Comment,
		obj.Outs,
	)
}

func (obj *Emission) Decode(data []byte) error {
	return bin.Decode(data,
		new(int),
		&obj.Asset,
		&obj.Comment,
		&obj.Outs,
	)
}

func (out *EmissionOutput) Encode() []byte {
	return bin.Encode(
		out.Type,
		out.Address,
		out.Value,
		out.Amount,
		out.Reserved1,
		out.Reserved2,
	)
}

func (out *EmissionOutput) Decode(data []byte) error {
	return bin.Decode(data,
		&out.Type,
		&out.Address,
		&out.Value,
		&out.Amount,
		&out.Reserved1,
		&out.Reserved2,
	)
}

func (out *EmissionOutput) IsDistributionReward() bool {
	return out.Type == EmissionTypeDistribution
}

func (out *EmissionOutput) IsAuthorReward() bool {
	return out.Type == EmissionTypeAuthor
}

func (out *EmissionOutput) IsReferralReward() bool {
	return out.Type == EmissionTypeReferer
}

var oneGiB = bignum.NewInt(consts.GiB)

func (obj *Emission) AvgRatePerGiB() (s bignum.Int) {
	var v int64
	for _, out := range obj.Outs {
		if out.IsDistributionReward() {
			s.Increment(out.Amount)
			v += out.Value
		}
	}
	if v == 0 {
		return
	}
	return s.Mul(oneGiB).Div(bignum.NewInt(v))
}

func (obj *Emission) OutputByAddress(addr []byte) *EmissionOutput {
	if len(addr) > 0 {
		for _, out := range obj.Outs {
			if bytes.Equal(out.Address, addr) {
				return out
			}
		}
	}
	return nil
}

func (obj *Emission) TotalAmount() (s bignum.Int) {
	for _, out := range obj.Outs {
		s.Increment(out.Amount)
	}
	return
}

func (obj *Emission) TotalValue() (s int64) {
	for _, out := range obj.Outs {
		s += out.Value
	}
	return
}

func (obj *Emission) Verify() error {

	if !obj.Sender().Equal(obj.ChainConfig().MasterPubKey()) { // Sender of emission-tx must be EmissionPublicKey
		return ErrTxIncorrectSender
	}

	for _, out := range obj.Outs {
		if out.Amount.Sign() <= 0 {
			return ErrTxIncorrectAmount
		}
		if out.Value < 0 {
			return ErrTxIncorrectValue
		}
		if !crypto.IsValidAddress(out.Address) {
			return ErrTxIncorrectAddress
		}
	}
	return nil
}

func (obj *Emission) Execute(st *state.State) {
	for _, out := range obj.Outs {
		// add coins to attached address
		st.Increment(obj.Asset, out.Address, out.Amount, 0)
	}
}

func (t *Emission) MarshalJSON() ([]byte, error) {
	return json.Object{
		"asset":   hex.Encode(t.Asset),
		"outs":    t.Outs,
		"comment": t.Comment,
	}.Bytes(), nil
}

func (out *EmissionOutput) TypeStr() string {
	return emissionTypeStr[out.Type]
}

func (out *EmissionOutput) MarshalJSON() ([]byte, error) {
	return json.Object{
		"type":    out.Type,
		"value":   out.Value,
		"address": crypto.EncodeAddress(out.Address),
		"amount":  out.Amount,
	}.Bytes(), nil
}
