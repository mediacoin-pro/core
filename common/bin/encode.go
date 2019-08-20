package bin

import (
	"io"
	"os"
)

type Encoder interface {
	Encode() []byte
}

type Decoder interface {
	Decode([]byte) error
}

func Encode(vv ...interface{}) []byte {
	w := NewBuffer(nil)
	if err := w.WriteVar(vv...); err != nil {
		panic(err)
	}
	return w.Bytes()
}

func Decode(data []byte, vv ...interface{}) error {
	return NewBuffer(data).ReadVar(vv...)
}

func Read(r io.Reader, v ...interface{}) error {
	return NewReader(r).ReadVar(v...)
}

func Write(w io.Writer, v ...interface{}) (int64, error) {
	buf := NewBuffer(nil, v...)
	return io.Copy(w, buf.Buffer())
}

func ReadFile(name string, v ...interface{}) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	if err = Read(f, v...); err != nil {
		return err
	}
	return f.Close()
}

func WriteFile(name string, v ...interface{}) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err = Write(f, v...); err != nil {
		return err
	}
	return f.Close()
}

type binaryEncoder interface {
	BinaryEncode(w io.Writer) error
}

type binaryDecoder interface {
	BinaryDecode(r io.Reader) error
}

type binWriter interface {
	BinWrite(writer *Writer)
}

type binReader interface {
	BinRead(reader *Reader)
}
