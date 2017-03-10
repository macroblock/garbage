package utils

import (
	"fmt"
	"runtime"
)

// Assert -
func Assert(ok bool, s string) {
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Sprintf("\n### assert ###\n%s:%d\n[error] %s", file, line, s))
	}
}
