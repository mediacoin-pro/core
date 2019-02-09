package sys

import (
	"io"
	"os"
)

type FileRange struct {
	file *os.File
	sec  *io.SectionReader
}

func OpenFileRange(path string, offset, size int64) (*FileRange, error) {
	if file, err := os.Open(path); err != nil {
		return nil, err
	} else {
		return &FileRange{
			file: file,
			sec:  io.NewSectionReader(file, offset, size),
		}, nil
	}
}

func (r *FileRange) Close() error {
	return r.file.Close()
}

func (r *FileRange) Read(b []byte) (int, error) {
	return r.sec.Read(b)
}

func (r *FileRange) ReadAt(p []byte, off int64) (n int, err error) {
	return r.sec.ReadAt(p, off)
}

func (r *FileRange) Seek(offset int64, whence int) (int64, error) {
	return r.sec.Seek(offset, whence)
}
