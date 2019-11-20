package chain

import (
	"fmt"

	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/crypto/patricia"
)

type BCContext interface {
	Config() *Config
	LastBlockHeader() *BlockHeader
	State() *state.State
	StateTree() *patricia.Tree
	ChainTree() *patricia.Tree
	TransactionByID(txID uint64) (*Transaction, error)
}

var UserNameByID = func(userID uint64) (nick string) {
	return fmt.Sprintf("0x%016x", userID)
}

type txContext struct {
	BCContext
	state *state.State
}

func (c *txContext) State() *state.State {
	return c.state
}

func NewSubContext(bc BCContext) BCContext {
	return &txContext{
		BCContext: bc,
		state:     bc.State().NewSubState(),
	}
}
