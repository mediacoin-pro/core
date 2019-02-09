package sys

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func RunFile(filename string) error {
	//return run("cmd", "/C", "start", escExecPath(filename))
	return run("cmd", "/C", escExecPath(filename))
}

func ShowInExplorer(filePath string) (err error) {
	if IsDir(filePath) {
		return run("explorer", SystemPath(filePath))
	}
	return run("explorer", "/select,", SystemPath(filePath))
}

// Link uses std window-cmd mklink
// see: mklink manual
//func Link(src, dst string) (err error) {
//	if IsDir(src) {
//		err = run("cmd", "/C", "mklink", "/J", SystemPath(dst), SystemPath(src))
//	} else {
//		err = run("cmd", "/C", "mklink", "/H", SystemPath(dst), SystemPath(src)) // hard-link
//		if err != nil {
//			err = run("cmd", "/C", "mklink", SystemPath(dst), SystemPath(src)) // sym-link
//		}
//	}
//	return
//}

var replExecPath = strings.NewReplacer(
	"/", `\`,
	" ", "^ ",
	"(", "^(",
	")", "^)",
	"$", "^$",
	"%", "^%",
	"@", "^@",
	"&", "^&",
	"*", "^*",
	"#", "^#",
	"!", "^!",
)

func escExecPath(s string) string {
	s = replExecPath.Replace(s)
	if strings.HasPrefix(s, `\`) {
		s = homeDrive() + s
	}
	return s
}

func homeDrive() string {
	return os.Getenv("HOMEDRIVE")
}

var replSystemPath = strings.NewReplacer("/", `\`)

func SystemPath(s string) string {
	return replSystemPath.Replace(NormPath(s))
}

var reSlash = strings.NewReplacer("\\", "/")

func NormPath(filePath string) string {
	if filePath == "" {
		return ""
	}
	filePath = reSlash.Replace(filePath)
	filePath = path.Clean(filePath)
	if strings.HasPrefix(filePath, "/") {
		filePath = homeDrive() + filePath
	}
	filePath = strings.TrimSuffix(filePath, "/")
	return filePath
}

func run(name string, arg ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd-run-error: %v (%s)", err, stderr.String())
	}
	return nil
}

//func makeShortcut(src, dst string) (err error) {
//	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
//	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
//	if err != nil {
//		return
//	}
//	defer oleShellObject.Release()
//	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
//	if err != nil {
//		return
//	}
//	defer wshell.Release()
//	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
//	if err != nil {
//		return
//	}
//	idispatch := cs.ToIDispatch()
//	oleutil.PutProperty(idispatch, "TargetPath", src)
//	oleutil.CallMethod(idispatch, "Save")
//	return
//}
