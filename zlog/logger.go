package zlog

import (
	"fmt"
	"io"
	"os"
	"strings"

	ansi "github.com/k0kubun/go-ansi"
	"github.com/macroblock/imed/pkg/misc"
)

type (
	// Logger -
	Logger struct {
		name          string
		writer        io.Writer
		admitFunc     func(Filter, string) bool
		formatFunc    func(*LogInfo) string
		levelFilter   Filter
		packageFilter []string
	}

	// ILogger -
	ILogger interface {
		io.Writer
		Admit(Filter, string) bool
		Format(*LogInfo) string
	}
)

// NewLogger -
func NewLogger(params ...interface{}) ILogger {
	unsupported := "logger error: unsupported type %v.\n"
	dupInit := "logger error: multiple initialization of %v is not allowed.\n"
	overlap := "logger error: %v and %v cannot be initialized at once due same functionality.\n"
	errTypes := []string{}
	o := &Logger{
		levelFilter: FilterAll,
		writer:      ansi.NewAnsiStdout(),
	}
	first := true
	w, a, f := 0, 0, 0
	for _, v := range params {
		switch t := v.(type) {
		default:
			errTypes = append(errTypes, fmt.Sprintf("%T", t))
		case Level:
			if first {
				o.levelFilter = 0
				first = false
			}
			o.levelFilter = t.OrLower()
		case Filter:
			if first {
				o.levelFilter = 0
				first = false
			}
			o.levelFilter |= t
		case os.File:
			o.writer = &t
			w++
		case *os.File:
			o.writer = t
			w++
		case io.Writer:
			o.writer = t
			w++
		case string:
			o.packageFilter = append(o.packageFilter, t)
		case func(Filter, string) bool:
			o.admitFunc = t
			a++
		case func(*LogInfo) string:
			o.formatFunc = t
			f++
		}
		assertf(o.admitFunc == nil, overlap)
		assertf(len(o.packageFilter) == 0, overlap)
	}
	exit := false
	if len(errTypes) > 0 {
		for _, s := range errTypes {
			fmt.Printf(unsupported, s)
		}
		exit = true
	}
	if w > 1 {
		fmt.Printf(dupInit, "logger.writer")
		exit = true
	}
	if a > 1 {
		fmt.Printf(dupInit, "logger.admitFunc")
		exit = true
	}
	if f > 1 {
		fmt.Printf(dupInit, "logger.formatFunc")
		exit = true
	}
	if len(o.packageFilter) > 0 && a > 1 {
		fmt.Printf(overlap, "logger.packageFilter", "logger.admitFunc")
		exit = true
	}
	if !first && a > 1 {
		fmt.Printf(overlap, "logger.levelFilter", "logger.admitFunc")
		exit = true
	}
	if exit {
		exitFunc()
	}
	return o
}

// Admit -
func (o *Logger) Admit(filter Filter, packageName string) bool {
	if o.admitFunc != nil {
		return o.admitFunc(filter, packageName)
	}
	if o.levelFilter&filter == 0 {
		return false
	}
	if len(o.packageFilter) == 0 {
		return true
	}
	for _, s := range o.packageFilter {
		if s == packageName {
			return true
		}
	}
	return false
}

// Format -
func (o *Logger) Format(info *LogInfo) string {
	if o.formatFunc != nil {
		return o.formatFunc(info)
	}

	head := fmtHeader(info)
	title := fmtTitle(info)
	msg := fmtMessage(info)
	err := fmtCause(info)

	ret := fmt.Sprintf("%v%v%v%v", head, title, msg, err)

	return ret
}

func fmtHeader(info *LogInfo) string {
	style := GetStyle(info.Level)
	if style.Header == showNone {
		return ""
	}
	ret := fmt.Sprintf("%v\n%v [%v]%v", style.C2, info.Time, info.State, style.C1)
	return ret
}

func fmtTitle(info *LogInfo) string {
	style := GetStyle(info.Level)
	if style.Title == showNone {
		return ""
	}
	ret := fmt.Sprintf("%v\n%v %v.%v at %v:%v %v",
		"", info.Level, info.PackageName, info.FuncName, info.FileName, info.LineNumber, "")
	return ret
}

func fmtMessage(info *LogInfo) string {
	style := GetStyle(info.Level)
	if style.Message == showNone {
		return ""
	}
	x := strings.Split(info.Message, "\n")
	msg := strings.Join(x, "\n    ")
	ret := fmt.Sprintf("%v\n>>> %v%v", style.C0, msg, "")
	return ret
}

func fmtCause(info *LogInfo) string {
	style := GetStyle(info.Level)
	if style.Cause == showNone || len(info.Cause) == 0 {
		return ""
	}
	ret := info.Cause
	x := strings.Split(ret, "\n")
	ret = fmt.Sprintf("%v\n  > %v%v",
		style.C1, strings.Join(x, "\n    "), misc.Color())
	return ret
}

// Write -
func (o *Logger) Write(p []byte) (n int, err error) {
	if o.writer != nil {
		return o.writer.Write(p)
	}
	fmt.Println(string(p))
	return len(p), nil
}
