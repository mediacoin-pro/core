package bin

import "bytes"

type Buffer struct {
	buf *bytes.Buffer
	Reader
	Writer
}

func NewBuffer(p []byte) *Buffer {
	buf := bytes.NewBuffer(p)
	return &Buffer{
		buf,
		Reader{rd: buf},
		Writer{wr: buf},
	}
}

func (b *Buffer) Error() error {
	if b.Reader.err != nil {
		return b.Reader.err
	}
	return b.Writer.err
}

func (w *Buffer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Buffer) Buffer() *bytes.Buffer {
	return w.buf
}
