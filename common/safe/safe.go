package safe

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/mediacoin-pro/core/common/sys"
	"github.com/mediacoin-pro/core/common/xlog"
)

func Go(fn func()) {
	go Exec(fn)
}

func Exec(fn func()) {
	defer Recover()
	fn()
}

func Func(fn func()) func() {
	return func() {
		defer Recover()
		fn()
	}
}

func OnceFunc(fn func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			defer Recover()
			fn()
		})
	}
}

func Loop(duration time.Duration, fn func()) {
	if duration == 0 {
		duration = time.Millisecond
	}
	for {
		Exec(fn)
		sys.Sleep(duration)
	}
}

func Watch(fn func() interface{}, event func()) {
	var _v = fn()
	go Loop(0, func() {
		if v := fn(); v != _v {
			_v = v
			Exec(event)
		}
	})
}

func Recover() {
	defer func() { recover() }()
	if r := recover(); r != nil {
		ss := debug.Stack()
		xlog.Fatal.Print("Panic:\n", string(ss))
		xlog.Fatal.Printf("FATAL-ERR: %v", r)
	}
}

func RecoverError(err *error) {
	defer func() { recover() }()
	if r := recover(); r != nil {
		if err != nil {
			*err = fmt.Errorf("PANIC-ERROR: %v", r)
		}
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
