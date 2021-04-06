package sys

import (
	"sync"
	"sync/atomic"
	"time"
)

const timeLag = 7101 * time.Microsecond // ≈7msec
const maxInt64 = 1<<63 - 1

var (
	now = time.Now().UnixNano()

	sleepMin int64 = maxInt64
	sleepMx  sync.RWMutex
	sleepMap = map[chan struct{}]int64{}
)

func init() {
	go func() {
		for {
			time.Sleep(timeLag)
			tCur := time.Now().UnixNano()
			atomic.StoreInt64(&now, tCur)

			sleepMx.RLock()
			nxt := sleepMin
			sleepMx.RUnlock()
			if nxt > tCur {
				continue
			}

			sleepMx.Lock()
			sleepMin = maxInt64
			for c, t := range sleepMap {
				if t <= tCur {
					delete(sleepMap, c)
					close(c)
				} else if t < sleepMin {
					sleepMin = t
				}
			}
			sleepMx.Unlock()
		}
	}()
}

func Unix() int64 {
	return atomic.LoadInt64(&now) / 1e9
}

func UnixNano() int64 {
	return atomic.LoadInt64(&now)
}

func Now() time.Time {
	return time.Unix(0, UnixNano())
}

func Sleep(dur time.Duration) {
	if dur <= 0 {
		return
	}
	if dur > 5*timeLag {
		time.Sleep(dur)
		return
	}

	c := make(chan struct{})
	deadline := UnixNano() + int64(dur)

	sleepMx.Lock()
	if deadline < sleepMin {
		sleepMin = deadline
	}
	sleepMap[c] = deadline
	sleepMx.Unlock()

	<-c
}
