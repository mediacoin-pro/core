package bin

import (
	"encoding/gob"
	"errors"
	"io"
	"math"
	"math/big"
	"reflect"
	"time"
)

type Reader struct {
	rd         io.Reader
	err        error
	CntRead    int64
	maxCntRead int64
}

var (
	errBinaryDataWasCorrupted = errors.New("bin.readVarInt-Error: binary data was corrupted")
	errExceededAllowableLimit = errors.New("bin.Reader.Read-Error: exceeded allowable limit")
)

func NewReader(rd io.Reader) *Reader {
	return &Reader{rd: rd}
}

func (r *Reader) Error() error {
	return r.err
}

func (r *Reader) ClearError() {
	r.err = nil
}

func (r *Reader) SetError(err error) {
	if err != nil {
		r.err = err
	}
}

func (r *Reader) Close() error {
	if c, ok := r.rd.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (r *Reader) SetReadLimit(sz int64) {
	if sz == 0 {
		r.maxCntRead = 0
	} else {
		r.maxCntRead = r.CntRead + sz
	}
}

func (r *Reader) Read(buf []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	defer func() {
		if e, _ := recover().(error); e != nil {
			err = e
		}
		if err != nil && r.err == nil {
			r.err = err
		}
	}()
	if r.maxCntRead > 0 && int64(len(buf))+r.CntRead > r.maxCntRead {
		err = errExceededAllowableLimit
		return
	}
	n, err = io.ReadFull(r.rd, buf)
	r.CntRead += int64(n)
	return
}

func (r *Reader) read(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := r.Read(buf)
	return buf, err
}

//----------- fixed types --------------
func (r *Reader) ReadUint8() (uint8, error) {
	bb, err := r.read(1)
	if len(bb) == 1 {
		return uint8(bb[0]), err
	}
	return 0, err
}

func (r *Reader) ReadUint16() (uint16, error) {
	b, err := r.read(2)
	return BytesToUint16(b), err
}

func (r *Reader) ReadUint32() (uint32, error) {
	b, err := r.read(4)
	return BytesToUint32(b), err
}

func (r *Reader) ReadUint64() (uint64, error) {
	b, err := r.read(8)
	return BytesToUint64(b), err
}

func (r *Reader) ReadFloat32() (float32, error) {
	b, err := r.read(4)
	return math.Float32frombits(BytesToUint32(b)), err
}

func (r *Reader) ReadFloat64() (float64, error) {
	b, err := r.read(8)
	return math.Float64frombits(BytesToUint64(b)), err
}

func (r *Reader) ReadTime() (time.Time, error) {
	v, err := r.ReadUint64()
	return time.Unix(0, int64(v)), err
}

func (r *Reader) ReadTime32() (time.Time, error) {
	v, err := r.ReadUint32()
	return time.Unix(int64(v), 0), err
}

func (r *Reader) ReadByte() (b byte, err error) {
	bb, err := r.read(1)
	if len(bb) > 0 {
		b = bb[0]
	}
	return
}

func (r *Reader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b != 0, err
}

//----------- var types ----------------
func (r *Reader) readVarInt() (i int64) {
	b0, err := r.ReadUint8()
	if err != nil {
		return
	}
	if b0&0x80 == 0 {
		return int64(b0)
	}
	n := int(b0 & 0x3f)
	if n > 8 {
		r.SetError(errBinaryDataWasCorrupted)
		return
	}
	bb, err := r.read(n)
	if err != nil {
		return
	}
	for _, c := range bb {
		i <<= 8
		i |= int64(c)
	}
	if b0&0x40 != 0 {
		i = -i
	}
	return
}

func (r *Reader) ReadBigInt() (i *big.Int, err error) {
	b0, err := r.ReadUint8()
	if err != nil {
		return
	}
	if b0&0x80 == 0 {
		i = big.NewInt(int64(b0))
		return
	}
	i = new(big.Int)
	var n int
	n = int(b0 & 0x3f)
	if n == 0x3f {
		n = int(r.readVarInt())
	}
	bb, err := r.read(n)
	if err != nil {
		return
	}
	i.SetBytes(bb)
	if b0&0x40 != 0 {
		i.Neg(i)
	}
	return
}

func (r *Reader) ReadVarInt() (int, error) {
	v := r.readVarInt()
	return int(v), r.err
}

func (r *Reader) ReadVarInt64() (int64, error) {
	v := r.readVarInt()
	return v, r.err
}

func (r *Reader) ReadVarUint64() (uint64, error) {
	v := r.readVarInt()
	return uint64(v), r.err
}

func (r *Reader) ReadSliceBytes() ([][]byte, error) {
	n, err := r.ReadVarInt()
	if err != nil || n == 0 {
		return nil, err
	}
	res := make([][]byte, n)
	for i := 0; i < n; i++ {
		if res[i], err = r.ReadBytes(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *Reader) ReadBytes() ([]byte, error) {
	if n, err := r.ReadVarInt(); err != nil {
		return nil, err
	} else if n > 0 {
		return r.read(n)
	}
	return nil, nil
}

func (r *Reader) ReadString() (string, error) {
	v, err := r.ReadBytes()
	return string(v), err
}

func (r *Reader) ReadSliceString() ([]string, error) {
	n, err := r.ReadVarInt()
	if err != nil || n == 0 {
		return nil, err
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		if res[i], err = r.ReadString(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *Reader) ReadError() (error, error) {
	if s, err := r.ReadString(); err != nil {
		return nil, err
	} else {
		return errors.New(s), nil
	}
}

func (r *Reader) readSlice(p reflect.Value) {
	if n, err := r.ReadVarInt(); err == nil {
		if n == 0 {
			p.Set(reflect.Zero(p.Type()))
			return
		}
		slice := reflect.MakeSlice(p.Type(), n, n)
		for i := 0; i < n && r.err == nil; i++ {
			r.ReadVar(slice.Index(i).Addr().Interface())
		}
		if r.err == nil {
			p.Set(slice)
		}
	}
}

func (r *Reader) readMap(p reflect.Value) {
	if n, err := r.ReadVarInt(); err == nil {
		if n == 0 {
			p.Set(reflect.Zero(p.Type()))
			return
		}
		mp := reflect.MakeMap(p.Type())
		key := reflect.New(p.Type().Key())
		val := reflect.New(p.Type().Elem())
		for i := 0; i < n && r.err == nil; i++ {
			if r.ReadVar(key.Interface()) == nil && r.ReadVar(val.Interface()) == nil {
				mp.SetMapIndex(key.Elem(), val.Elem())
			}
		}
		if r.err == nil {
			p.Set(mp)
		}
	}
}

func (r *Reader) ReadSlice(val interface{}) error {
	if pp := reflect.ValueOf(val); pp.Kind() == reflect.Ptr && !pp.IsNil() {
		if p := pp.Elem(); p.Kind() == reflect.Slice {
			r.readSlice(p)
		}
	}
	return r.err //break
}

func (r *Reader) ReadVar(val ...interface{}) error {
	for _, v := range val {
		if err := r.readVar(v); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) readVar(val interface{}) error {
	switch v := val.(type) {
	case *int:
		*v = int(r.readVarInt())
	case *int8:
		*v = int8(r.readVarInt())
	case *int16:
		*v = int16(r.readVarInt())
	case *int32:
		*v = int32(r.readVarInt())
	case *int64:
		*v = int64(r.readVarInt())

	case *uint:
		*v = uint(r.readVarInt())
	case *uint8:
		*v = uint8(r.readVarInt())
	case *uint16:
		*v = uint16(r.readVarInt())
	case *uint32:
		*v = uint32(r.readVarInt())
	case *uint64:
		*v = uint64(r.readVarInt())

	case *float32:
		*v, _ = r.ReadFloat32()
	case *float64:
		*v, _ = r.ReadFloat64()
	case *time.Time:
		*v, _ = r.ReadTime()
	case *bool:
		*v, _ = r.ReadBool()

	case *string:
		*v, _ = r.ReadString()
	case *[]string:
		*v, _ = r.ReadSliceString()
	case *[]byte:
		*v, _ = r.ReadBytes()
	case *Bytes:
		*v, _ = r.ReadBytes()
	case *[][]byte:
		*v, _ = r.ReadSliceBytes()

	case **big.Int:
		*v, _ = r.ReadBigInt()
	case *big.Int:
		if x, err := r.ReadBigInt(); err == nil {
			v.Set(x)
		}

	case Decoder:
		if bb, err := r.ReadBytes(); err == nil {
			r.err = v.Decode(bb)
		}
	case binaryDecoder:
		r.err = v.BinaryDecode(r.rd)

	case binReader:
		v.BinRead(r)

	case *error:
		*v, _ = r.ReadError()

	default:

		if pp := reflect.ValueOf(val); pp.Kind() == reflect.Ptr && !pp.IsNil() {
			p := pp.Elem()
			switch p.Kind() {
			case reflect.Ptr:
				// read object in case:  var obj*Object; r.Read(&obj)
				buf, err := r.ReadBytes()
				if err != nil {
					return err
				}
				if len(buf) == 0 { // set nil pointer object
					p.Set(reflect.Zero(p.Type()))
					return nil
				}
				objPtr := reflect.New(p.Type().Elem())
				if obj, ok := objPtr.Interface().(Decoder); ok {
					if r.err = obj.Decode(buf); r.err == nil {
						p.Set(objPtr)
					}
					return r.err
				}

			case reflect.Map:
				r.readMap(p)
				return r.err

			case reflect.Slice:
				r.readSlice(p)
				return r.err
			}
		}

		//case reflect.Chan, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Map, reflect.Func:
		//	if !v.IsNil() {
		//		obj.Encode(w)
		//	}
		//}

		// other type
		r.err = gob.NewDecoder(r).Decode(v)
	}
	return r.err
}
