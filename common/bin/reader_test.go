package bin

import (
	"bytes"
	"testing"
	"time"

	"math/big"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestReader_ReadVar(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.WriteVar(uint64(123))
	w.WriteVar("abc")
	w.WriteVar(3.1415)
	w.WriteVar([]byte{5, 6, 7})
	w.WriteVar(Bytes("binary"))
	w.WriteVar(Point{88, 99})
	w.WriteVar([]int{7, 8, 9})
	w.WriteVar(&User{666, "Devil"})
	w.WriteVar([]*User{{1, "Alice"}, {2, "Bob"}})

	var (
		i  int
		s  string
		f  float64
		bb []byte
		bt Bytes
		p  Point
		ii []int
		u  *User
		uu []*User
	)
	r := NewReader(buf)
	r.ReadVar(&i)
	r.ReadVar(&s)
	r.ReadVar(&f)
	r.ReadVar(&bb)
	r.ReadVar(&bt)
	r.ReadVar(&p)
	r.ReadVar(&ii)
	r.ReadVar(&u)
	r.ReadVar(&uu)

	assert.Equal(t, 123, i)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 3.1415, f)
	assert.Equal(t, []byte{5, 6, 7}, bb)
	assert.Equal(t, "binary", string(bt))
	assert.Equal(t, Point{88, 99}, p)
	assert.Equal(t, []int{7, 8, 9}, ii)
	assert.Equal(t, User{666, "Devil"}, *u)
	assert.Equal(t, []*User{{1, "Alice"}, {2, "Bob"}}, uu)
}

func TestReader_ReadVarInt(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVarInt(0x1234)

	r := w.Reader
	iDec, err := r.ReadVarInt()

	assert.NoError(t, err)
	assert.Equal(t, 0x1234, iDec)
}

func TestReader_ReadBigInt(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar(0)                                                                             // 0 bit
	w.WriteVar(127)                                                                           // 8 bit
	w.WriteVar(int64(0x0123456789abcdef))                                                     // 64 bit
	w.WriteVar(newBigInt("ffffffff00000000ffffffff00000000"))                                 // 128 bit
	w.WriteVar(newBigInt("ffffffff00000000ffffffff00000000ffffffff00000000ffffffff00000000")) // 256 bit
	w.WriteVar(newBigInt(strings.Repeat("a", 384/4)))                                         // 384 bit
	w.WriteVar(newBigInt(strings.Repeat("b", 512/4)))                                         // 512 bit
	w.WriteVar(newBigInt(strings.Repeat("c", 1024/4)))                                        // 1024 bit

	r := w.Reader
	b0, err0 := r.ReadBigInt()
	b8, err8 := r.ReadBigInt()
	b64, err64 := r.ReadBigInt()
	b128, err128 := r.ReadBigInt()
	b256, err256 := r.ReadBigInt()
	b384, err384 := r.ReadBigInt()
	b512, err512 := r.ReadBigInt()
	b1024, err1024 := r.ReadBigInt()

	assert.NoError(t, err0)
	assert.NoError(t, err8)
	assert.NoError(t, err64)
	assert.NoError(t, err128)
	assert.NoError(t, err256)
	assert.NoError(t, err384)
	assert.NoError(t, err512)
	assert.NoError(t, err1024)

	assert.Equal(t, b0, big.NewInt(0))
	assert.Equal(t, b8, big.NewInt(127))
	assert.Equal(t, b64, big.NewInt(0x0123456789abcdef))
	assert.Equal(t, b128, newBigInt("ffffffff00000000ffffffff00000000"))
	assert.Equal(t, b256, newBigInt("ffffffff00000000ffffffff00000000ffffffff00000000ffffffff00000000"))
	assert.Equal(t, b384, newBigInt(strings.Repeat("a", 384/4)))
	assert.Equal(t, b512, newBigInt(strings.Repeat("b", 512/4)))
	assert.Equal(t, b1024, newBigInt(strings.Repeat("c", 1024/4)))
}

