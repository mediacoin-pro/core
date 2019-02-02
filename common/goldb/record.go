package goldb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/mediacoin-pro/core/common/bin"
)

type Record struct {
	Key   []byte
	Value []byte
}

var errInvalidKeyData = errors.New("goldb: invalid key data")

func NewRecord(key []byte, v interface{}) Record {
	return Record{key, encodeValue(v)}
}

func (r Record) String() string {
	return fmt.Sprintf("record(%x:%x)", r.Key, r.Value)
}

//------- key ---------
func (r Record) Table() Entity {
	id, err := decodeUint(r.Key)
	panicOnErr(err)
	return Entity(id)
}

func (r Record) RowID() (id uint64) {
	r.MustDecodeKey(&id)
	return
}

func (r Record) MustDecodeKey(vv ...interface{}) {
	panicOnErr(r.DecodeKey(vv...))
}

func (r Record) DecodeKey(vv ...interface{}) (err error) {
	buf := bin.NewBuffer(r.Key)
	if _, err = buf.ReadVarInt64(); err != nil { // read tableID
		return
	}
	for _, v := range vv {
		if str, ok := v.(*string); ok { // special case - read string in Key
			if n := bytes.IndexByte(r.Key[int(buf.CntRead):], 0); n < 0 {
				return errInvalidKeyData
			} else {
				s := make([]byte, n+1)
				if _, err = buf.Read(s); err != nil {
					return
				}
				*str = string(s[:n])
			}
		} else {
			buf.ReadVar(v)
		}
		if err = buf.Error(); err != nil {
			return
		}
	}
	return
}

func (r Record) KeyOffset(q *Query) []byte {
	return r.Key[len(q.filter):]
}

//------ value --------------
func (r Record) MustDecode(v interface{}) {
	panicOnErr(r.Decode(v))
}

func (r Record) Decode(v interface{}) error {
	return decodeValue(r.Value, v)
}

func (r Record) ValueID() (id uint64) {
	r.MustDecode(&id)
	return
}

func (r Record) ValueStr() (v string) {
	r.MustDecode(&v)
	return
}

func (r Record) ValueInt() (v int64) {
	r.MustDecode(&v)
	return
}

func (r Record) ValueBigInt() (v *big.Int) {
	r.MustDecode(&v)
	return
}
