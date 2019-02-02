package goldb

import "syscall"

func init() {
	// set limit open files in process
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{999999, 999999})
}
