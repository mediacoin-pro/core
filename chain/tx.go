package chain

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/common/hex"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/crypto/merkle"
	"github.com/mediacoin-pro/core/model"
)

const (
	MaxTxDataSize = 4 * 1024
)

type Transaction struct {
	// Tx data
	Type      int               // tx-type
	Version   int               // tx version
	Network   int               // networkID
	ChainID   uint64            //
	Nonce     uint64            // sender nonce (by default: Unix-time in µsec)
	Data      []byte            // encoded tx-object
	Reserved1 []byte            //
	Reserved2 []byte            //
	Sender    *crypto.PublicKey // tx-sender
	Sig       []byte            // tx-sender signature

	// Chain data
	StateUpdates state.Values // state changes (is not filled by sender)

	// not imported fields
	blockNum uint64            // block-num
	blockIdx int               // tx-index in block
	blockTs  int64             // block-timestamp in µsec
	_obj     ITransaction      //
	bc       BCContext         //
	_users   map[uint64]string // cache of user nicks for current transaction
}

func NewTx(
	bc BCContext,
	sender *crypto.PublicKey,
	prv *crypto.PrivateKey,
	nonce uint64,
	obj ITransaction,
) *Transaction {
	if nonce == 0 {
		nonce = NewNonce()
	}
	cfg := DefaultConfig
	if bc != nil {
		cfg = bc.Config()
	}
	if sender == nil {
		sender = prv.PublicKey()
	}
	tx := &Transaction{
		Type:    model.TypeOf(obj), //
		Version: 0,                 //
		Network: cfg.NetworkID,     //
		ChainID: cfg.ChainID,       //
		Sender:  sender,            //
		Nonce:   nonce,             //
		Data:    obj.Encode(),      // encoded tx-object

		bc:   bc,
		_obj: obj,
	}
	obj.SetContext(tx)
	tx.Sig = prv.Sign(tx.Hash()) // set sender`s signature
	return tx
}

var (
	ErrTxEmptySender      = errors.New("tx-verify-error: empty tx-sender")
	ErrTxEmptyData        = errors.New("tx-verify-error: empty tx-data")
	ErrTxInvalidData      = errors.New("tx-verify-error: invalid tx-data")
	ErrTxInvalidChainID   = errors.New("tx-verify-error: invalid chain-id")
	ErrTxInvalidNetworkID = errors.New("tx-verify-error: invalid network-id")
	ErrTxDataIsTooLong    = errors.New("tx-verify-error: tx is too long")
)

func (tx *Transaction) String() string {
	//if obj, err := tx.Object(); err == nil {
	//	return enc.IndentJSON(obj)
	//}
	return enc.JSON(tx)
}

func (tx *Transaction) ID() uint64 {
	return TxIDByHash(tx.Hash())
}

func (tx *Transaction) StrID() string {
	return hex.EncodeUint(tx.ID())
}

func (tx *Transaction) SenderID() uint64 {
	if tx != nil && tx.Sender != nil {
		return tx.Sender.ID()
	}
	return 0
}

func (tx *Transaction) SenderAddress() []byte {
	if tx != nil && tx.Sender != nil {
		return tx.Sender.Address()
	}
	return nil
}

func (tx *Transaction) SenderAddressStr() string {
	return crypto.EncodeAddress(tx.SenderAddress())
}

func (tx *Transaction) SenderNick() (nick string, err error) {
	return tx.UsernameByID(tx.Sender.ID())
}

// Hash returns hash of senders data
func (tx *Transaction) Hash() []byte {
	if tx == nil {
		return nil
	}
	return bin.Hash256(
		tx.Type,
		tx.Version,
		tx.Network,
		tx.ChainID,
		tx.Nonce,
		tx.Sender,
		tx.Data,
		tx.Reserved1,
		tx.Reserved2,
	)
}

func (tx *Transaction) TxStHash() []byte {
	return merkle.Root(tx.Hash(), tx.StateUpdates.Hash())
}

func (tx *Transaction) Size() int {
	return len(tx.Encode())
}

func (tx *Transaction) StrType() string {
	return model.TypeStr(tx.Type)
}

func (tx *Transaction) Equal(tx1 *Transaction) bool {
	return bytes.Equal(tx.Encode(), tx1.Encode())
}

func (tx *Transaction) StateAddressTotal(asset, addr []byte) (v bignum.Int) {
	if s := tx.StateUpdates.Find(asset, addr); s != nil {
		v = s.Balance
	}
	return
}

func (tx *Transaction) SetBlockInfo(bc BCContext, blockNum uint64, blockTxIdx int, blockTs int64) {
	tx.bc, tx.blockNum, tx.blockIdx, tx.blockTs = bc, blockNum, blockTxIdx, blockTs
}

