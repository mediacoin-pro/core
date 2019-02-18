package mempool

import (
	"bytes"
	"sync"

	"github.com/mediacoin-pro/core/chain"
)

type Storage struct {
	mx  sync.RWMutex
	txs map[uint64]*chain.Transaction
}

type Info struct {
	Size int `json:"size"`
}

func NewStorage() *Storage {
	return &Storage{
		txs: map[uint64]*chain.Transaction{},
	}
}

func (s *Storage) Info() (i Info) {
	i.Size = s.Size()
	return
}

func (s *Storage) Size() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return len(s.txs)
}

// todo: add counters by txType
func (s *Storage) SizeOf(txType int) (count int) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, tx := range s.txs {
		if tx.Type == txType {
			count++
		}
	}
	return
}

func (s *Storage) Put(tx *chain.Transaction) (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.txs[tx.ID()] = tx
	return
}

func (s *Storage) PutTxs(txs []*chain.Transaction) (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, tx := range txs {
		s.txs[tx.ID()] = tx
	}
	return
}

func (s *Storage) Pop() (tx *chain.Transaction) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if len(s.txs) > 0 {
		for txID, tx := range s.txs {
			delete(s.txs, txID)
			return tx
		}
	}
	return
}

func (s *Storage) PopAll() (txs []*chain.Transaction) {
	s.mx.Lock()
	vv := s.txs
	s.txs = map[uint64]*chain.Transaction{}
	s.mx.Unlock()

	txs = make([]*chain.Transaction, 0, len(vv))
	for _, tx := range vv {
		txs = append(txs, tx)
	}
	return
}

func (s *Storage) TxsByAddress(addr []byte) (txs []*chain.Transaction, err error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, tx := range s.txs {
		if bytes.Equal(tx.SenderAddress(), addr) {
			txs = append(txs, tx)
		}
	}
	return
}

func (s *Storage) AllTxs() (txs []*chain.Transaction, err error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	txs = make([]*chain.Transaction, 0, len(s.txs))
	for _, tx := range s.txs {
		txs = append(txs, tx)
	}
	return
}

func (s *Storage) RemoveTxs(txIDs []uint64) (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	for _, txID := range txIDs {
		delete(s.txs, txID)
	}
	return
}
