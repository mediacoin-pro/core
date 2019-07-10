package gosync

import "sync"

type Mutex struct {
	mx sync.Mutex
}

func (m *Mutex) Lock(fn func()) {
	m.mx.Lock()
	defer m.mx.Unlock()
	fn()
}
