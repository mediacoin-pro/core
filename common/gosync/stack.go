package gosync

import "sync"

type Stack struct {
	mx   sync.RWMutex
	vals []interface{}
}

func (s *Stack) Push(value interface{}) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.vals = append(s.vals, value)
}

func (s *Stack) Pop() (val interface{}) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if len(s.vals) > 0 {
		val = s.vals[len(s.vals)-1]
		s.vals = s.vals[:len(s.vals)-1]
	}
	return
}

func (s *Stack) PopAll() (vals []interface{}) {
	s.mx.Lock()
	defer s.mx.Unlock()
	vals = s.vals
	s.vals = nil
	return
}

func (s *Stack) Size() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return len(s.vals)
}

func (s *Stack) Values() []interface{} {
	s.mx.RLock()
	defer s.mx.RUnlock()

	vv := make([]interface{}, len(s.vals))
	copy(vv, s.vals)
	return vv
}

func (s *Stack) Clear() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.vals = nil
}
