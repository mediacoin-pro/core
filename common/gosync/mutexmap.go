package gosync

import "sync"

type MutexMap struct {
	mx sync.RWMutex
	mm map[interface{}]*mxMapValue
}

type mxMapValue struct {
	cnt int
	ch  chan struct{}
}

func (m *MutexMap) Lock(key interface{}) {
	m.mx.Lock()
	if m.mm == nil {
		m.mm = map[interface{}]*mxMapValue{}
	}
	v, ok := m.mm[key]
	if !ok {
		v = &mxMapValue{ch: make(chan struct{}, 1)}
		m.mm[key] = v
	}
	v.cnt++
	m.mx.Unlock()

	v.ch <- struct{}{} // lock
}

func (m *MutexMap) Unlock(key interface{}) {
	m.mx.RLock()
	v := m.mm[key]
	m.mx.RUnlock()

	<-v.ch // unlock

	m.mx.Lock()
	if v.cnt--; v.cnt == 0 {
		//close(v.ch)
		delete(m.mm, key)
	}
	m.mx.Unlock()
}
