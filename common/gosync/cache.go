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
	if _, ok := m.vals[key]; !ok && len(m.vals) >= m.Limit {
		for k, _ := range m.vals { // remove any key
			delete(m.vals, k)
			break
		}
	}
	m.vals[key] = value
	m.ver++
}
