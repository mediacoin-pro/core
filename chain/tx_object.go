package chain

import (
	"github.com/mediacoin-pro/core/chain/state"
	"github.com/mediacoin-pro/core/model"
)

type ITransaction interface {
	Encode() []byte
	Decode([]byte) error
	SetContext(*Transaction)
	Verify() error
	Execute(state *state.State)
}

func RegisterTxType(txType int, txObj ITransaction) error {
	model.RegisterModel(txType, txObj)
	return nil
}

//func TxTypeStr(typ int ) string {
//	return txTypeStr[typ]
//}
//
//func typeByObject(obj ITransaction) TxType {
//	typ := reflect.TypeOf(obj)
//	if typ.Kind() == reflect.Ptr {
//		typ = typ.Elem()
//	}
//	return txObjTypes[typ]
//}

//func newObjectByType(typ TxType) (ITransaction, error) {
//	rt, ok := txTypes[typ]
//	if !ok {
//		return nil, ErrUnsupportedTxType
//	}
//	ptr := reflect.New(rt)
//	if obj, ok := ptr.Interface().(ITransaction); ok {
//		return obj, nil
//	} else {
//		return nil, ErrUnsupportedTxType
//	}
//}
