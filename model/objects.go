package model

import "github.com/mediacoin-pro/core/common/bin"

type Objects []IObject

func NewObjects(oo ...IObject) Objects {
	return oo
}

func (oo *Objects) Push(obj IObject) {
	*oo = append(*oo, obj)
}

func (oo Objects) Encode() []byte {
	var dd = make([][]byte, len(oo))
	for i, o := range oo {
		dd[i] = Encode(o)
	}
	return bin.Encode(dd)
}

func (oo *Objects) Decode(data []byte) (err error) {
	var dd [][]byte
	if err = bin.Decode(data, &dd); err != nil {
		return
	}
	objs := make(Objects, len(dd))
	for i, d := range dd {
		if obj, err := Decode(d); err != nil {
			return err
		} else {
			objs[i] = obj
		}
	}
	*oo = objs
	return
}
