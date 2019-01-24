package json

import "encoding/json"

type Object map[string]interface{}

func ParseObject(data []byte) (obj Object, err error) {
	err = json.Unmarshal(data, &obj)
	if obj == nil {
		obj = Object{}
	}
	return
}

func NewObject(v interface{}) (obj Object) {
	if v == nil {
		return
	}
	if data, err := json.Marshal(v); err != nil {
		panic(err)
	} else if err := json.Unmarshal(data, &obj); err != nil {
		panic(err)
	}
	return
}

func (obj Object) String() string {
	return string(obj.Bytes())
}

func (obj Object) Bytes() []byte {
	b, _ := json.Marshal(map[string]interface{}(obj))
	return b
}

func (obj Object) IndentString() string {
	b, _ := json.MarshalIndent(obj, "", "  ")
	return string(b)
}

func (obj Object) Get(name string) (v Value) {
	if obj != nil {
		v.val = obj[name]
	}
	return
}

func (obj Object) GetBool(name string) bool {
	return obj.Get(name).Bool()
}

func (obj Object) GetStr(name string) string {
	return obj.Get(name).String()
}

func (obj Object) GetNum(name string) float64 {
	return obj.Get(name).Float64()
}

func (obj Object) GetInt(name string) int {
	return obj.Get(name).Int()
}

func (obj Object) GetInt64(name string) int64 {
	return obj.Get(name).Int64()
}

func (obj Object) GetUint64(name string) uint64 {
	return obj.Get(name).Uint64()
}

func (obj Object) GetObj(name string) Object {
	return obj.Get(name).Object()
}

func (obj Object) GetArr(name string) Array {
	return obj.Get(name).Array()
}

func (obj Object) GetNoNil(name ...string) (v Value) {
	if obj != nil {
		for _, n := range name {
			if v.val = obj[n]; v.val != nil {
				return
			}
		}
	}
	return
}

func (obj Object) Unmarshal(v interface{}) error {
	return json.Unmarshal(obj.Bytes(), v)
}

func (obj Object) Encode() []byte {
	return obj.Bytes()
}

func (obj *Object) Decode(data []byte) (err error) {
	return json.Unmarshal(data, &obj)
}

func (obj Object) Set(name string, v interface{}) Object {
	obj[name] = v
	return obj
}

func ValueToObject(v interface{}) (obj Object, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	return ParseObject(data)
}
