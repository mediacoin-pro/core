package X15

import (
	"crypto/md5"
	"log"
	"runtime"
	"testing"
	"time"
)

/**
	go test -v xnet/crypto/X15 -bench=. -benchmem

	PASS
	BenchmarkGenerateKeyBySecret-4	       1	1928276011 ns/op	55095848 B/op	 1238393 allocs/op
	BenchmarkGenerateKey-4        	   50000	     25783 ns/op	     976 B/op	      18 allocs/op
	BenchmarkSign-4               	   20000	     77236 ns/op	    3328 B/op	      52 allocs/op
	BenchmarkVerify-4             	    5000	    291804 ns/op	   50546 B/op	    1107 allocs/op
	ok  	xnet/crypto	7.305s
**/

func init() {
	// redefine trace function
	fnNum := 0
	st := new(runtime.MemStats)
	runtime.ReadMemStats(st)
	memAloc := st.TotalAlloc
	traceTime := time.Now()
	trace = func(alg string, hash []byte) {
		fnNum++
		runtime.ReadMemStats(st)
		t := time.Now()
		log.Printf("%2d %10s: %7d μs, mem:%6d KB, [%d]\t%x\n",
			fnNum,
			alg,
			t.Sub(traceTime).Nanoseconds()/1000,
			(st.TotalAlloc-memAloc)/1000,
			len(hash), md5.Sum(hash),
		)
		traceTime = time.Now()
		memAloc = st.TotalAlloc
	}
}

func BenchmarkGenerateKeyBySecret(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKeyByPassword([]byte("secret-string"), 256)
	}
}
