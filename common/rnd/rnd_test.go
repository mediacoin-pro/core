package rnd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandUint64(t *testing.T) {

	N := 100000
	m8, m1 := 0., 0.

	for i := 0; i < N; i++ {
		if Uint64() < 0x8000000000000000 {
			m8++
		}
		if Uint64() < 0x1000000000000000 {
			m1++
		}
	}

	assert.InDelta(t, 0.5000, m8/float64(N), 0.01)
	assert.InDelta(t, 0.0625, m1/float64(N), 0.01)
}

func TestIntn(t *testing.T) {

	r := New("abc", 123)
	v := [3]float64{}
	N := 100000.

	for i := 0.; i < N; i++ {
		v[r.Intn(3)]++
	}

	assert.InDelta(t, 0.33333333333333, v[0]/N, 0.01)
	assert.InDelta(t, 0.33333333333333, v[1]/N, 0.01)
	assert.InDelta(t, 0.33333333333333, v[2]/N, 0.01)
}

func TestNew(t *testing.T) {
	r1 := New("abc", 123)
	r2 := New("abc", 123)

	for i := 1; i < 100000; i++ {

		v1 := r1.Intn(i)
		v2 := r2.Intn(i)

		assert.True(t, v1 == v2)
	}
}

func TestShuffle(t *testing.T) {
	slice := []byte("abcdefghijklmnopqrstuvwxyz")

	Shuffle(slice, 0)

	assert.Equal(t, "sobhlicfmetpjvyuxawqnrgzdk", string(slice))
}
