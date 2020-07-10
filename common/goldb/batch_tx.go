package goldb

import (
	"github.com/mediacoin-pro/core/common/xlog"
	"github.com/syndtr/goleveldb/leveldb"
)

func (s *Storage) ExecBatch(fn func(tx *Transaction)) error {
	// put tx to batch
	s.batchMx.Lock()
	if s.batchEx == nil {
		s.batchEx = make(chan struct{})
		go s.startBatchSync()
	}
	if s.batchCl == nil { // new batch
		s.batchCl, s.batchErr = make(chan struct{}), new(error)
	}
	cl, pErr := s.batchCl, s.batchErr
	s.batchTxs = append(s.batchTxs, fn)
	if len(s.batchTxs) == 1 {
		close(s.batchEx)
	}
	s.batchMx.Unlock()
	//---

	<-cl // waiting for batch commit
	return *pErr
}

func (s *Storage) startBatchSync() {
	for {
		<-s.batchEx // waiting for txs

		// pop all txs
		s.batchMx.Lock()
		txs, cl, pErr := s.batchTxs, s.batchCl, s.batchErr
		s.batchEx, s.batchTxs, s.batchCl, s.batchErr = make(chan struct{}), nil, nil, nil
		s.batchMx.Unlock()

		// commit
		err := s.Exec(func(t *Transaction) {
			for _, fn := range txs {
				fn(t)
			}
		})
		*pErr = err
		close(cl)

		if err == leveldb.ErrClosed {
			break
		}
		if err != nil {
			xlog.Error.Printf("-- goldb.ExecBatch() txs:%3d ERR:%v", len(txs), err)
		}
	}
}
