package goldb

import (
	"time"

	"github.com/mediacoin-pro/core/common/xlog"
	"github.com/syndtr/goleveldb/leveldb"
)

func (s *Storage) ExecBatch(fn func(tx *Transaction)) error {
	var cl chan struct{}
	var pErr *error

	// put tx to batch
	s.batchMx.Lock()
	if !s.batchSync {
		s.batchSync = true
		go s.startBatchSync()
	}
	if cl, pErr = s.batchCl, s.batchErr; cl == nil { // new batch
		cl, pErr = make(chan struct{}), new(error)
		s.batchCl, s.batchErr = cl, pErr
	}
	s.batchTxs = append(s.batchTxs, fn)
	s.batchMx.Unlock()
	//---

	<-cl // waiting for batch commit
	return *pErr
}

func (s *Storage) startBatchSync() {
	for {
		// pop all txs
		s.batchMx.Lock()
		txs, cl, pErr := s.batchTxs, s.batchCl, s.batchErr
		s.batchTxs, s.batchCl, s.batchErr = nil, nil, nil
		s.batchMx.Unlock()
		//
		if len(txs) == 0 {
			time.Sleep(time.Millisecond)
			continue
		}
		t0 := time.Now()
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

		if xlog.TraceIsOn() {
			xlog.Trace.Printf("-- goldb.ExecBatch() txs:%3d,  avgT:%6.3fs", len(txs), time.Since(t0).Seconds())
		}
		if err != nil {
			xlog.Error.Printf("-- goldb.ExecBatch() txs:%3d,  avgT:%6.3fs ERR:%v", len(txs), time.Since(t0).Seconds(), err)
		}
	}
}
