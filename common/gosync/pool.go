package gosync

import "sync"

type Pool struct {
	mx   sync.RWMutex
	vals []interface{}
}

func (q *Pool) Push(value interface{}) {
	q.mx.Lock()
	defer q.mx.Unlock()
	q.vals = append(q.vals, value)
}

func (q *Pool) Pop() (val interface{}) {
	q.mx.Lock()
	defer q.mx.Unlock()
	if len(q.vals) > 0 {
		val = q.vals[0]
		q.vals = q.vals[1:]
	}
	return
}

func (q *Pool) PopAll() (vals []interface{}) {
	q.mx.Lock()
	defer q.mx.Unlock()
	vals = q.vals
	q.vals = nil
	return
}

func (q *Pool) Size() int {
	q.mx.RLock()
	defer q.mx.RUnlock()
	return len(q.vals)
}

func (q *Pool) Values() []interface{} {
	q.mx.RLock()
	defer q.mx.RUnlock()

	vv := make([]interface{}, len(q.vals))
	copy(vv, q.vals)
	return vv
}

func (q *Pool) String() string {
	q.mx.RLock()
	defer q.mx.RUnlock()

	return encString(q.vals)
}

func (q *Pool) Strings() []string {
	q.mx.RLock()
	defer q.mx.RUnlock()

	ss := make([]string, 0, len(q.vals))
	for _, v := range q.vals {
		ss = append(ss, encString(v))
	}
	return ss
}

func (q *Pool) Clear() {
	q.mx.Lock()
	defer q.mx.Unlock()
	q.vals = nil
}
