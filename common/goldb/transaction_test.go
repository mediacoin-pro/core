package goldb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

func TestTransaction_PutVar(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	// put vars
	err := store.Exec(func(tr *Transaction) {
		tr.PutVar(Key(TestTable, "keyA"), "Alice")           // string
		tr.PutVar(Key(TestTable, "keyB"), 0xB0B)             // int
		tr.PutVar(Key(TestTable, "keyC"), User{"Carol", 18}) // any object
	})

	// get vars
	var (
		v1 string
		v2 int
		v3 User
	)
	store.GetVar(Key(TestTable, "keyA"), &v1)
	store.GetVar(Key(TestTable, "keyB"), &v2)
	store.GetVar(Key(TestTable, "keyC"), &v3)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", v1)
	assert.Equal(t, 0xB0B, v2)
	assert.Equal(t, User{"Carol", 18}, v3)
}

func TestTransaction_GetVar(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	var a, b, c struct {
		v   string
		ok  bool
		err error
	}
	key := Key(TestTable, "id")
	err := store.Exec(func(tr *Transaction) {
		tr.PutVar(key, "Alice")

		// get data from transaction
		a.ok, a.err = tr.GetVar(key, &a.v)

		// get data from storage (before commit)
		b.ok, b.err = store.GetVar(key, &b.v)
	})
	// get data from storage (after commit)
	c.ok, c.err = store.GetVar(key, &c.v)

	assert.NoError(t, err)
	assert.NoError(t, a.err)
	assert.NoError(t, b.err)
	assert.NoError(t, c.err)
	assert.True(t, a.ok)
	assert.False(t, b.ok)
	assert.True(t, c.ok)
	assert.Equal(t, "Alice", a.v)
	assert.Equal(t, "", b.v)
	assert.Equal(t, "Alice", c.v)
}

func TestTransaction_Discard(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	var a, b struct {
		v   string
		ok  bool
		err error
	}
	key := Key(TestTable, "id")
	err := store.Exec(func(tr *Transaction) {
		tr.PutVar(key, "Alice")

		a.ok, a.err = tr.GetVar(key, &a.v) // get from transaction

		tr.Fail(errors.New("transaction-fail")) // discard transaction
	})
	b.ok, b.err = store.GetVar(key, &b.v) // get from storage

	assert.Error(t, err)
	assert.NoError(t, b.err)
	assert.True(t, a.ok)
	assert.False(t, b.ok)
	assert.Equal(t, "Alice", a.v)
	assert.Equal(t, "", b.v)
}

func TestTransaction_SequenceNextVal(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	var a, b, c, d uint64

	store.Exec(func(tr *Transaction) {
		a = tr.SequenceNextVal(TestTable)
		b = tr.SequenceNextVal(TestTable)
	})
	store.Exec(func(tr *Transaction) {
		c = tr.SequenceNextVal(TestTable)
		tr.Fail(errors.New("transaction-fail")) // discard transaction
	})
	store.Exec(func(tr *Transaction) {
		d = tr.SequenceNextVal(TestTable)
	})

	assert.True(t, b == a+1)
	assert.True(t, c == b+1)
	assert.True(t, d == b+1)
}

func TestTransaction_Increment(t *testing.T) {
	var v1 int64
	store := newTestStorage()
	defer store.Drop()
	key := Key(TestTable, "id")
	store.Exec(func(tr *Transaction) {
		tr.PutVar(key, 100)
	})

	err := store.Exec(func(tr *Transaction) {
		// increment value (transaction context)
		v1 = tr.Increment(key, 23)
	})
	// get value (storage context)
	v2, _ := store.GetInt(key)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), v1)
	assert.Equal(t, int64(123), v2)
}

func TestTransaction_Exec_fail(t *testing.T) {
	store := newTestStorage()
	defer store.Drop()

	key := Key(TestTable, "id")
	store.Exec(func(tr *Transaction) {
		tr.PutVar(key, 666) // put int
	})

	err := store.Exec(func(tr *Transaction) {
		tr.GetVar(key, &User{}) // fail read value
	})

	assert.Error(t, err)
}
