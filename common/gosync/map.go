package gosync

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
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

func NewMap(values map[interface{}]interface{}) Map {
	return Map{
		vals: values,
	}
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
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	m.vals[normKey(key)] = value
	m.ver++
}

func (m *Map) Increment(key interface{}, val int64) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	if v, ok := m.vals[normKey(key)].(int64); ok {
		val += v
	}
	m.vals[normKey(key)] = val
	m.ver++
}

func (m *Map) Delete(key interface{}) {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.vals, normKey(key))
	m.ver++
}

func (m *Map) Get(key interface{}) interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()
	if v, ok := m.vals[normKey(key)]; ok {
		return v
	}
	return nil
}

func (m *Map) Exists(key interface{}) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()
	_, ok := m.vals[normKey(key)]
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

func (m *Map) KeyValues() map[interface{}]interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	res := map[interface{}]interface{}{}
	for k, v := range m.vals {
		res[k] = v
	}
	return res
}

func (m *Map) Keys() []interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	vv := make([]interface{}, 0, len(m.vals))
	for key := range m.vals {
		vv = append(vv, key)
	}
	return vv
}

func (m *Map) Values() []interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()

	vv := make([]interface{}, 0, len(m.vals))
	for _, v := range m.vals {
		vv = append(vv, v)
	}
	return vv
}

func (m *Map) String() string {
	m.mx.RLock()
	defer m.mx.RUnlock()

	ss := map[string]string{}
	for k, v := range m.vals {
		ss[encString(k)] = encString(v)
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
	for key, value = range m.vals {
		delete(m.vals, key)
		m.ver++
		return
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
		n := rand.Intn(cnt)
		for k, v := range m.vals {
			if n == 0 {
				return k, v
			}
			n--
		}
		panic(1)
	}
	return nil, nil
}

func (m *Map) Dump() ([]byte, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(m.vals)
	return buf.Bytes(), err
}

func (m *Map) Restore(data []byte) (err error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&m.vals)
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
