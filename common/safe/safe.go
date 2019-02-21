package safe

import (
	"runtime/debug"
	"time"

	"github.com/mediacoin-pro/core/common/xlog"
)

func Exec(fn func()) {
	defer func() { recover() }()
	defer Recover()
	fn()
}

func Loop(duration time.Duration, fn func()) {
	for {
		Exec(fn)
		time.Sleep(duration)
	}
}

func Recover() {
	defer func() { recover() }()
	if r := recover(); r != nil {
		ss := debug.Stack()
		xlog.Fatal.Print("Panic:\n", string(ss))
		xlog.Fatal.Printf("FATAL-ERR: %v", r)
	}
}

func RecoverAndReport() {
	defer func() { recover() }()
	if r := recover(); r != nil {
		xlog.Fatal.Print("Panic:\n", string(debug.Stack()))
		//SendReport(fnName, r)
		xlog.Fatal.Printf("FATAL-ERR: %v", r)
		panic(r)
	}
}

func TracePanic() {
	if r := recover(); r != nil {
		xlog.Fatal.Print("Panic:\n", string(debug.Stack()))
		xlog.Fatal.Printf("PANIC-ERROR: %v", r)
		panic(r)
	}
}
