package replication

import (
	"time"

	"github.com/mediacoin-pro/core/chain/bcstore"
	"github.com/mediacoin-pro/core/common/safe"
	"github.com/mediacoin-pro/core/common/xlog"
	"github.com/mediacoin-pro/node/rest"
)

type Service struct {
	client *rest.Client
	bc     *bcstore.ChainStorage
}

const blocksReplicationBatchSize = 100

func Start(bc *bcstore.ChainStorage) *Service {
	s := NewService(nil, bc)
	s.StartReplication()
	return s
}

func NewService(c *rest.Client, bc *bcstore.ChainStorage) *Service {
	if c == nil {
		c = rest.NewClient("")
	}
	return &Service{
		client: c,
		bc:     bc,
	}
}

func (s *Service) StartReplication() {
	go s.startBlockchainReplication()
	go s.startMempoolReplication()
}

func (s *Service) startBlockchainReplication() {
	for {
		ok, err := s.loadBlocksBatch(s.bc.LastBlock().Num, blocksReplicationBatchSize)
		if err != nil {
			xlog.Error.Printf("replication> loadBlocksBatch Error: %v", err)
		}
		if !ok || err != nil {
			time.Sleep(5 * time.Second)
		}
	}
}

func (s *Service) loadBlocksBatch(blockOffset uint64, batchSize int) (ok bool, err error) {
	defer safe.RecoverAndReport()

	//xlog.Printf("replication> loading blocks: offset:%v size:%d ...", blockOffset, batchSize)
	blocks, err := s.client.GetBlocks(blockOffset, batchSize)
	//xlog.PrintfErr("replication> loaded blocks: %v", len(blocks), err)
	if err != nil {
		xlog.Error.Printf("replication> client.GetBlocks-Error: %v", err)
		return
	}
	if len(blocks) == 0 {
		return
	}
	if err = s.bc.PutBlock(blocks...); err != nil {
		xlog.Error.Printf("replication> bc.PutBlocks-Error: %v", err)
		return
	}
	xlog.Printf("replication> ✅ replicated block#%d ", blocks[len(blocks)-1].Num)
	return true, nil
}

func (s *Service) startMempoolReplication() {

	// todo: (it`s temporary scheme) refactor me! use decentralize replication;

	for ; ; time.Sleep(100 * time.Millisecond) {
		if ok, err := s.putMempoolTxs(); err != nil {
			xlog.Error.Printf("replication> loadBlocksBatch-Error: %v", err)
		} else if !ok {
			time.Sleep(time.Second)
		}
	}
}

func (s *Service) putMempoolTxs() (ok bool, err error) {
	defer safe.Recover()

	//-- get from local mempool
	txs, _ := s.bc.Mempool.AllTxs()
	if len(txs) == 0 {
		return
	}

	//-- put to remote node
	for _, tx := range txs {
		if err = s.client.PutTx(tx); err != nil {
			xlog.Error.Printf("replication> client.PutTx-Error: %v", err)
			return
		}
		//-- remove tx from mempool
		s.bc.Mempool.RemoveTx(tx.ID())
	}

	xlog.Printf("replication> putTxs(%d). OK", len(txs))
	return true, nil
}
