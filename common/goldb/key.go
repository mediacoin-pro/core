package goldb

import (
	"reflect"

	"github.com/mediacoin-pro/core/common/bin"
)

type Entity int

type DBEncoder interface {
	DBEncode() []byte
}

type DBDecoder interface {
	DBDecode([]byte) error
}

func Key(entityID Entity, vv ...any) []byte {
	w := bin.NewBuffer(nil)
	w.WriteVarInt(int(entityID))
	encodeKeyValues(w, vv)
	return w.Bytes()
}

func PrimaryKey(tableID Entity, id uint64) []byte {
	return Key(tableID, id)
}

func encodeKeyValues(w *bin.Buffer, vv []any) *bin.Buffer {
	for _, v := range vv {
		if s, ok := v.(string); ok {
			w.Write(append([]byte(s), 0x00))
		} else {
			w.WriteVar(v)
		}
	}
	return w
}

func encodeValue(v any) []byte {
	switch obj := v.(type) {
	case DBEncoder:
		return obj.DBEncode()
	case bin.Encoder:
		return obj.Encode()
	default:
		return bin.Encode(v)
	}
}

func decodeValue(data []byte, v any) error {
	switch obj := v.(type) {
	case DBDecoder:
		return obj.DBDecode(data)
	case bin.Decoder:
		return obj.Decode(data)
	default:
		if pp := reflect.ValueOf(v); pp.Kind() == reflect.Ptr && !pp.IsNil() {
			p := pp.Elem()
			if p.Kind() == reflect.Ptr {
				// read object in case:  var obj*Object; r.Read(&obj)
				objPtr := reflect.New(reflect.TypeOf(p.Interface()).Elem())
				if obj, ok := objPtr.Interface().(bin.Decoder); ok {
					if err := obj.Decode(data); err != nil {
						return err
					}
					p.Set(objPtr)
					return nil
				}
			}
		}
		return bin.Decode(data, v)
	}
}

func encodeUint(id uint64) []byte {
	buf := bin.NewBuffer(nil)
	buf.WriteVarUint64(id)
	return buf.Bytes()
}

func decodeUint(data []byte) (uint64, error) {
	return bin.NewBuffer(data).ReadVarUint64()
}
