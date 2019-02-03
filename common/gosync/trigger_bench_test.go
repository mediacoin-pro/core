package gosync

import (
	"sync"
	"testing"
)

func BenchmarkWaitingViaTrigger(b *testing.B) {
	for i := 0; i < b.N; i++ {

		t := NewTrigger()

		go func() {
			// do somethings
			t.Trigger()
		}()

		t.Wait()
	}
}

func BenchmarkWaitingViaWaitGroups(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			// do somethings
			wg.Done()
		}()

		wg.Wait()
	}
}

func BenchmarkWaitingViaMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var mx sync.Mutex

		mx.Lock()
		go func() {
			// do somethings
			mx.Unlock()
		}()

		mx.Lock()
		mx.Unlock()
	}
}
