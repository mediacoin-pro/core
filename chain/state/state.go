package state

import (
	"errors"

	"github.com/mediacoin-pro/core/common/bignum"
)

type State struct {
	chainID uint64
	getter  func(asset, addr []byte) bignum.Int //

	vals map[string]bignum.Int //
	sets Values                //
}

var (
	ErrNegativeValue = errors.New("blockchain/state: not enough funds")
	ErrInvalidKey    = errors.New("blockchain/state: invalid key")
)

func NewState(chainID uint64, getter func(asset, addr []byte) bignum.Int) *State {
	return &State{
		chainID: chainID,
		getter:  getter,
		vals:    map[string]bignum.Int{},
	}
}

func (s *State) NewSubState() *State {
	return NewState(s.chainID, s.Get)
}

func (s *State) Copy() *State {
	a := NewState(s.chainID, nil)
	for _, v := range s.sets {
		a.set(v)
	}
	return a
}

func strKey(asset, addr []byte) string {
	return string(asset) + string(addr)
}

func (s *State) Get(asset, addr []byte) bignum.Int {
	sKey := strKey(asset, addr)
	val, ok := s.vals[sKey]
	if !ok {
		if s.getter != nil {
			val = s.getter(asset, addr)
		}
		s.vals[sKey] = val
	}
	return val
}

func (s *State) GetBytes(asset, addr []byte) []byte {
	b := s.Get(asset, addr).Bytes()
	if len(b) > 0 {
		b = b[1:]
	}
	return b
}

func (s *State) SetBytes(asset, addr, value []byte) {
	b := make([]byte, len(value)+1)
	b[0] = 1
	copy(b[1:], value)
	s.Set(asset, addr, bignum.NewFromBytes(b), 0)
}

func (s *State) Values() Values {
	return s.sets
}

func (s *State) set(v *Value) {
	if v.Balance.Sign() < 0 {
		s.Fail(ErrNegativeValue)
		return
	}
	if v.ChainID == s.chainID {
		s.vals[strKey(v.Asset, v.Address)] = v.Balance
	}
	s.sets = append(s.sets, v)
}

func (s *State) Apply(vv Values) {
	for _, v := range vv {
		s.set(v)
	}
}

func (s *State) Set(asset, addr []byte, v bignum.Int, memo uint64) {
	s.set(&Value{s.chainID, asset, addr, v, memo})
}

func (s *State) CrossChainSet(chainID uint64, asset, addr []byte, v bignum.Int, memo uint64) {
	s.set(&Value{chainID, asset, addr, v, memo})
}

func (s *State) Increment(asset, addr []byte, delta bignum.Int, memo uint64) {
	if delta.IsZero() {
		return
	}
	v := s.Get(asset, addr).Add(delta)
	s.Set(asset, addr, v, memo)
}

func (s *State) Decrement(asset, addr []byte, delta bignum.Int, memo uint64) {
	s.Increment(asset, addr, delta.Neg(), memo)
}

func (s *State) Fail(err error) {
	panic(err)
}
