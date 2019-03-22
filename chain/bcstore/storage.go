package bcstore

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/chain/mempool"
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/chain/txobj"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/goldb"
	"github.com/mediacoin-pro/core/common/gosync"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/mediacoin-pro/core/crypto/patricia"
	"github.com/mediacoin-pro/core/model"
)

type ChainStorage struct {
	Dir     string
	Cfg     *chain.Config
	db      *goldb.Storage
	Mempool *mempool.Storage

	// blocks
	mxW          sync.Mutex
	mxR          sync.RWMutex
	lastBlock    *chain.Block  //
	stat         *Statistic    //
	cacheHeaders *gosync.Cache // blockNum => *BlockHeader
	cacheTxs     *gosync.Cache // blockNum => []*Transaction
	cacheIdxTx   *gosync.Cache // idxKey => *Transaction
	middleware   []Middleware  //
}

type Middleware func(*goldb.Transaction, *chain.Block)

const (
	// tables
	dbTabHeaders   = 0x01 // (blockNum) => BlockHeader
	dbTabTxs       = 0x02 // (blockNum, txIdx) => Transaction
	dbTabChainTree = 0x03 //
	dbTabStateTree = 0x04 // (asset, addr) => sateValue
	dbTabStat      = 0x05 // (ts) => Statistic

	// indexes
	dbIdxTxID          = 0x20 // (txID)                        => txNum
	dbIdxAsset         = 0x21 // (asset, txNum)                => sateValue
	dbIdxAssetAddr     = 0x22 // (asset, addr, txNum)          => sateValue
	dbIdxAssetAddrMemo = 0x23 // (asset, addr, addrTag, txNum) => sateValue
	dbIdxUserID        = 0x24 // (userID) => txUID
	dbIdxUserNick      = 0x25 // (nick) => txUID

	dbIdxInvites    = 0x27 // (userID, txNum)               => invitedUserID
	dbIdxSrcInvites = 0x28 // (userID, txNum)               => invitedUserID
)

var (
	ErrBlockNotFound         = errors.New("block not found")
	errTxHasBeenRegistered   = errors.New("tx has been registered")
	errTxNotFound            = errors.New("tx not found")
	errUserHasBeenRegistered = errors.New("user has been registered")
	errUserNotFound          = errors.New("user not found")
	ErrAddrNotFound          = errors.New("address not found")
	errIncorrectAddress      = errors.New("incorrect address")
	errIncorrectAssetVal     = errors.New("incorrect asset value")
	errIncorrectTxState      = errors.New("incorrect tx state")
	errIncorrectChainRoot    = errors.New("incorrect chain root")
	errIncorrectStateRoot    = errors.New("incorrect state root")
)

func NewChainStorage(dir string, cfg *chain.Config) (s *ChainStorage) {
	if cfg == nil {
		cfg = chain.NewConfig()
	}
	s = &ChainStorage{
		Dir:          dir,
		Cfg:          cfg,
		db:           goldb.NewStorage(dir, nil),
		cacheHeaders: gosync.NewCache(10000),
		cacheTxs:     gosync.NewCache(1000),
		cacheIdxTx:   gosync.NewCache(50000),
		Mempool:      mempool.NewStorage(),
	}

	//if cfg.VacuumDB {
	//	s.db.Vacuum()
	//}

	// query last block
	if b, err := s.queryLastBlock(); err != nil {
		panic(err)
	} else {
		s.lastBlock = b
	}
	// query actual totals
	if st, err := s.TotalsAt(time.Time{}); err != nil {
		panic(err)
	} else {
		s.stat = st
	}
	//s.stat = &Statistic{}
	//if err := s.db.QueryValue(goldb.NewQuery(dbTabStat).Last(), &s.stat); err != nil {
	//	panic(err)
	//}

	return
}

func (s *ChainStorage) Config() *chain.Config {
	return s.Cfg
}

func (s *ChainStorage) Close() (err error) {
	return s.db.Close()
}

func (s *ChainStorage) Drop() (err error) {
	s.db.Close()
	return s.db.Drop()
}

func (s *ChainStorage) Dump(filePath string) (err error) {
	return s.db.Dump(filePath, nil)
}

func (s *ChainStorage) VacuumDB() error {
	return s.db.Vacuum()
}

