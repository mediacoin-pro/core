package json

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Value struct {
	val interface{}
}

func NewValue(v interface{}) Value {
	return Value{v}
}

func Parse(data []byte) (v Value, err error) {
	err = json.Unmarshal(data, &v.val)
	return
}

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.val)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.val)
}

func (v Value) Bytes() (data []byte) {
	data, _ = json.Marshal(v.val)
	return
}

func (v Value) JSON() string {
	return string(v.Bytes())
}

func (v Value) Empty() bool {
	if v.val == nil {
		return true
	}
	switch v := v.val.(type) {
	case bool:
		return !v
	case string:
		return v == ""
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return v == 0
	case map[string]interface{}:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	}
	return false
}

func (v Value) IsNull() bool {
	return v.val == nil
}

func (v Value) IsObject() bool {
	return v.Object() != nil
}

func (v Value) IsArray() bool {
	return v.Array() != nil
}

func (v Value) Array() Array {
	switch val := v.val.(type) {
	case []interface{}:
		return val

	case Array:
		return val

	default:
		return nil
	}
}

func (v Value) Object() Object {
	switch val := v.val.(type) {
	case map[string]interface{}:
		return val

	case Object:
		return val

	default:
		return nil
	}
}

func (v Value) Bool() bool {
	return v.Int64() != 0
}

func (v Value) String() string {
	if v.val == nil {
		return ""
	}
	switch val := v.val.(type) {
	case string:
		return val
	case int:
		return strconv.FormatInt(int64(val), 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case int64:
		return strconv.FormatInt(int64(val), 10)
	case uint64:
		return strconv.FormatUint(uint64(val), 10)
	case float64:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case interface {
		String() string
	}:
		return val.String()
	}
	return v.JSON()
}

func (v Value) Int() int {
	return int(v.Int64())
}

func (v Value) Uint64() uint64 {
	return uint64(v.Int64())
}

func (v Value) Int64() (num int64) {
	if v.val == nil {
		return
	}
	switch v := v.val.(type) {
	case int:
		num = int64(v)
	case uint:
		num = int64(v)
	case int64:
		num = int64(v)
	case uint64:
		num = int64(v)
	case int32:
		num = int64(v)
	case uint32:
		num = int64(v)
	case float32:
		num = int64(v)
	case float64:
		num = int64(v)
	case bool:
		if v {
			num = 1
		}
	case string:
		if strings.IndexByte(v, '.') >= 0 {
			f, _ := strconv.ParseFloat(v, 64)
			return int64(f)
		}
		num, _ = strconv.ParseInt(v, 0, 64)
	}
	return
}

func (v Value) Float64() (num float64) {
	switch v := v.val.(type) {
	case int:
		num = float64(v)
	case uint:
		num = float64(v)
	case int32:
		num = float64(v)
	case uint32:
		num = float64(v)
	case int64:
		num = float64(v)
	case uint64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = float64(v)
	case string:
		num, _ = strconv.ParseFloat(v, 64)
	}
	return
}

func IsNum(v interface{}) bool {
	switch v.(type) {
	case int, uint, int32, uint32, int64, uint64, float32, float64:
		return true
	}
	return false
}
