package gosync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Locker chan struct{}

func NewLocker() Locker {
	return make(Locker, 1)
}
func (l Locker) Locked() bool {
	return len(l) > 0
}
func (l Locker) Lock() {
	l <- struct{}{}
}
func (l Locker) Unlock() {
	<-l
}

func TestLocker(t *testing.T) {
	locker := NewLocker()

	const n = 10000

	var wg sync.WaitGroup
	wg.Add(n)
	k := 0
	for i := 0; i < n; i++ {
		go func() {
			locker.Lock()
			k++
			locker.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, n, k)
}