func (s *ChainStorage) AddMiddleware(fn Middleware) {
	s.middleware = append(s.middleware, fn)
}

func (s *ChainStorage) ChainTree() *patricia.Tree {
	return patricia.NewTree(patricia.NewMemoryStorage(patricia.NewSubStorage(s.db, goldb.Key(dbTabChainTree))))
}

func (s *ChainStorage) StateTree() *patricia.Tree {
	return patricia.NewTree(patricia.NewMemoryStorage(patricia.NewSubStorage(s.db, goldb.Key(dbTabStateTree))))
}

// State returns state struct from db
func (s *ChainStorage) State() *state.State {
	return state.NewState(s.Cfg.ChainID, func(asset, addr []byte) (v bignum.Int) {
		if err := s.db.QueryValue(goldb.NewQuery(dbIdxAssetAddr, asset, addr).Last(), &v); err != nil {
			panic(err)
		}
		return
	})
}

//----------------- put block --------------------------
func (s *ChainStorage) PutNewBlock(txs []*chain.Transaction, miner *crypto.PrivateKey) (block *chain.Block, err error) {
	block, err = chain.GenerateNewBlock(s, txs, miner)
	if err != nil || block == nil {
		return
	}
	err = s.PutBlock(block)
	return
}

// open db.transaction; verify block; save block and index-records
func (s *ChainStorage) PutBlock(blocks ...*chain.Block) error {
	if len(blocks) == 0 {
		return nil
	}
	// lock tx exec
	s.mxW.Lock()
	defer s.mxW.Unlock()

	// verify blocks
	lastBlockHeader := s.lastBlock.BlockHeader
	for _, block := range blocks {
		for txIdx, tx := range block.Txs {
			tx.SetBlockInfo(s, block.Num, txIdx, block.Timestamp)
		}
		if err := block.Verify(lastBlockHeader, s.Cfg); err != nil {
			return err
		}
		lastBlockHeader = block.BlockHeader
	}

	var stat = s.stat.Clone()
	var txsIDs []uint64

	// open db transaction
	err := s.db.Exec(func(tr *goldb.Transaction) {

		stateTree := patricia.NewSubTree(tr, goldb.Key(dbTabStateTree))
		chainTree := patricia.NewSubTree(tr, goldb.Key(dbTabChainTree))

		for _, block := range blocks {

			// init new block statistic
			//blockStat = blockStat.New(block.Num, len(block.Txs))

			// add index on transactions
			for txIdx, tx := range block.Txs {

				txID := tx.ID()
				txUID := encodeTxUID(block.Num, txIdx)
				txsIDs = append(txsIDs, txID)

				// check transaction by txID
				if id, _ := tr.GetID(goldb.Key(dbIdxTxID, txID)); id != 0 {
					tr.Fail(errTxHasBeenRegistered)
				}

				if s.Cfg.VerifyTxsLevel >= chain.VerifyTxLevel1 {

					//-- verify sender signature
					if err := tx.Verify(s.Cfg); err != nil {
						tr.Fail(err)
					}

					//-- verify transaction state
					// make state by dbTransaction
					st := state.NewState(s.Cfg.ChainID, func(a, addr []byte) (v bignum.Int) {
						// get state from db
						tr.QueryValue(goldb.NewQuery(dbIdxAssetAddr, a, addr).Last(), &v)
						return
					})

					// execute transaction
					stateUpdates, err := tx.Execute(st)
					if err != nil {
						tr.Fail(err)
					}

					// compare result state
					if !tx.StateUpdates.Equal(stateUpdates) {
						tr.Fail(errIncorrectTxState)
					}
				}

				stat.Txs++

				obj := tx.TxObject()

				switch tx.Type {

				case model.TxEmission:
					emission := obj.(*txobj.Emission)

					//if emission.IsPrimaryEmission() {
					//	for _, out := range emission.Outs {
					//		// set last tx by source
					//		tr.Put(goldb.Key(dbIdxSourceTx, emission.Asset, out.SourceID, txUID), nil)
					//
					//		// increment last tx by source
					//		if out.Delta > 0 {
					//			delta := emission.Amount(out.Delta)
					//			tr.IncrementBig(goldb.Key(dbIdxSourceAddr, emission.Asset, out.SourceID, out.Address), delta.BigInt())
					//		}
					//	}
					//}
					// else if emission.IsReferralReward() {
					//	for _, out := range emission.Outs {
					//		if out.Delta > 0 {
					//			delta := emission.Amount(out.Delta)
					//			tr.IncrementBig(goldb.Key(dbIdxSrcInvites, emission.Asset, out.SourceID, out.Address), delta.BigInt())
					//		}
					//	}

					stat.IncrementSupplyStat(emission) // refresh totals statistic

				case model.TxTransfer:
					//tr := obj.(*txobj.Transfer)
					//stat.IncVolumeStat(tr) // refresh statistic of total transfers
					stat.Transfers++

				case model.TxUser:
					//userID := tx.Sender.ID()
					usr := obj.(*txobj.User)
					userID := usr.UserID()

					// get user by userID
					if usrTxUID, _ := tr.GetID(goldb.Key(dbIdxUserID, userID)); usrTxUID != 0 {
						tr.Fail(errUserHasBeenRegistered)
					}
					tr.PutID(goldb.Key(dbIdxUserID, userID), txUID)

					// get user by nick
					if usrTxUID, _ := tr.GetID(goldb.Key(dbIdxUserNick, usr.Nick)); usrTxUID != 0 {
						tr.Fail(errUserHasBeenRegistered)
					}
					tr.PutID(goldb.Key(dbIdxUserNick, usr.Nick), txUID)

					// referrals
					if usr.ReferrerID != 0 {
						tr.PutID(goldb.Key(dbIdxInvites, usr.ReferrerID, txUID), txUID)
					}

					// increment users counter
					stat.Users++
				}

				// put transaction data
				tr.PutVar(goldb.Key(dbTabTxs, block.Num, txIdx), tx)

				// put index transaction by txID
				tr.PutID(goldb.Key(dbIdxTxID, txID), txUID)

				// save state to db-storage
				for stIdx, v := range tx.StateUpdates {
					if v.ChainID == s.Cfg.ChainID {
						stateTree.Put(v.StateKey(), v.Balance.Bytes())

						tr.PutVar(goldb.Key(dbIdxAssetAddr, v.Asset, v.Address, txUID, stIdx), v.Balance)

						if v.Memo != 0 { // change state with memo
							tr.PutVar(goldb.Key(dbIdxAssetAddrMemo, v.Asset, v.Address, v.Memo, txUID, stIdx), v.Balance)
						}
					}
				}
			}

			// verify state root
			if stateRoot, _ := stateTree.Root(); !bytes.Equal(block.StateRoot, stateRoot) {
				tr.Fail(errIncorrectStateRoot)
			}

			// verify chain root
			chainTree.PutVar(block.Num, block.Hash())
			if chainRoot, _ := chainTree.Root(); !bytes.Equal(block.ChainRoot, chainRoot) {
				tr.Fail(errIncorrectChainRoot)
			}

			// put block
			tr.PutVar(goldb.Key(dbTabHeaders, block.Num), block.BlockHeader)

			// save totals
			stat.Blocks = block.Num
			stat.BCSize += block.Size()
			tr.PutVar(goldb.Key(dbTabStat, block.Timestamp, block.Num), stat)

			// middleware for each block
			for _, fn := range s.middleware {
				fn(tr, block)
			}
		}
	})

	if err != nil {
		return err
	}

	//--- success block commit ------

	// refresh last block and totals info
	s.mxR.Lock()
	s.lastBlock = blocks[len(blocks)-1]
	s.stat = stat
	s.mxR.Unlock()

	for _, block := range blocks {
		s.cacheHeaders.Set(block.Num, block.BlockHeader)
	}

	// remove txs from Mempool
	s.Mempool.RemoveTx(txsIDs...)

	return nil
}

