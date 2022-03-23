package gosync

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"sync"
)

// A Map is a set of temporary objects that may be individually set, get and deleted.
//
// A Map is safe for use by multiple goroutines simultaneously.
type Map struct {
	mx   sync.RWMutex
	ver  uint64
	vals map[interface{}]interface{}
}

func normKey(key interface{}) interface{} {
	if v, ok := key.([]byte); ok {
		return string(v)
	}
	return key
}

func (m *Map) Clear() {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.vals = map[interface{}]interface{}{}
	m.ver++
}

func (m *Map) Set(key, value interface{}) {
	key = normKey(key)

	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	m.vals[key] = value
	m.ver++
}

func (m *Map) Increment(key interface{}, val int64) int64 {
	key = normKey(key)

	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	if v, ok := m.vals[key].(int64); ok {
		val += v
	}
	m.vals[key] = val
	m.ver++
	return val
}

func (m *Map) Delete(key interface{}) {
	key = normKey(key)

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.vals != nil {
		delete(m.vals, key)
		m.ver++
	}
}

func (m *Map) Get(key interface{}) interface{} {
	key = normKey(key)

	m.mx.RLock()
	defer m.mx.RUnlock()

	if m.vals != nil {
		return m.vals[key]
	}
	return nil
}

func (m *Map) Exists(key interface{}) bool {
	key = normKey(key)

	m.mx.RLock()
	defer m.mx.RUnlock()

	if m.vals == nil {
		return false
	}
	_, ok := m.vals[key]
	return ok
}

func (m *Map) Size() int {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return len(m.vals)
}

func (m *Map) Version() uint64 {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return m.ver
}

func (m *Map) Info() (size int, ver uint64) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return len(m.vals), m.ver
}

func (m *Map) ForEach(fn func(key, value interface{})) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	if m.vals != nil {
		for k, v := range m.vals {
			fn(k, v)
		}
	}
}

func (m *Map) KeyValues() map[interface{}]interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	res := map[interface{}]interface{}{}
	if m.vals != nil {
		for k, v := range m.vals {
			res[k] = v
		}
	}
	return res
}

func (m *Map) Keys() []interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	vv := make([]interface{}, 0, len(m.vals))
	if m.vals != nil {
		for key := range m.vals {
			vv = append(vv, key)
		}
	}
	return vv
}

func (m *Map) Values() []interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	vv := make([]interface{}, 0, len(m.vals))
	if m.vals != nil {
		for _, v := range m.vals {
			vv = append(vv, v)
		}
	}
	return vv
}

func (m *Map) String() string {
	m.mx.RLock()
	defer m.mx.RUnlock()

	ss := map[string]string{}
	if m.vals != nil {
		for k, v := range m.vals {
			ss[encString(k)] = encString(v)
		}
	}
	return encString(ss)
}

func (m *Map) Strings() []string {
	ss := make([]string, 0, len(m.vals))
	for _, v := range m.Values() {
		ss = append(ss, encString(v))
	}
	return ss
}

func (m *Map) Pop() (key, value interface{}) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals != nil {
		for key, value = range m.vals {
			delete(m.vals, key)
			m.ver++
			return
		}
	}
	return
}

func (m *Map) PopAll() (values map[interface{}]interface{}) {
	m.mx.Lock()
	defer m.mx.Unlock()
	values = m.vals
	m.vals = nil
	m.ver++
	return
}

func (m *Map) RandomValue() interface{} {
	_, v := m.Random()
	return v
}

func (m *Map) RandomKey() interface{} {
	k, _ := m.Random()
	return k
}

func (m *Map) Random() (key, value interface{}) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	if cnt := len(m.vals); cnt > 0 {
		// todo: optimize it!  (add keys slice)
		i := rand.Intn(cnt)
		for k, v := range m.vals {
			if i--; i < 0 {
				return k, v
			}
		}
	}
	return nil, nil
}

func (m *Map) BinaryEncode(w io.Writer) error {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return gob.NewEncoder(w).Encode(m.vals)
}

func (m *Map) BinaryDecode(r io.Reader) (err error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	err = gob.NewDecoder(r).Decode(&m.vals)
	m.ver++
	return
}

// String returns object as string (encode to json)
func encString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
