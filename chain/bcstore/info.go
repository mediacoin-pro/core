package bcstore

import "github.com/mediacoin-pro/core/chain"

type Info struct {
	Network   int                `json:"network"`    //
	ChainID   uint64             `json:"chain"`      //
	Stat      Statistic          `json:"stat"`       //
	LastBlock *chain.BlockHeader `json:"last_block"` //
}

func (s *ChainStorage) Info() (inf Info, err error) {
	inf.Network = s.Cfg.NetworkID
	inf.ChainID = s.Cfg.ChainID
	inf.Stat = s.Totals()
	inf.LastBlock = s.LastBlockHeader()
	return
}