func (s *ChainStorage) LastBlock() *chain.Block {
	s.mxR.RLock()
	defer s.mxR.RUnlock()
	return s.lastBlock
}

func (s *ChainStorage) LastBlockHeader() *chain.BlockHeader {
	return s.LastBlock().BlockHeader
}

func (s *ChainStorage) DBSize() int64 {
	return s.db.Size()
}

func (s *ChainStorage) Totals() *Statistic {
	s.mxR.RLock()
	defer s.mxR.RUnlock()
	return s.stat.Clone()
}

func (s *ChainStorage) CountBlocks() uint64 {
	return s.LastBlock().Num
}

func (s *ChainStorage) CountTxs() int64 {
	s.mxR.RLock()
	defer s.mxR.RUnlock()
	return s.stat.Txs
}

func (s *ChainStorage) TotalsAt(t time.Time) (totals *Statistic, err error) {
	q := goldb.NewQuery(dbTabStat).OrderDesc().Limit(1)
	if !t.IsZero() {
		q.Offset(t.UnixNano() / 1e3)
	}
	if err = s.db.QueryValue(q, &totals); err == nil && totals == nil {
		totals = &Statistic{}
	}
	return
}

func (s *ChainStorage) BlockStat(num uint64) (st *Statistic, err error) {
	h, err := s.BlockHeader(num)
	if err != nil {
		return
	}
	_, err = s.db.GetVar(goldb.Key(dbTabStat, h.Timestamp, h.Num), &st)
	return
}

