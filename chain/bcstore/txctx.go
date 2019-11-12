package bcstore

import (
	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/goldb"
)

type txContext struct {
	*ChainStorage
	state *state.State
}

func (s *ChainStorage) newTxContext(dbTx *goldb.Transaction) chain.BCContext {

	// make state by dbTransaction
	st := state.NewState(s.Cfg.ChainID, func(a, addr []byte) (v bignum.Int) {
		// get state from db
		dbTx.QueryValue(goldb.NewQuery(dbIdxAssetAddr, a, addr).Last(), &v)
		return
	})

	return &txContext{
		ChainStorage: s,
		state:        st,
	}
}

func (c *txContext) State() *state.State {
	return c.state
}
