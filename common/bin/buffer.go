package bin

import "bytes"

type Buffer struct {
	buf *bytes.Buffer
	Reader
	Writer
}

func NewBuffer(p []byte, v ...interface{}) *Buffer {
	buf := bytes.NewBuffer(p)
	rw := &Buffer{
		buf,
		Reader{rd: buf},
		Writer{wr: buf},
	}
	if len(v) > 0 {
		rw.WriteVar(v...)
	}
	return rw
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

func (w *Buffer) Read(p []byte) (int, error) {
	return w.buf.Read(p)
}

func (w *Buffer) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *Buffer) Close() error {
	return nil
}