func (s *ChainStorage) FetchStat(fn func(s *Statistic) error) (err error) {
	q := goldb.NewQuery(dbTabStat)
	return s.db.Fetch(q, func(rec goldb.Record) error {
		var st *Statistic
		rec.MustDecode(&st)
		return fn(st)
	})
}

func (s *ChainStorage) queryLastBlock() (block *chain.Block, err error) {
	err = s.FetchBlocks(0, 1, true, func(b *chain.Block) error {
		block = b
		return nil
	})
	if err == nil && block == nil {
		block = chain.NewBlock(chain.GenesisBlockHeader(s.Cfg), nil)
	}
	return
}

func (s *ChainStorage) GetBlock(num uint64) (block *chain.Block, err error) {
	h, err := s.BlockHeader(num)
	if err != nil {
		return
	}
	txs, err := s.BlockTxs(num)
	if err != nil {
		return
	}
	return chain.NewBlock(h, txs), nil
}

func (s *ChainStorage) GetBlocks(offset uint64, limit int64, desc bool) (blocks []*chain.Block, err error) {
	err = s.FetchBlocks(offset, limit, desc, func(block *chain.Block) error {
		blocks = append(blocks, block)
		return nil
	})
	return
}

func (s *ChainStorage) BlockHeader(num uint64) (h *chain.BlockHeader, err error) {
	if num == 0 {
		return chain.GenesisBlockHeader(s.Cfg), nil
	}
	if h, _ = s.cacheHeaders.Get(num).(*chain.BlockHeader); h != nil {
		return
	}

	// get block from db-storage
	h = new(chain.BlockHeader)
	if ok, err := s.db.GetVar(goldb.Key(dbTabHeaders, num), h); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrBlockNotFound
	}

	s.cacheHeaders.Set(num, h)
	return h, nil
}

func (s *ChainStorage) FetchBlocks(offset uint64, limit int64, desc bool, fn func(block *chain.Block) error) error {
	return s.FetchBlockHeaders(offset, limit, desc, func(h *chain.BlockHeader) error {
		if txs, err := s.BlockTxs(h.Num); err != nil {
			return err
		} else {
			return fn(chain.NewBlock(h, txs))
		}
	})
}

func (s *ChainStorage) FetchBlockHeaders(offset uint64, limit int64, desc bool, fn func(block *chain.BlockHeader) error) error {
	q := goldb.NewQuery(dbTabHeaders)
	if offset > 0 {
		q.Offset(offset)
	}
	q.Order(desc)
	if limit > 0 {
		q.Limit(limit)
	}
	return s.db.Fetch(q, func(rec goldb.Record) error {
		var block = new(chain.BlockHeader)
		rec.MustDecode(block)
		return fn(block)
	})
}

//------------ txs --------------------
func encodeTxUID(blockNum uint64, txIdx int) uint64 {
	return (blockNum << 32) | uint64(txIdx)
}

func decodeTxUID(txUID uint64) (blockNum uint64, txIdx int) {
	return txUID >> 32, int(txUID & 0xffffffff)
}

func (s *ChainStorage) addBlockInfoToTx(tx *chain.Transaction, blockNum uint64, txIdx int) (err error) {
	block, err := s.BlockHeader(blockNum)
	if err == nil {
		tx.SetBlockInfo(s, blockNum, txIdx, block.Timestamp)
	}
	return
}

