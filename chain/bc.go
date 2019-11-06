package chain

import (
	"errors"

	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/crypto/patricia"
)

type BCContext interface {
	Config() *Config
	LastBlockHeader() *BlockHeader
	State() *state.State
	StateTree() *patricia.Tree
	ChainTree() *patricia.Tree
	TransactionByID(txID uint64) (*Transaction, error)
	UsernameByID(userID uint64) (nick string, err error)
	UserAuthInfo(*crypto.PublicKey) *crypto.PublicKey
}

// todo: StateTree() move to *State
// todo: tx.Execute() -> stateUpdates, stateRoot, err

var EmptyBCContext emptyBCContext

type emptyBCContext struct {
}

var errInvalidBCContext = errors.New("chain> invalid BCContext")

func (bc emptyBCContext) Config() *Config {
	return DefaultConfig
}

func (bc emptyBCContext) LastBlockHeader() *BlockHeader {
	return nil
}

func (bc emptyBCContext) State() *state.State {
	panic(errInvalidBCContext)
}

func (bc emptyBCContext) StateTree() *patricia.Tree {
	panic(errInvalidBCContext)
}

func (bc emptyBCContext) ChainTree() *patricia.Tree {
	panic(errInvalidBCContext)
}

func (bc emptyBCContext) UserAuthInfo(pub *crypto.PublicKey) *crypto.PublicKey {
	return pub
}

func (bc emptyBCContext) TransactionByID(txID uint64) (tx *Transaction, err error) {
	err = errInvalidBCContext
	return
}

func (bc emptyBCContext) UsernameByID(userID uint64) (nick string, err error) {
	err = errInvalidBCContext
	return
}
