package gosync

import (
	"sync"
	"testing"
)

func BenchmarkLocker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		locker := NewLocker()
		locker.Lock()
		locker.Unlock()
	}
}

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var mx sync.Mutex
		mx.Lock()
		mx.Unlock()
	}
}