func TestReader_ReadFloat64(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar(-0.123456789)

	r := w.Reader
	res, err := r.ReadFloat64()

	assert.NoError(t, err)
	assert.Equal(t, -0.123456789, res)
}

func TestReader_ReadFloat32(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar(float32(-1. / 3))

	r := w.Reader
	res, err := r.ReadFloat32()

	assert.NoError(t, err)
	assert.Equal(t, float32(-1./3), res)
}

func TestReader_ReadTime(t *testing.T) {
	time.Local = time.UTC
	w := NewBuffer(nil)
	w.WriteVar(time.Date(2016, 07, 06, 18, 24, 45, 0, time.Local))

	r := w.Reader
	res, err := r.ReadTime()

	assert.NoError(t, err)
	assert.Equal(t, "2016-07-06 18:24:45 UTC", res.Format("2006-01-02 15:04:05 MST"))
}

func TestReader_ReadTime32(t *testing.T) {
	time.Local = time.UTC
	w := NewBuffer(nil)
	w.WriteTime32(time.Date(2016, 07, 06, 18, 24, 45, 0, time.Local))

	r := w.Reader
	res, err := r.ReadTime32()

	assert.NoError(t, err)
	assert.Equal(t, "2016-07-06 18:24:45 UTC", res.Format("2006-01-02 15:04:05 MST"))
}

func TestReader_ReadMap(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar(map[int]string{2: "aa", 3: "bb", 4: "cc", 5: "dd"})
	r := w.Reader

	var mm = map[int]string{1: "00"}
	r.ReadVar(&mm)

	assert.Equal(t, map[int]string{2: "aa", 3: "bb", 4: "cc", 5: "dd"}, mm)
}

func TestReader_ReadSlice(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar([]Point{{1, 2}, {33, 44}})
	r := w.Reader

	var points = []Point{{-1, -1}}
	r.ReadSlice(&points)

	assert.Equal(t, []Point{{1, 2}, {33, 44}}, points)
}

func TestReader_ReadNilSlice(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar([]Point{})
	r := w.Reader

	var points = []Point{{-1, -1}}
	r.ReadSlice(&points)

	assert.Nil(t, points)
}

func TestReader_ReadEncoder(t *testing.T) {
	data := Encode(&User{123, "Alice"})

	var b *User
	r := NewBuffer(data)
	err := r.ReadVar(&b)

	assert.NoError(t, err)
	assert.Equal(t, User{123, "Alice"}, *b)
}

func TestReader_ReadNilPointerEncoder(t *testing.T) {
	var nilUser *User
	data := Encode(nilUser)

	var b = new(User)
	r := NewBuffer(data)
	err := r.ReadVar(&b)

	assert.Equal(t, []byte{0}, data)
	assert.NoError(t, err)
	assert.Nil(t, b)
}

func TestReader_SetReadLimit(t *testing.T) {
	arr := make([]byte, 100)
	w := NewBuffer(nil)
	w.WriteVar(arr)

	r := w.Reader
	r.SetReadLimit(101) // 100 + 1 byte of length
	res, err := r.ReadBytes()

	assert.NoError(t, err)
	assert.Equal(t, 100, len(res))
}

func TestReader_ReadLimit_Fail(t *testing.T) {
	arr := make([]byte, 100)
	w := NewBuffer(nil)
	w.WriteVar(arr)

	r := w.Reader
	r.SetReadLimit(99)
	_, err := r.ReadBytes()

	assert.Error(t, err)
}

//-----------------------------------
type Point struct {
	X int
	Y int
}

type User struct {
	ID   uint64
	Name string
}

func (u *User) Encode() []byte {
	return Encode(u.ID, u.Name)
}

func (u *User) Decode(data []byte) error {
	return Decode(data, &u.ID, &u.Name)
}
