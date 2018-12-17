package zlog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	defaultLog *Log

	defaultLogger = NewLogger()
)

type (
	// LogInfo -
	LogInfo struct {
		Level           Level
		Cause           string
		Message         string
		Time            time.Time
		State           Filter
		LogName         string
		PackageFullName string
		PackageName     string
		FileName        string
		FuncName        string
		LineNumber      int
	}
)

func init() {
	defaultLog = NewLog("global")
	defaultLog.Add(
		NewLogger(),
	)
}

func exitFunc() {
	os.Exit(-1)
}

func assert(ok bool, msg string) {
	if ok {
		return
	}
	println(msg)
	exitFunc()
}

func assertf(ok bool, format string, text ...interface{}) {
	assert(ok, fmt.Sprintf(format, text...))
}

func isTerminal(f *os.File) bool {
	if terminal.IsTerminal(int(f.Fd())) {
		return true
	}
	return false
}

func fillCallInfo(info *LogInfo) {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	name := runtime.FuncForPC(pc).Name()
	parts := strings.Split(name, ".")
	pl := len(parts)
	packageFullName := ""
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageFullName = strings.Join(parts[0:pl-2], ".")
		x := strings.Split(parts[pl-3], "/")
		packageName = x[len(x)-1]
	} else {
		packageFullName = strings.Join(parts[0:pl-1], ".")
		x := strings.Split(parts[pl-2], "/")
		packageName = x[len(x)-1]
	}

	// fmt.Printf("name: %v\nfile: %v\nline: %v\nfunc: %v\npkg: %v\n\n", name, fileName, line, funcName, packageName)
	// fmt.Printf("\n%v\n### %v *** %v\n", name, packageName, packageFullName)
	info.PackageFullName = packageFullName
	info.PackageName = packageName
	info.FileName = fileName
	info.FuncName = funcName
	info.LineNumber = line
}
