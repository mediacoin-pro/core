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
	if len(m.vals) >= m.Limit {
		if _, ok := m.vals[key]; !ok {
			for k := range m.vals { // remove any key
				delete(m.vals, k)
				break
			}
		}
	}
	m.vals[key] = value
	m.ver++
}
