package gosync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutexMap_Lock(t *testing.T) {
	var mx MutexMap
	counters := make([]int, 10)
	totalCount := 0

	asyncIncrementByKey := func(key int) {
		mx.Lock(key)
		defer mx.Unlock(key)

		counters[key] = counters[key] + 1
		totalCount = totalCount + 1
	}

	var wg sync.WaitGroup
	wg.Add(99990)
	for i := 0; i < 99990; i++ {
		go func(i int) {
			defer wg.Done()
			asyncIncrementByKey(i % 10)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, 0, len(mx.mm))
	assert.Equal(t, 10, len(counters))

	for _, cnt := range counters {
		assert.Equal(t, 9999, cnt)
	}
	assert.NotEqual(t, 99990, totalCount)
}
