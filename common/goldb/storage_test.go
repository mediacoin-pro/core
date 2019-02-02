package goldb

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

//------------------------------------
const (
	TestTable = iota + 1
)

func newTestStorage() *Storage {
	return NewStorage(fmt.Sprintf("%s/test-goldb-%x.db", os.TempDir(), rand.Int()), nil)
}

func TestStorage_Close(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	err1 := store.Close()
	err2 := store.Close()

	assert.NotNil(t, store)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestContext_Fetch(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	// put data
	store.Exec(func(tr *Transaction) {
		tr.PutVar(Key(TestTable, "A", 1), "Alice")
		tr.PutVar(Key(TestTable, "B", 2), "Bob")
		tr.PutVar(Key(TestTable, "C", 3), "Cat")
		tr.PutVar(Key(TestTable, "A", 4), "Alina")
	})

	// query all rows of TestTable
	q := NewQuery(TestTable)
	store.Fetch(q, nil)

	// query rows of TestTable where second part of key is "A"
	qA := NewQuery(TestTable, "A")
	store.Fetch(qA, nil)

	// query rows of TestTable where second part of key is "A" and third part more than 1
	qA2 := NewQuery(TestTable, "A").Offset(1)
	store.Fetch(qA2, nil)

	assert.Equal(t, 4, int(q.NumRows))
	assert.Equal(t, 2, int(qA.NumRows))
	assert.Equal(t, 1, int(qA2.NumRows))
}

func fileExists(path string) bool {
	st, _ := os.Stat(path)
	return st != nil
}

func TestStorage_Vacuum(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	//------- insert test data ------------
	const countRows = 3000
	for i := 0; i < countRows; i++ {
		store.Exec(func(tr *Transaction) {
			tr.PutVar(Key(TestTable, "LongLongLongKey%d", i*15551%countRows), "String value")
		})
	}
	sizeBefore := store.Size()
	t.Log("\tStorage.Vacuum: start.  Storage-size: ", sizeBefore)

	// run several read routines
	var wg sync.WaitGroup
	var fFinishVacuum bool
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; !fFinishVacuum; i++ {
				i %= countRows
				v, _ := store.GetStr(Key(TestTable, "LongLongLongKey%d", i*15551%countRows))
				if !assert.Equal(t, "String value", v) {
					break
				}
			}
		}()
	}

	//----- vacuum db ------------
	err := store.Vacuum()

	sizeAfter := store.Size()
	t.Log("\tStorage.Vacuum: finish. Storage-size: ", sizeAfter)

	fFinishVacuum = true
	wg.Wait()

	//----- asserts ------------
	assert.NoError(t, err)
	assert.True(t, sizeAfter < sizeBefore/50)
	assert.True(t, fileExists(store.dir))
	assert.False(t, fileExists(store.dir+".reindex"))
	assert.False(t, fileExists(store.dir+".old"))
}

func TestStorage_Vacuum_Parallel(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	// insert test data
	const countRows = 1000
	for i := 0; i < countRows; i++ {
		store.Exec(func(tr *Transaction) {
			tr.PutVar(Key(TestTable, i), "First value")
		})
	}

	// parallel update all rows
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < countRows; i++ {
			err := store.Exec(func(tr *Transaction) {
				tr.PutVar(Key(TestTable, i), "Second value")
			})
			if !assert.NoError(t, err) {
				break
			}
		}
	}()

	// start vacuum
	err := store.Vacuum()

	assert.NoError(t, err)

	// check data
	wg.Wait()
	for i := 0; i < countRows; i++ {
		val, err := store.GetStr(Key(TestTable, i))
		assert.NoError(t, err)
		assert.Equal(t, "Second value", val)
	}
}
