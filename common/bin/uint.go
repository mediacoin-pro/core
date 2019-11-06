package bin

import "math"

func Uint64ToBytes(val uint64) []byte {
	return []byte{
		byte(val >> 56),
		byte(val >> 48),
		byte(val >> 40),
		byte(val >> 32),
		byte(val >> 24),
		byte(val >> 16),
		byte(val >> 8),
		byte(val),
	}
}

func Uint32ToBytes(val uint32) []byte {
	return []byte{
		byte(val >> 24),
		byte(val >> 16),
		byte(val >> 8),
		byte(val),
	}
}

func Uint16ToBytes(val uint16) []byte {
	return []byte{
		byte(val >> 8),
		byte(val),
	}
}

func Uint8ToBytes(val uint8) []byte {
	return []byte{
		byte(val),
	}
}

func Float64ToBytes(val float64) []byte {
	return Uint64ToBytes(math.Float64bits(val))
}

func Float32ToBytes(val float32) []byte {
	return Uint32ToBytes(math.Float32bits(val))
}

func BoolToBytes(val bool) []byte {
	if val {
		return []byte{1}
	}
	return []byte{0}
}

func BytesToUint16(b []byte) uint16 {
	if n := len(b); n < 2 {
		if n == 0 {
			return 0
		}
		b = append(make([]byte, 2-n, 2), b...)
	}
	return uint16(b[1]) |
		(uint16(b[0]) << 8)
}

func BytesToUint32(b []byte) uint32 {
	if n := len(b); n < 4 {
		if n == 0 {
			return 0
		}
		b = append(make([]byte, 4-n, 4), b...)
	}
	return uint32(b[3]) |
		(uint32(b[2]) << 8) |
		(uint32(b[1]) << 16) |
		(uint32(b[0]) << 24)
}

func BytesToUint64(b []byte) uint64 {
	if n := len(b); n < 8 {
		if n == 0 {
			return 0
		}
		b = append(make([]byte, 8-n, 8), b...)
	}
	return uint64(b[7]) |
		(uint64(b[6]) << 8) |
		(uint64(b[5]) << 16) |
		(uint64(b[4]) << 24) |
		(uint64(b[3]) << 32) |
		(uint64(b[2]) << 40) |
		(uint64(b[1]) << 48) |
		(uint64(b[0]) << 56)
}

func BytesToFloat64(b []byte) float64 {
	return math.Float64frombits(BytesToUint64(b))
}

func BytesToFloat32(b []byte) float32 {
	return math.Float32frombits(BytesToUint32(b))
}
