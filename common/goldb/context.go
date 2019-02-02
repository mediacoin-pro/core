package goldb

import (
	"bytes"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Context is context of reading data via get or fetch-methods.
// Context is implemented by Transaction and Storage
type context struct {
	qCtx         queryContext
	fPanicOnErr  bool
	rmx          sync.RWMutex
	ReadOptions  *opt.ReadOptions
	WriteOptions *opt.WriteOptions
	countReads   int64
}

type queryContext interface {
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator
}

// CountReads returns number of waiting reads
func (c *context) CountReads() int64 {
	return atomic.LoadInt64(&c.countReads)
}

// Get returns raw data by key
func (c *context) Get(key []byte) ([]byte, error) {
	atomic.AddInt64(&c.countReads, 1)
	defer atomic.AddInt64(&c.countReads, -1)

	c.rmx.RLock()
	defer c.rmx.RUnlock()

	data, err := c.qCtx.Get(key, c.ReadOptions)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil && c.fPanicOnErr {
		panic(err)
	}
	return data, err
}

// GetInt returns uint64-data by key
func (c *context) GetInt(key []byte) (num int64, err error) {
	_, err = c.GetVar(key, &num)
	return
}

// GetBigInt returns bigint-number by key
func (c *context) GetBigInt(key []byte) (num *big.Int, err error) {
	_, err = c.GetVar(key, &num)
	return
}

// GetID returns uint64-data by key
func (c *context) GetID(key []byte) (v uint64, err error) {
	data, err := c.Get(key)
	if err != nil {
		return
	}
	if data == nil {
		return 0, nil
	}
	v, err = decodeUint(data)
	if err != nil && c.fPanicOnErr {
		panic(err)
	}
	return
}

// GetStr returns string-data by key
func (c *context) GetStr(key []byte) (s string, err error) {
	_, err = c.GetVar(key, &s)
	return
}

// GetVar get data by key and unmarshal to to variable;
// Returns true when data by key existed
func (c *context) GetVar(key []byte, v interface{}) (bool, error) {
	if data, err := c.Get(key); err != nil {
		return false, err
	} else if data == nil {
		return false, nil
	} else if err = decodeValue(data, v); err != nil {
		if c.fPanicOnErr {
			panic(err)
		}
		return false, err
	} else {
		return true, nil
	}
}

// GetNumRows fetches data by query and calculates count rows
func (c *context) GetNumRows(q *Query) (numRows uint64, err error) {
	err = c.execute(q, nil)
	numRows = q.NumRows
	return
}

// Exists returns true when exists results by query
func (c *context) Exists(q *Query) (ok bool, err error) {
	var qCopy = *q
	qCopy.Limit(1)
	err = c.execute(&qCopy, nil)
	ok = qCopy.NumRows > 0
	return
}

// Fetch fetches data by query
func (c *context) Fetch(q *Query, fnRecord func(rec Record) error) error {
	return c.execute(q, fnRecord)
}

// FetchID fetches uint64-ID by query
func (c *context) FetchID(q *Query, fnRow func(id uint64) error) error {
	return c.execute(q, func(rec Record) error {
		if id, err := decodeUint(rec.Value); err != nil {
			return err
		} else {
			return fnRow(id)
		}
	})
}

var Break = errors.New("break of fetching")

// QueryValue returns first row-value by query
func (c *context) QueryValue(q *Query, v interface{}) error {
	q.Limit(1)
	return c.Fetch(q, func(rec Record) error {
		rec.MustDecode(v)
		return nil
	})
}

// QueryIDs returns slice of row-id by query
func (c *context) QueryIDs(q *Query) (ids []uint64, err error) {
	err = c.FetchID(q, func(id uint64) error {
		ids = append(ids, id)
		return nil
	})
	return
}

// QueryID returns first row-id by query
func (c *context) QueryID(q *Query) (id uint64, err error) {
	err = c.QueryValue(q, &id)
	return
}

// LastRowID returns rowID of last record in table
func (c *context) LastRowID(tableID Entity) (rowID uint64, err error) {
	q := NewQuery(tableID).Last()
	err = c.Fetch(q, func(rec Record) error {
		rec.MustDecodeKey(&rowID)
		return nil
	})
	return
}

//------ private ------
var tail1024 = bytes.Repeat([]byte{255}, 1024)

func (c *context) execute(q *Query, fnRow func(rec Record) error) (err error) {
	q.NumRows = 0
	pfx := q.filter
	pfxLen := len(pfx)
	start := append(pfx, q.offset...)
	nStart := len(start)
	limit := q.limit
	if limit < 0 {
		limit = 1e15
	}
	skipFirst := len(q.offset) > 0

	var iter iterator.Iterator
	var iterNext func() bool
	var fnRecordFilter = q.fnFilter

	c.rmx.RLock()
	defer c.rmx.RUnlock()

	if !q.desc { // ask
		iter = c.qCtx.NewIterator(&util.Range{Start: start}, nil)
		iterNext = func() bool { return iter.Next() }

	} else { // desc
		iter = c.qCtx.NewIterator(nil, nil)
		iter.Seek(append(start, tail1024...))
		iterNext = func() bool { return iter.Prev() }
	}

	defer func() {
		if r, _ := recover().(error); r != nil && r != Break {
			err = r
		}
		iter.Release()
		if err == nil {
			err = iter.Error()
		}
		if err != nil && c.fPanicOnErr {
			panic(err)
		}
	}()

	for limit > 0 && iterNext() {
		key := iter.Key()
		if !bytes.HasPrefix(key, pfx) {
			break
		}
		if skipFirst { // skip first record if record.key == startOffset
			if len(key) >= nStart && bytes.Equal(key[:nStart], start) {
				continue
			}
			skipFirst = false
		}
		val := iter.Value()
		if fnRecordFilter != nil && !fnRecordFilter(Record{key, val}) {
			continue
		}
		q.offset = key[pfxLen:]
		limit--
		if fnRow != nil {
			if err = fnRow(Record{key, val}); err != nil {
				if err == Break {
					err = nil
				}
				break
			}
		}
		q.NumRows++
	}
	return
}