func (s *ChainStorage) GetTransaction(blockNum uint64, txIdx int) (tx *chain.Transaction, err error) {
	if blockNum == 0 {
		return
	}
	ok, err := s.db.GetVar(goldb.Key(dbTabTxs, blockNum, txIdx), &tx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errTxNotFound
	}
	err = s.addBlockInfoToTx(tx, blockNum, txIdx)
	return
}

func (s *ChainStorage) transactionByUID(txUID uint64) (*chain.Transaction, error) {
	if txUID == 0 {
		return nil, nil
	}
	return s.GetTransaction(decodeTxUID(txUID))
}

func (s *ChainStorage) BlockTxs(blockNum uint64) (txs []*chain.Transaction, err error) {
	if tt, ok := s.cacheTxs.Get(blockNum).([]*chain.Transaction); ok {
		return tt, nil
	}
	var bNum uint64
	var txIdx int
	err = s.db.Fetch(goldb.NewQuery(dbTabTxs, blockNum), func(rec goldb.Record) error {
		var tx *chain.Transaction
		rec.MustDecode(&tx)
		rec.MustDecodeKey(&bNum, &txIdx)
		txs = append(txs, tx)
		return s.addBlockInfoToTx(tx, bNum, txIdx)
	})
	if err == nil && len(txs) > 0 {
		s.cacheTxs.Set(blockNum, txs)
	}
	return
}

func (s *ChainStorage) BlockTxsCount(blockNum uint64) (count int, err error) {
	num, err := s.db.GetNumRows(goldb.NewQuery(dbTabTxs, blockNum))
	return int(num), err
}

func (s *ChainStorage) TransactionByHash(txHash []byte) (*chain.Transaction, error) {
	tx, err := s.TransactionByID(chain.TxIDByHash(txHash))
	if err == nil && tx != nil && !bytes.Equal(txHash, tx.Hash()) { // collision
		return nil, nil
	}
	return tx, err
}

func (s *ChainStorage) TransactionByID(txID uint64) (*chain.Transaction, error) {
	return s.transactionByIdxKey(goldb.Key(dbIdxTxID, txID))
}

func (s *ChainStorage) transactionByIdxKey(idxKey []byte) (tx *chain.Transaction, err error) {
	if tx, _ = s.cacheIdxTx.Get(idxKey).(*chain.Transaction); tx != nil {
		return
	}
	txUID, err := s.db.GetID(idxKey)
	if err != nil {
		return
	}
	tx, err = s.transactionByUID(txUID)
	if tx != nil {
		s.cacheIdxTx.Set(idxKey, tx)
	}
	return
}

func (s *ChainStorage) fetchTransactionsByIndex(q *goldb.Query, fn func(tx *chain.Transaction) error) error {
	return s.db.Fetch(q, func(rec goldb.Record) (err error) {
		var txUID uint64
		rec.MustDecode(&txUID)
		tx, err := s.transactionByUID(txUID)
		if tx != nil && err == nil {
			return fn(tx)
		}
		return
	})
}

func (s *ChainStorage) FetchTransactions(
	offset uint64,
	limit int64,
	orderDesc bool,
	fn func(tx *chain.Transaction) error,
) error {
	q := goldb.NewQuery(dbTabTxs)
	if offset > 0 {
		q.Offset(offset>>32, int(offset&0xffffffff)) // blockNum, txIdx
	}
	if limit > 0 {
		q.Limit(limit)
	}
	q.Order(orderDesc)
	return s.db.Fetch(q, func(rec goldb.Record) error {
		var blockNum uint64
		var txIdx int
		var tx *chain.Transaction
		rec.MustDecodeKey(&blockNum, &txIdx)
		rec.MustDecode(&tx)
		if err := s.addBlockInfoToTx(tx, blockNum, txIdx); err != nil {
			return err
		}
		return fn(tx)
	})
}

func (s *ChainStorage) TransactionsByAddr(
	asset []byte,
	addr []byte,
	memo uint64,
	offset uint64,
	limit int64,
	orderDesc bool,
) (txs []*chain.Transaction, nextOffset uint64, err error) {
	err = s.FetchTransactionsByAddr(asset, addr, memo, offset, limit, orderDesc, func(tx *chain.Transaction, _ bignum.Int) error {
		txs = append(txs, tx)
		nextOffset = tx.TxUID()
		return nil
	})
	return
}

