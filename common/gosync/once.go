package gosync

import (
	"sync"
	"sync/atomic"
)

// Once is an object that will perform exactly one action.
type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}

func (o *Once) Do(f func()) {
	if o.Done() {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
