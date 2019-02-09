package sys

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

func RunFile(filename string) (err error) {
	return exec.Command("open", filename).Run()
}

func ShowInExplorer(filePath string) (err error) {
	return exec.Command("open", "-R", filePath).Run()
}

func Link(src, dst string) (err error) {
	if err = os.Link(src, dst); err != nil {
		err = os.Symlink(src, dst)
	}
	return
}

func SystemPath(s string) string {
	return s
}

var reSlash = strings.NewReplacer("\\", "/")

func NormPath(filePath string) string {
	if filePath == "" {
		return ""
	}
	filePath = reSlash.Replace(filePath)
	filePath = path.Clean(filePath)
	filePath = strings.TrimSuffix(filePath, "/")
	return filePath
}
