package sys

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExt(t *testing.T) {
	path := "/abc.def/file.JPEG"

	ext := FileExt(path)

	assert.Equal(t, "jpeg", ext)
}

func TestFileSize(t *testing.T) {
	org, sym := makeTestSymlinkFile(123)
	defer os.Remove(org)
	defer os.Remove(sym)

	sizeOrg := FileSize(org)
	sizeSym := FileSize(sym)

	assert.Equal(t, int64(123), sizeOrg)
	assert.Equal(t, int64(123), sizeSym)
}

func TestFileExists(t *testing.T) {
	org, sym := makeTestSymlinkFile(1)
	defer os.Remove(org)
	defer os.Remove(sym)

	ok := FileExists(sym)

	assert.True(t, ok)
}

func TestFileExists_Fail(t *testing.T) {
	org, sym := makeTestSymlinkFile(1)
	defer os.Remove(sym)
	os.Remove(org)

	ok := FileExists(sym)

	assert.False(t, ok)
}

func makeTestSymlinkFile(size int64) (org, sym string) {
	org = TempFilename("systest")
	sym = org + ".symlink"
	if err := ioutil.WriteFile(org, make([]byte, size), 0666); err != nil {
		panic(err)
	}
	if err := os.Symlink(org, sym); err != nil {
		panic(err)
	}
	return
}

func TestCopyFile(t *testing.T) {
	const fileSize = 1234
	name1 := TempFilename("test-copy-file-test-1")
	name2 := TempFilename("test-copy-file-test-2")
	defer os.Remove(name1)
	defer os.Remove(name2)
	if err := ioutil.WriteFile(name1, make([]byte, fileSize), 0666); err != nil {
		panic(err)
	}

	var copied int64
	err := CopyFile(name1, name2, func(n int64) {
		copied = n
	})

	assert.NoError(t, err)
	assert.Equal(t, fileSize, int(copied))
}

func TestMoveFile(t *testing.T) {
	const fileSize = 1234
	name1 := TempFilename("test-copy-file-test-1")
	name2 := TempFilename("test-copy-file-test-2")
	defer os.Remove(name1)
	defer os.Remove(name2)
	if err := ioutil.WriteFile(name1, make([]byte, fileSize), 0666); err != nil {
		panic(err)
	}

	var copied int64
	err := MoveFile(name1, name2, func(n int64) {
		copied = n
	})

	assert.NoError(t, err)
	assert.True(t, FileExists(name2))
	assert.True(t, !FileExists(name1))
	assert.Equal(t, fileSize, int(copied))
	assert.Equal(t, fileSize, int(FileSize(name2)))
}