func (s *ChainStorage) FetchTransactionsByAddr(
	asset []byte,
	addr []byte,
	memo uint64,
	offset uint64,
	limit int64,
	orderDesc bool,
	fn func(tx *chain.Transaction, val bignum.Int) error,
) error {
	if len(asset) == 0 {
		asset = assets.Default
	}
	var q *goldb.Query
	if memo == 0 { // fetch transactions by address
		q = goldb.NewQuery(dbIdxAssetAddr, asset, addr)
	} else { // fetch transactions by address+memo
		q = goldb.NewQuery(dbIdxAssetAddrMemo, asset, addr, memo)
	}
	if offset > 0 {
		q.Offset(offset)
	}
	if limit <= 0 {
		limit = 1000
	}
	q.Order(orderDesc)

	var txUID uint64
	return s.db.Fetch(q, func(rec goldb.Record) error {
		if limit <= 0 {
			return goldb.Break
		}
		var _memo, _txUID uint64
		if memo == 0 {
			rec.MustDecodeKey(&asset, &addr, &_txUID)
		} else {
			rec.MustDecodeKey(&asset, &addr, &_memo, &_txUID)
		}
		if txUID == _txUID { // exclude multiple records with the same txUID
			return nil
		}
		txUID = _txUID
		tx, err := s.transactionByUID(txUID)
		if err != nil {
			return err
		}
		var v bignum.Int
		rec.MustDecode(&v)
		limit--
		return fn(tx, v)
	})
}

func (s *ChainStorage) QueryTransaction(
	asset []byte,
	addr []byte,
	memo uint64,
	offset uint64,
	orderDesc bool,
) (tx *chain.Transaction, val bignum.Int, err error) {
	err = s.FetchTransactionsByAddr(asset, addr, memo, offset, 1, orderDesc, func(t *chain.Transaction, v bignum.Int) error {
		tx, val = t, v
		return goldb.Break
	})
	return
}

func (s *ChainStorage) QueryTransactions(
	asset []byte,
	addr []byte,
	memo uint64,
	offset uint64,
	limit int64,
	orderDesc bool,
) (txs []*chain.Transaction, err error) {
	err = s.FetchTransactionsByAddr(asset, addr, memo, offset, limit, orderDesc, func(tx *chain.Transaction, _ bignum.Int) error {
		txs = append(txs, tx)
		return nil
	})
	return
}

// AddressByStr returns address by nickname "@nick", "0x<hexUserID>" or by address "MDCxxxxxxxxxxxx"
func (s *ChainStorage) AddressByStr(str string) (addr []byte, memo uint64, err error) {
	if str == "" {
		err = ErrAddrNotFound
		return
	}
	if str[0] == '@' { // address by nickname "@<nickname>"
		if u, err := s.UserByNick(str); err != nil || u == nil {
			return nil, 0, err
		} else {
			return u.Address(), 0, err
		}
	}
	if len(str) == 18 && str[:2] == "0x" { // address by userID "0x<userID:hex>"
		userID, err := strconv.ParseUint(str[2:], 16, 64)
		if err != nil {
			return nil, 0, errIncorrectAddress
		}
		addr, err := s.AddressByUserID(userID)
		return addr, 0, err
	}
	return crypto.DecodeAddress(str)
}

func (s *ChainStorage) AddressByUserID(userID uint64) (addr []byte, err error) {
	u, err := s.UserByID(userID)
	if err != nil || u == nil {
		return
	}
	return u.Address(), nil
}

func (s *ChainStorage) UsernameByID(userID uint64) (nick string, err error) {
	u, err := s.UserByID(userID)
	if u != nil {
		nick = u.Nick
	}
	return
}

func (s *ChainStorage) UserByID(userID uint64) (u *txobj.User, err error) {
	if userID == 0 {
		return
	}
	tx, err := s.transactionByIdxKey(goldb.Key(dbIdxUserID, userID))
	if err != nil || tx == nil {
		return
	}
	obj, err := tx.Object()
	if err != nil {
		return
	}
	u, ok := obj.(*txobj.User)
	if !ok || u == nil {
		err = errUserNotFound
	}
	return
}

