package gosync

import "sync"

type RWMutex struct {
	mx sync.RWMutex
}

func (m *RWMutex) Lock(fn func()) {
	m.mx.Lock()
	defer m.mx.Unlock()
	fn()
}

func (m *RWMutex) RLock(fn func()) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	fn()
}