func (tx *Transaction) BCContext() BCContext {
	if tx != nil {
		return tx.bc
	}
	return nil
}

func (tx *Transaction) UsernameByID(userID uint64) (nick string, err error) {
	if tx._users != nil {
		if s, ok := tx._users[userID]; ok {
			return s, nil
		}
	}
	if tx.bc != nil {
		if nick, err = tx.bc.UsernameByID(userID); err == nil && nick != "" {
			if tx._users == nil {
				tx._users = map[uint64]string{}
			}
			tx._users[userID] = nick
		}
	}
	return
}

func (tx *Transaction) BlockNum() uint64 {
	return tx.blockNum
}

func (tx *Transaction) BlockIdx() int {
	return tx.blockIdx
}

// BlockTs returns timestamp in µsec
func (tx *Transaction) BlockTs() int64 {
	return tx.blockTs
}

func (tx *Transaction) TxUID() uint64 {
	return makeTxUID(tx.blockNum, tx.blockIdx)
}

func (tx *Transaction) StrTxUID() string {
	return EncodeTxUID(tx.TxUID())
}

func (tx *Transaction) Seq() uint64 {
	return (tx.blockNum << 32) | uint64(tx.blockIdx)
}

func (tx *Transaction) Encode() []byte {
	if len(tx.Data) == 0 {
		panic(ErrTxEmptyData)
	}
	return bin.Encode(
		tx.Type,
		tx.Version,
		tx.Network,
		tx.ChainID,
		tx.Nonce,
		tx.Data,
		tx.Reserved1,
		tx.Reserved2,
		tx.Sender,
		tx.Sig,
		tx.StateUpdates,
	)
}

func (tx *Transaction) Decode(data []byte) (err error) {
	return bin.Decode(data,
		&tx.Type,
		&tx.Version,
		&tx.Network,
		&tx.ChainID,
		&tx.Nonce,
		&tx.Data,
		&tx.Reserved1,
		&tx.Reserved2,
		&tx.Sender,
		&tx.Sig,
		&tx.StateUpdates,
	)
}

func (tx *Transaction) TxObject() ITransaction {
	obj, _ := tx.Object()
	return obj
}

func (tx *Transaction) Object() (obj ITransaction, err error) {
	if tx == nil {
		return
	}
	if tx._obj != nil {
		return tx._obj, nil
	}
	o, err := model.ObjectByType(tx.Type)
	if err != nil {
		return
	}
	obj = o.(ITransaction)
	obj.SetContext(tx)
	if err = obj.Decode(tx.Data); err != nil {
		return
	}
	tx._obj = obj
	return
}

// Timestamp returns user timestamp from nonce
func (tx *Transaction) Timestamp() time.Time {
	return time.Unix(0, int64(tx.blockTs)*1e3)
}

func (tx *Transaction) Verify(cfg *Config) error {

	//-- verify transaction data
	if tx.Network != cfg.NetworkID {
		return ErrTxInvalidNetworkID
	}
	if tx.ChainID != cfg.ChainID {
		return ErrTxInvalidChainID
	}
	if len(tx.Data) == 0 {
		return ErrTxEmptyData
	}
	if tx.Type != 0 && len(tx.Data) > MaxTxDataSize {
		return ErrTxDataIsTooLong
	}
	if tx.Sender == nil || tx.Sender.Empty() {
		return ErrTxEmptySender
	}
	txObj, err := tx.Object()
	if err != nil {
		return err
	}
	if err := txObj.Verify(); err != nil {
		return err
	}

	//-- verify sender signature
	if txHash := tx.Hash(); !tx.senderAuth().Verify(txHash, tx.Sig) &&
		!cfg.MasterPubKey().Verify(txHash, tx.Sig) {
		return ErrInvalidTxSig
	}
	return nil
}

func (tx *Transaction) senderAuth() *crypto.PublicKey {
	if tx.bc != nil {
		return tx.bc.UserAuthInfo(tx.Sender)
	}
	return tx.Sender
}

// Execute executes tx, changes state, returns state-updates
func (tx *Transaction) Execute(s *state.State) (updates state.Values, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tx.Execute-panic: %v", r)
		}
	}()

	obj, err := tx.Object()
	if err != nil {
		return
	}

	newState := s.NewSubState()

	obj.Execute(newState)

	updates = newState.Values()

	return
}

func TxIDByHash(txHash []byte) uint64 {
	return bin.BytesToUint64(txHash)
}

func makeTxUID(blockNum uint64, txIdx int) uint64 {
	return (blockNum << 32) | uint64(txIdx)
}

func EncodeTxUID(txUID uint64) string {
	return hex.EncodeUint(txUID)
}

func DecodeTxUID(s string) (txUID uint64, err error) {
	return strconv.ParseUint(s, 16, 64)
}