func (s *ChainStorage) UserByNick(nick string) (u *txobj.User, err error) {
	nick = strings.ToLower(strings.TrimPrefix(nick, "@"))
	if nick == "" {
		return
	}
	tx, err := s.transactionByIdxKey(goldb.Key(dbIdxUserNick, nick))
	if err != nil || tx == nil {
		return
	}
	obj, err := tx.Object()
	if err != nil {
		return
	}
	u, ok := obj.(*txobj.User)
	if !ok || u == nil {
		err = errUserNotFound
	}
	return
}

func (s *ChainStorage) UserByAddress(addr []byte) (*txobj.User, error) {
	if u, err := s.UserByID(crypto.AddressToUserID(addr)); err == nil && u != nil && bytes.Equal(u.Address(), addr) {
		return u, err
	} else {
		return nil, err
	}
}

// UserInfoByStr returns user-info by nickname "@nick" or by address "MDCxxxxxxxxxxxxxx"
func (s *ChainStorage) UserByStr(usernameOrAddr string) (u *txobj.User, err error) {
	switch {
	case usernameOrAddr == "":
		return

	case usernameOrAddr[0] == '@': // search by nick
		return s.UserByNick(usernameOrAddr)

	case strings.HasPrefix(usernameOrAddr, "0x"): // search by ID
		userID, err := strconv.ParseUint(strings.TrimPrefix(usernameOrAddr, "0x"), 16, 64)
		if err != nil {
			return nil, err
		}
		return s.UserByID(userID)

	default:
		addr, _, err := crypto.DecodeAddress(usernameOrAddr)
		if err != nil {
			return nil, err
		}
		return s.UserByAddress(addr)
	}
}

func (s *ChainStorage) FetchAllUsers(fn func(u *txobj.User) error) (err error) {
	q := goldb.NewQuery(dbIdxUserID)
	return s.fetchTransactionsByIndex(q, func(tx *chain.Transaction) error {
		if u, ok := tx.TxObject().(*txobj.User); ok && u != nil {
			return fn(u)
		}
		return nil
	})
}

func (s *ChainStorage) FetchInvitedUsers(
	userID uint64,
	offset string,
	limit int64,
	orderDesc bool,
	fn func(u *txobj.User) error,
) (nextOffset string, err error) {
	q := goldb.NewQuery(dbIdxInvites, userID)
	if offset != "" {
		offsetUID, _ := chain.DecodeTxUID(offset)
		q.Offset(offsetUID)
	}
	q.Order(orderDesc).Limit(limit)
	err = s.fetchTransactionsByIndex(q, func(tx *chain.Transaction) error {
		if user, ok := tx.TxObject().(*txobj.User); ok && user != nil {
			nextOffset = tx.StrTxUID()
			return fn(user)
		}
		return nil
	})
	return
}

func (s *ChainStorage) QueryInvitedUsers(
	userID uint64,
	offset string,
	limit int64,
) (users []*txobj.User, nextOffset string, err error) {
	nextOffset, err = s.FetchInvitedUsers(userID, offset, limit, false, func(u *txobj.User) error {
		users = append(users, u)
		return nil
	})
	return
}

func (s *ChainStorage) GetBalance(addr, asset []byte) (balance bignum.Int, lastTx *chain.Transaction, err error) {
	lastTx, balance, err = s.QueryTransaction(asset, addr, 0, 0, true)
	return
}

func (s *ChainStorage) LastTx(addr []byte, memo uint64, asset []byte) (lastTx *chain.Transaction, err error) {
	lastTx, _, err = s.QueryTransaction(asset, addr, memo, 0, true)
	return
}

func (s *ChainStorage) AddressInfo(addr []byte, memo uint64, asset []byte) (inf *chain.AddressInfo, err error) {
	if len(asset) == 0 {
		asset = assets.Default
	}
	bal, tx, err := s.GetBalance(addr, asset)
	if err != nil {
		return
	}
	if memo != 0 {
		if tx, err = s.LastTx(addr, memo, asset); err != nil {
			return
		}
	}
	inf = &chain.AddressInfo{
		Address: addr,
		Memo:    memo,
		Asset:   asset,
		Balance: bal,
	}
	if tx != nil {
		inf.LastTxID = tx.ID()
		inf.LastTxTs = tx.BlockTs()
	}
	if user, _ := s.UserByAddress(addr); user != nil {
		inf.UserID = user.UserID()
		inf.UserNick = user.Nick
	}
	return
}
