package bcstore

import "github.com/mediacoin-pro/core/chain"

type Info struct {
	Version   string             `json:"ver"`        //
	ChainID   uint64             `json:"chain"`      //
	Stat      *Statistic         `json:"stat"`       //
	LastBlock *chain.BlockHeader `json:"last_block"` //
	//Mempool   mempool.Info            `json:"mempool"`    //
}

func (s *ChainStorage) Info() (inf Info, err error) {
	inf.Version = "" //, s.Cfg.Version
	inf.ChainID = s.Cfg.ChainID
	inf.Stat = s.Totals()
	inf.LastBlock = s.LastBlockHeader()
	//inf.Mempool = s.Mempool.Info()
	return
}
