package goldb

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func generateData(iteration int) (key, val []byte) {
	val = []byte(fmt.Sprintf(`{"Property":%d,"Атрибут":"%x"}`, iteration, iteration*13))
	h := md5.Sum(val) // pseudo random key
	key = h[:8]
	return
}

func putTestValue(iteration int, tr *Transaction) {
	key, val := generateData(iteration)
	tr.Put(key, val)
}

func BenchmarkContext_Get(b *testing.B) {
	b.StopTimer()

	const CountTestData = 100000
	var store = NewStorage(fmt.Sprintf("%s/test-bench-get-%d", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	// put test data
	err := store.Exec(func(tr *Transaction) {
		for i := 0; i < CountTestData; i++ {
			putTestValue(i, tr)
		}
	})
	if err != nil {
		b.Fatal(err)
	}

	// read test data
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key, _ := generateData(i % CountTestData)
		_, err := store.Get(key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTransaction_Put_ByOneRecord(b *testing.B) {
	b.StopTimer()
	store := NewStorage(fmt.Sprintf("%s/test-bench-leveldb-%x", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := store.Exec(func(tr *Transaction) {
			putTestValue(i, tr)
		})
		if err != nil {
			b.Fail()
		}
	}
	b.StopTimer()
}

func BenchmarkTransaction_Put_By10Records(b *testing.B) {
	b.StopTimer()
	store := NewStorage(fmt.Sprintf("%s/test-bench-leveldb-%x", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	b.StartTimer()
	for i := 1; i <= b.N; i++ {
		err := store.Exec(func(tr *Transaction) {
			for ; i <= b.N && i%10 != 0; i++ {
				putTestValue(i, tr)
			}
		})
		if err != nil {
			b.Fail()
		}
	}
	b.StopTimer()
}

func BenchmarkTransaction_Put_By100Records(b *testing.B) {
	b.StopTimer()
	store := NewStorage(fmt.Sprintf("%s/test-bench-leveldb-%x", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	b.StartTimer()
	for i := 1; i <= b.N; i++ {
		err := store.Exec(func(tr *Transaction) {
			for ; i <= b.N && i%100 != 0; i++ {
				putTestValue(i, tr)
			}
		})
		if err != nil {
			b.Fail()
		}
	}
	b.StopTimer()
}

func BenchmarkTransaction_Put_By1000Records(b *testing.B) {
	b.StopTimer()
	store := NewStorage(fmt.Sprintf("%s/test-bench-leveldb-%x", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	b.StartTimer()
	for i := 1; i <= b.N; i++ {
		err := store.Exec(func(tr *Transaction) {
			for ; i <= b.N && i%1000 != 0; i++ {
				putTestValue(i, tr)
			}
		})
		if err != nil {
			b.Fail()
		}
	}
	b.StopTimer()
}

func BenchmarkTransaction_Put_By10000Records(b *testing.B) {
	b.StopTimer()
	store := NewStorage(fmt.Sprintf("%s/test-bench-leveldb-%x", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	b.StartTimer()
	for i := 1; i <= b.N; i++ {
		err := store.Exec(func(tr *Transaction) {
			for ; i <= b.N && i%10000 != 0; i++ {
				putTestValue(i, tr)
			}
		})
		if err != nil {
			b.Fail()
		}
	}
	b.StopTimer()
}

func BenchmarkTransaction_Put_By100000Records(b *testing.B) {
	b.StopTimer()
	store := NewStorage(fmt.Sprintf("%s/test-bench-leveldb-%x", os.TempDir(), rand.Int()), nil)
	defer store.Drop()

	b.StartTimer()
	for i := 1; i <= b.N; i++ {
		err := store.Exec(func(tr *Transaction) {
			for ; i <= b.N && i%100000 != 0; i++ {
				putTestValue(i, tr)
			}
		})
		if err != nil {
			b.Fail()
		}
	}
	b.StopTimer()
}
