package gosync

type Cache struct {
	Map
	Limit int
}

func NewCache(size int) *Cache {
	return &Cache{
		Limit: size,
	}
}

func (m *Cache) Set(key, value interface{}) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if m.vals == nil {
		m.vals = map[interface{}]interface{}{}
	}
	key = normKey(key)
	m.vals[key] = value
	m.ver++

	// clear
	if len(m.vals) >= m.Limit {
		n := m.Limit - 10
		for k := range m.vals { // remove random keys
			delete(m.vals, k)
			if len(m.vals) <= n {
				break
			}
		}
	}
}
