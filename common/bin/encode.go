package bin

import "io"

type Encoder interface {
	Encode() []byte
}

type Decoder interface {
	Decode([]byte) error
}

func Encode(vv ...interface{}) []byte {
	w := NewBuffer(nil)
	w.WriteVar(vv...)
	return w.Bytes()
}

func Decode(data []byte, vv ...interface{}) error {
	return NewBuffer(data).ReadVar(vv...)
}

func Read(r io.Reader, v ...interface{}) error {
	return NewReader(r).ReadVar(v...)
}

func Write(w io.Writer, v ...interface{}) error {
	buf := NewBuffer(nil, v...)
	_, err := io.Copy(w, buf.Buffer())
	return err
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
