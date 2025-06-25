package exit

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Now(msg string, args ...any) {
	OnErr(fmt.Errorf(msg, args...))
}

func OnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "Error:", err.Error())
	if os.Getenv("GO_BACKTRACE") == "1" {
		debug.PrintStack()
	}
	os.Exit(1)
}

func OnErr2[T any](v T, err error) T {
	OnErr(err)
	return v
}

func OnErr3[A, B any](a A, b B, err error) (v1 A, v2 B) {
	OnErr(err)
	return a, b
}
