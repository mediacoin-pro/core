package sys

import "time"

func Sleep(dur time.Duration) {
	time.Sleep(dur)
}

func UnixNano() int64 {
	return time.Now().UnixNano()
}

func Now() time.Time {
	return time.Now()
}
