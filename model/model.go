package model

import (
	"errors"
	"reflect"
	"strings"

	"github.com/mediacoin-pro/core/common/bin"
)

const (
	TxEmission = 1
	TxTransfer = 2
	TxUser     = 3

	ObjDocument = 10
	ObjFile     = 11
	ObjLink     = 12
	ObjCounter  = 13
	ObjNode     = 14
)

// Usage:
//		var _ = model.RegisterModel(model.ObjectID, &Object{})

type IObject interface {
	bin.Encoder
	bin.Decoder
}

var (
	ErrUnsupportedType = errors.New("Unsupported model-type")
)

var (
	modelsByType  = map[int]reflect.Type{}
	modelsStrType = map[int]string{}
)

func RegisterModel(typ int, v IObject) error {
	modelsByType[typ] = reflect.TypeOf(v)
	modelsStrType[typ] = strings.ToLower(reflect.TypeOf(v).Elem().Name())
	return nil
}

func isNilPointer(obj interface{}) bool {
	if obj == nil {
		return true
	}
	switch v := reflect.ValueOf(obj); v.Kind() {
	case reflect.Chan, reflect.Slice, reflect.Interface, reflect.Ptr, reflect.Map, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

func TypeStr(typ int) string {
	return modelsStrType[typ]
}

func TypeOf(obj IObject) int {
	if isNilPointer(obj) {
		return 0
	}
	t := reflect.TypeOf(obj)
	for tp, rt := range modelsByType {
		if rt == t {
			return tp
		}
	}
	panic(ErrUnsupportedType)
}

func ObjectByType(typ int) (IObject, error) {
	if typ == 0 {
		return nil, nil
	}
	rt, ok := modelsByType[typ]
	if !ok {
		return nil, ErrUnsupportedType
	}
	ptr := reflect.New(rt.Elem())
	if obj, ok := ptr.Interface().(IObject); ok {
		return obj, nil
	} else {
		return nil, ErrUnsupportedType
	}
}

func Encode(obj IObject) []byte {
	w := bin.NewBuffer(nil)
	typ := int(TypeOf(obj))
	w.WriteVarInt(typ)
	if typ != 0 && obj != nil {
		w.WriteVar(obj)
	}
	return w.Bytes()
}

func Decode(data []byte) (obj IObject, err error) {
	r := bin.NewBuffer(data)
	t, err := r.ReadVarInt()
	if t != 0 && err == nil {
		if obj, err = ObjectByType(t); err == nil {
			err = r.ReadVar(obj)
		}
	}
	return
}
