package sys

import (
	"io"
	"os"
)

func MoveFile(old, new string, progress func(copied int64)) (err error) {
	if err = os.Rename(old, new); err == nil { // OK
		if progress != nil {
			progress(FileSize(new))
		}
		return
	}
	if err = CopyFile(old, new, progress); err == nil {
		err = os.Remove(old)
	}
	return
}

func CopyFile(old, new string, progress func(copied int64)) (err error) {
	f1, err := os.Open(old)
	if err != nil {
		return
	}
	defer func() {
		f1.Close()
	}()
	f2, err := os.OpenFile(new, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f2.Close()

	var w io.Writer
	if progress == nil {
		w = f2
	} else {
		w = NewProgressWriter(f2, progress)
	}
	_, err = io.Copy(w, f1)
	return
}
