package sys

import "io"

var DevNull *emptyStream

//----------------------------------------------
type emptyStream struct {
}

func (*emptyStream) Read(buf []byte) (int, error) {
	return len(buf), nil
}

func (*emptyStream) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func (*emptyStream) Close() error {
	return nil
}

func NewProgressWriter(w io.Writer, progress func(written int64) error) io.Writer {
	return &progressStream{w, 0, progress}
}

type progressStream struct {
	w  io.Writer
	n  int64
	fn func(written int64) error
}

func (s *progressStream) Write(buf []byte) (n int, err error) {
	n, err = s.w.Write(buf)
	s.n += int64(n)
	if err == nil {
		err = s.fn(s.n)
	}
	return
}

func (s *progressStream) Close() error {
	s.fn(s.n)
	if c, _ := s.w.(io.Closer); c != nil {
		return c.Close()
	}
	return nil
}

func Copy(w io.Writer, r io.Reader, progress func(written int64) error) (int64, error) {
	if progress != nil {
		w = NewProgressWriter(w, progress)
	}
	return io.Copy(w, r)
}
