package json

import "encoding/json"

type Array []interface{}

func (arr Array) String() string {
	return string(arr.Bytes())
}

func (arr Array) Bytes() []byte {
	b, _ := json.Marshal(arr)
	return b
}

func (arr Array) Len() int {
	return len(arr)
}

func (arr Array) IsNull() bool {
	return arr == nil
}

func (arr *Array) Append(v interface{}) {
	*arr = append(*arr, v)
}

func (arr Array) Eq(i int) (v Value) {
	if arr != nil && i >= 0 && i < len(arr) {
		v.val = arr[i]
	}
	return
}

func (arr Array) EqObj(i int) Object {
	return arr.Eq(i).Object()
}

func (arr Array) ForEach(fn func(Value)) {
	for i := range arr {
		fn(arr.Eq(i))
	}
}

func (arr Array) ForEachObject(fn func(obj Object)) {
	for i := range arr {
		fn(arr.Eq(i).Object())
	}
}

func (arr Array) Objects() (ss []Object) {
	ss = make([]Object, len(arr))
	for i := range arr {
		ss[i] = arr.Eq(i).Object()
	}
	return
}

func (arr Array) Strings() (ss []string) {
	ss = make([]string, len(arr))
	for i := range arr {
		ss[i] = arr.Eq(i).String()
	}
	return
}

func (arr Array) Ints() (ss []int) {
	ss = make([]int, len(arr))
	for i := range arr {
		ss[i] = arr.Eq(i).Int()
	}
	return
}

func (arr Array) Nums() (ss []float64) {
	ss = make([]float64, len(arr))
	for i := range arr {
		ss[i] = arr.Eq(i).Float64()
	}
	return
}

func (arr Array) Unmarshal(v interface{}) error {
	return json.Unmarshal(arr.Bytes(), v)
}
