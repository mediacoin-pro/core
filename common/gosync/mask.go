package gosync

import (
	"sync"
)

type Mask struct {
	size int
	cnt  int
	mm   []byte
	mx   sync.RWMutex
}

func NewMask(size int) *Mask {
	n := size >> 3
	if size&7 != 0 {
		n++
	}
	return &Mask{
		mm:   make([]byte, n),
		size: size,
	}
}

func (m *Mask) Size() int {
	return m.size
}

func (m *Mask) IsDone() bool {
	return m.DoneCount() >= m.size
}

func (m *Mask) DoneCount() int {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return m.cnt
}

func (m *Mask) IsSet(i int) bool {
	i, f := i>>3, byte(1<<uint8(i&7))

	m.mx.RLock()
	defer m.mx.RUnlock()
	return m.mm[i]&f != 0
}

func (m *Mask) Set(i int) bool {
	i, f := i>>3, byte(1<<uint8(i&7))

	m.mx.Lock()
	defer m.mx.Unlock()

	if m.mm[i]&f == 0 {
		m.cnt++
		m.mm[i] |= f
		return true
	}
	return false
}

func (m *Mask) Bytes() []byte {
	return m.Encode()
}

func (m *Mask) Decode(mask []byte) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	copy(m.mm, mask)
	m.cnt = 0
	for _, b := range m.mm {
		m.cnt += int(b&1 +
			(b>>1)&1 +
			(b>>2)&1 +
			(b>>3)&1 +
			(b>>4)&1 +
			(b>>5)&1 +
			(b>>6)&1 +
			(b>>7)&1)
	}
	return nil
}

func (m *Mask) Encode() []byte {
	bb := make([]byte, len(m.mm))

	m.mx.RLock()
	defer m.mx.RUnlock()
	copy(bb, m.mm)
	return bb
}
