package bin

import (
	"bytes"
	"encoding/gob"
	"io"
	"math"
	"math/big"
	"reflect"
	"time"
)

type Writer struct {
	wr         io.Writer
	err        error
	CntWritten int64
}

func NewWriter(w io.Writer) *Writer {
	if w == nil {
		w = bytes.NewBuffer(nil)
	}
	return &Writer{wr: w}
}

func (w *Writer) Error() error {
	return w.err
}

func (w *Writer) SetError(err error) {
	if err != nil {
		w.err = err
	}
}

func (w *Writer) Close() error {
	if c, ok := w.wr.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (w *Writer) write(bb []byte) error {
	w.Write(bb)
	return w.err
}

func (w *Writer) Write(bb []byte) (n int, err error) {
	if buf, ok := w.wr.(*bytes.Buffer); ok {
		n, err = buf.Write(bb)
	} else {
		var n64 int64
		n64, err = io.Copy(w.wr, bytes.NewBuffer(bb))
		n = int(n64)
	}
	w.CntWritten += int64(n)
	w.SetError(err)
	return
}

//----------- fixed types --------------
func (w *Writer) WriteNil() error {
	return w.write([]byte{0})
}

func (w *Writer) WriteByte(b byte) error {
	return w.write([]byte{b})
}

func (w *Writer) WriteUint8(i uint8) error {
	return w.WriteByte(byte(i))
}

func (w *Writer) WriteUint16(i uint16) error {
	return w.write(Uint16ToBytes(i))
}

func (w *Writer) WriteUint32(i uint32) error {
	return w.write(Uint32ToBytes(i))
}

func (w *Writer) WriteUint64(i uint64) error {
	return w.write(Uint64ToBytes(i))
}

func (w *Writer) WriteFloat32(f float32) error {
	return w.write(Uint32ToBytes(math.Float32bits(f)))
}

func (w *Writer) WriteFloat64(f float64) error {
	return w.write(Uint64ToBytes(math.Float64bits(f)))
}

func (w *Writer) WriteTime(t time.Time) error {
	return w.write(Uint64ToBytes(uint64(t.UnixNano())))
}

func (w *Writer) WriteTime32(t time.Time) error {
	return w.write(Uint32ToBytes(uint32(t.Unix())))
}

func (w *Writer) WriteBool(f bool) error {
	if f {
		return w.WriteByte(1)
	} else {
		return w.WriteByte(0)
	}
}

//----------- var types ----------------
func (w *Writer) WriteVarInt(num int) error {
	return w.WriteVarInt64(int64(num))
}

func (w *Writer) WriteVarUint64(num uint64) error {
	return w.WriteVarInt64(int64(num))
}

func (w *Writer) WriteVarInt64(i int64) error {
	if i >= 0 && i < 128 {
		return w.write([]byte{byte(i)})
	}
	var h byte = 0x80
	if i < 0 {
		h |= 0x40
		i = -i
	}
	const bufMaxLen = 8 + 1
	buf := make([]byte, bufMaxLen)
	var n byte = 0
	for i > 0 {
		n++
		buf[bufMaxLen-n] = byte(i)
		i >>= 8
	}
	buf[bufMaxLen-1-n] = h | n
	return w.write(buf[bufMaxLen-1-n:])
}

func (w *Writer) WriteBigInt(i *big.Int) error {
	if i == nil {
		return w.write([]byte{0})
	}
	sign := i.Sign()
	if sign == 0 {
		return w.write([]byte{0})
	}
	b := i.Bytes()
	n := len(b)
	if sign > 0 && n == 1 && b[0] < 128 {
		return w.write(b)
	}
	var h = byte(0x80)
	if sign < 0 {
		h |= 0x40
	}
	if n < 0x3f {
		h |= byte(n)
	} else {
		h |= 0x3f
	}
	if w.write([]byte{h}) != nil {
		return w.err
	}
	if n >= 0x3f && w.WriteVarInt(n) != nil {
		return w.err
	}
	return w.write(b)
}

func (w *Writer) WriteSliceBytes(bb [][]byte) error {
	w.WriteVarInt(len(bb))
	for _, d := range bb {
		w.WriteBytes(d)
	}
	return w.err
}

func (w *Writer) WriteBytes(bb []byte) error {
	w.WriteVarInt(len(bb))
	w.Write(bb)
	return w.err
}

func (w *Writer) WriteString(s string) error {
	w.WriteBytes([]byte(s))
	return w.err
}

func (w *Writer) WriteSliceString(ss []string) error {
	w.WriteVarInt(len(ss))
	for _, s := range ss {
		if w.WriteString(s) != nil {
			break
		}
	}
	return w.err
}

func (w *Writer) WriteError(err error) error {
	return w.WriteString(err.Error())
}

func (w *Writer) WriteVar(val ...interface{}) error {
	for _, v := range val {
		if err := w.writeVar(v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeVar(val interface{}) error {
	switch v := val.(type) {
	case nil:
		w.WriteNil()

	case int:
		w.WriteVarInt64(int64(v))
	case int8:
		w.WriteVarInt64(int64(v))
	case int16:
		w.WriteVarInt64(int64(v))
	case int32:
		w.WriteVarInt64(int64(v))
	case int64:
		w.WriteVarInt64(int64(v))

	case uint:
		w.WriteVarInt64(int64(v))
	case uint8:
		w.WriteVarInt64(int64(v))
	case uint16:
		w.WriteVarInt64(int64(v))
	case uint32:
		w.WriteVarInt64(int64(v))
	case uint64:
		w.WriteVarInt64(int64(v))

	case float32:
		w.WriteFloat32(v)
	case float64:
		w.WriteFloat64(v)
	case time.Time:
		w.WriteTime(v)
	case bool:
		w.WriteBool(v)

	case string:
		w.WriteString(v)
	case []string:
		w.WriteSliceString(v)
	case []byte:
		w.WriteBytes(v)
	case Bytes:
		w.WriteBytes(v)
	case [][]byte:
		w.WriteSliceBytes(v)

	case *big.Int:
		w.WriteBigInt(v)
	case big.Int:
		w.WriteBigInt(&v)

	case Encoder:
		if isNil(val) {
			w.WriteNil()
		} else {
			w.WriteBytes(v.Encode())
		}

	case binaryEncoder:
		if isNil(val) {
			w.WriteNil()
		} else {
			w.err = v.BinaryEncode(w.wr)
		}

	case binWriter:
		v.BinWrite(w)

	case error:
		w.WriteError(v)

	default:

		rv := reflect.ValueOf(v)
		switch rv.Kind() {

		case reflect.Slice:
			n := rv.Len()
			w.WriteVarInt(n)
			for i := 0; i < n; i++ {
				vi := rv.Index(i)
				w.WriteVar(vi.Interface())
			}

		case reflect.Map:
			keys := rv.MapKeys()
			w.WriteVarInt(len(keys))
			for _, key := range keys {
				w.WriteVar(key.Interface())
				w.WriteVar(rv.MapIndex(key).Interface())
			}

		default:
			w.err = gob.NewEncoder(w).Encode(v)
		}
	}
	return w.err
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v := reflect.ValueOf(v); v.Kind() {
	case reflect.Chan, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Map, reflect.Func:
		return v.IsNil()
	}
	return false
}
