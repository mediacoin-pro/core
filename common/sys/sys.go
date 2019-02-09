package sys

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func TempDir() string {
	return NormPath(os.TempDir())
}

func TempFilename(ext string) string {
	if ext != "" {
		ext = "." + ext
	}
	return fmt.Sprintf("%s/tmp%016x%s", TempDir(), rand.Uint64(), ext)
}

func FileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st != nil
}

var reFileExt = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)

// FileExt returns file extension  without dot.
// For ex: FileExt("//path/path/file.ZIP") -> "zip"
//         FileExt("//path/path/file")     -> ""
func FileExt(path string) string {
	if s := reFileExt.FindString(path); s != "" {
		return strings.ToLower(s[1:])
	}
	return ""
}

func SetFileExt(path, newExt string) string {
	if n := len(reFileExt.FindString(path)); n > 0 {
		path = path[:len(path)-n]
	}
	return path + "." + newExt
}

func FileSize(path string) int64 {
	st, err := os.Stat(path)
	for err == nil && (st.Mode()&os.ModeSymlink) != 0 { // is symlink
		if path, err = os.Readlink(path); err == nil {
			st, err = os.Stat(path)
		} else {
			return 0
		}
	}
	if err == nil && st != nil {
		return st.Size()
	}
	return 0
}

func DirAllFiles(dir string) (files []string, err error) {
	err = filepath.Walk(dir, func(filePath string, _ os.FileInfo, _ error) error {
		if !IsDir(filePath) {
			files = append(files, NormPath(filePath))
		}
		return nil
	})
	return
}

func Mkdir(name string, perm os.FileMode) (err error) {
	if err = os.Mkdir(name, perm); os.IsExist(err) {
		err = nil
	}
	return
}

func IsDir(path string) bool {
	st, _ := os.Stat(RealPath(path))
	return st != nil && st.IsDir()
}

func RealPath(path string) string {
	path, _ = filepath.EvalSymlinks(path)
	return NormPath(path)
}

// is not correct for Windows!
// deprecated
func IsLink(path string) bool {
	st, err := os.Stat(path)
	return err == nil && (st.Mode()&os.ModeSymlink) != 0
}

func UserHomeDir() (dir string) {
	for _, name := range [...]string{
		"HOME",                   // /Users/{username} *-nix
		"HOMEPATH",               // \Users\{username}
		"CSIDL_PROFILE",          // C:\Users\username
		"CSIDL_PERSONAL",         // C:\Documents and Settings\username\My Documents
		"CSIDL_MYDOCUMENTS",      // The virtual folder that represents the My Documents desktop item. This value is equivalent to CSIDL_PERSONAL.
		"CSIDL_DESKTOPDIRECTORY", // C:\Documents and Settings\username\Desktop
		"LOCALAPPDATA",           // C:\Users\{username}\AppData\Local
		"APPDATA",                // C:\Users\{username}\AppData\Roaming
		"CSIDL_APPDATA",          // C:\Users\{username}\AppData\Roaming
		"ProgramData",            // C:\ProgramData
		"CommonProgramFiles",     // C:\Program Files\Common Files
		"CSIDL_LOCAL_APPDATA",    //
		"CSIDL_PROGRAM_FILES",    // C:\Program Files
		"CD", // The current directory
	} {
		if dir = os.Getenv(name); dir != "" {
			if dir = NormPath(dir); dir != "" && FileExists(dir) && isAvailableOnWrite(dir) { // dir exists
				break
			}
		}
	}
	return
}

func isAvailableOnWrite(dir string) bool {
	testDir := NormPath(dir) + "/.mdc-test"
	defer os.Remove(testDir)

	return FileExists(testDir) || os.Mkdir(testDir, 0755) == nil && FileExists(testDir)
}

func DirSize(dirname string) (n int64) {
	FetchDir(dirname, func(info os.FileInfo) error {
		if info.IsDir() {
			n += DirSize(dirname + "/" + info.Name())
		} else {
			n += info.Size()
		}
		return nil
	})
	return
}

func FetchDir(dirname string, fn func(os.FileInfo) error) error {
	f, err := os.Open(dirname)
	if err != nil {
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, filename := range names {
		if info, err := os.Lstat(dirname + "/" + filename); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		} else if err = fn(info); err != nil {
			return err
		}
	}
	return nil
}

func ClearDir(dir string) (err error) {
	return FetchDir(dir, func(info os.FileInfo) error {
		return os.RemoveAll(dir + "/" + info.Name())
	})
}

func ShellCommand(shellScript string) *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = bytes.NewBuffer([]byte(shellScript))
	return cmd
}

func RunShellCommand(shellScript string, stdOut, stdErr io.Writer) (err error) {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = bytes.NewBuffer([]byte(shellScript))
	if stdOut == nil {
		stdOut = DevNull
	}
	cmd.Stdout = stdOut
	if stdErr == nil {
		stdErr = DevNull
	}
	cmd.Stderr = stdErr
	return cmd.Run()
}
