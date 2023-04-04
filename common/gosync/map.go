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
	vals map[any]any
}

func normKey(key any) any {
	if v, ok := key.([]byte); ok {
		return string(v)
	}
	return key
}

func (m *Map) Clear() {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.vals = map[any]any{}
	m.ver++
}

func (m *Map) Set(key, value any) {
	key = normKey(key)

	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[any]any{}
	}
	m.vals[key] = value
	m.ver++
}

func (m *Map) Increment(key any, val int64) int64 {
	key = normKey(key)

	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[any]any{}
	}
	if v, ok := m.vals[key].(int64); ok {
		val += v
	}
	m.vals[key] = val
	m.ver++
	return val
}

func (m *Map) Delete(key any) {
	key = normKey(key)

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.vals != nil {
		delete(m.vals, key)
		m.ver++
	}
}

func (m *Map) Get(key any) any {
	key = normKey(key)

	m.mx.RLock()
	defer m.mx.RUnlock()

	if m.vals != nil {
		return m.vals[key]
	}
	return nil
}

func (m *Map) Exists(key any) bool {
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

func (m *Map) ForEach(fn func(key, value any)) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	if m.vals != nil {
		for k, v := range m.vals {
			fn(k, v)
		}
	}
}

func (m *Map) KeyValues() map[any]any {
	m.mx.RLock()
	defer m.mx.RUnlock()

	res := map[any]any{}
	if m.vals != nil {
		for k, v := range m.vals {
			res[k] = v
		}
	}
	return res
}

func (m *Map) Keys() []any {
	m.mx.RLock()
	defer m.mx.RUnlock()

	vv := make([]any, 0, len(m.vals))
	if m.vals != nil {
		for key := range m.vals {
			vv = append(vv, key)
		}
	}
	return vv
}

func (m *Map) Values() []any {
	m.mx.RLock()
	defer m.mx.RUnlock()

	vv := make([]any, 0, len(m.vals))
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

func (m *Map) Pop() (key, value any) {
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

func (m *Map) PopAll() (values map[any]any) {
	m.mx.Lock()
	defer m.mx.Unlock()
	values = m.vals
	m.vals = nil
	m.ver++
	return
}

func (m *Map) RandomValue() any {
	_, v := m.Random()
	return v
}

func (m *Map) RandomKey() any {
	k, _ := m.Random()
	return k
}

func (m *Map) Random() (key, value any) {
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
func encString(v any) string {
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
