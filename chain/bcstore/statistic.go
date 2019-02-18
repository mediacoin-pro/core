package bcstore

import (
	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/chain/txobj"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/enc"
)

type Statistic struct {
	Blocks    uint64     `json:"blocks"`    //
	Txs       int64      `json:"txs"`       //
	Transfers int64      `json:"transfers"` //
	Users     int64      `json:"users"`     //
	Supply    bignum.Int `json:"supply"`    //
	Traffic   bignum.Int `json:"traffic"`   //
	Rate      bignum.Int `json:"rate"`      // nanocoins for 1 GB
	BCSize    int64      `json:"bcsize"`    //

	_reserve2 int `json:"-"` //
	_reserve3 int `json:"-"` //
	_reserve4 int `json:"-"` //
	_reserve5 int `json:"-"` //
}

func (s *Statistic) New(blockNum uint64, blockTxs int) *Statistic {
	c := s.Clone()
	c.Blocks = blockNum
	c.Txs += int64(blockTxs)
	return c
}

func (s *Statistic) Clone() *Statistic {
	var c = *s
	return &c
}

func (s *Statistic) String() string {
	return enc.JSON(s)
}

func (s *Statistic) Encode() []byte {
	return bin.Encode(
		0,

		s.Blocks,
		s.Txs,
		s.Transfers,
		s.Users,
		s.Supply,
		s.Traffic,
		s.Rate,
		s.BCSize,

		s._reserve2,
		s._reserve3,
		s._reserve4,
		s._reserve5,
	)
}

func (s *Statistic) Decode(data []byte) error {
	return bin.Decode(data,
		new(int), // version

		&s.Blocks,
		&s.Txs,
		&s.Transfers,
		&s.Users,
		&s.Supply,
		&s.Traffic,
		&s.Rate,
		&s.BCSize,

		&s._reserve2,
		&s._reserve3,
		&s._reserve4,
		&s._reserve5,
	)
}

func (s *Statistic) IncrementSupplyStat(emission *txobj.Emission) {
	if !assets.IsMDC(emission.Asset) {
		return
	}
	// total supply
	s.Supply.Increment(emission.TotalAmount())

	// total traffic
	if rate := emission.AvgRatePerGiB(); !rate.IsZero() {
		s.Rate = rate
		s.Traffic.Increment(bignum.NewInt(emission.TotalValue()))
	}
}
