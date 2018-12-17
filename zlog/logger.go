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
	style := GetStyle(info.Level)
	header := fmtHeader(info, style)
	body := fmtBody(info, style)
	footer := fmtFooter(info, style)
	return fmt.Sprintf("%v%v%v", header, body, footer)
}

func fmtBriefState(state Filter) string {
	ret := ""
	for i := levelTooLow + 1; i < levelTooHigh; i++ {
		r := "-"
		if state&i.Only() != 0 {
			r = i.String()[0:1]
		}
		ret += r
	}
	return ret
}

func fmtVerboseState(state Filter) string {
	ret := " x "
	if state&LevelPanic.Only() != 0 {
		ret = filterToStr[LevelPanic]
	}
	for i := levelTooLow + 2; i < levelTooHigh; i++ {
		r := " x "
		if state&i.Only() != 0 {
			r = filterToStr[i]
		}
		ret += "|" + r
	}
	return ret
}

func fmtHeader(info *LogInfo, style *LevelStyle) string {
	switch style.Header {
	default:
		fallthrough
	case showVerbose:
		return fmt.Sprintf("%s%v %s\n%v %v %v/%v:%v %s\n",
			style.C2, info.Time, style.C1,
			info.Level, info.FuncName, info.PackageFullName, info.FileName, info.LineNumber, style.C0)
	case showBrief:
		return fmt.Sprintf("%s%v %v %v %v/%v:%v %s\n",
			style.C2, info.Level.Only(), info.Time.Format("2006-01-02 15:04:05"),
			info.FuncName, info.PackageName, info.FileName, info.LineNumber, style.C0)
	case showEssential:
		return fmt.Sprintf("%s%v %v %v %s\n",
			style.C2, info.Level.Only(), info.Time.Format("2006-01-02 15:04:05"),
			info.PackageName, style.C0)
	case showNone:
		return ""
	}
}

func fmtMsg(msg string, indent int, prefix string, C0, C1 string) string {
	if len(msg) == 0 {
		return ""
	}
	if indent == 0 && len(prefix) == 0 {
		return fmt.Sprintf("%s%v%s\n", C0, msg, C1)
	}
	offs := "\n" + strings.Repeat(" ", indent)
	msg = strings.Join(strings.Split(msg, "\n"), offs)
	return fmt.Sprintf("%s%v%v%s\n", C0, prefix, msg, C1)
}

func fmtBody(info *LogInfo, style *LevelStyle) string {
	if style.Body == showNone {
		return ""
	}
	cR := misc.Color()
	C1 := style.C1
	if len(info.Cause) == 0 {
		C1 = cR
	}
	switch style.Body {
	default:
		fallthrough
	case showVerbose:
		return fmt.Sprintf("%v%v",
			fmtMsg(info.Message, 4, ">>> ", style.C0, C1),
			fmtMsg(info.Cause, 4, "  > ", style.C1, cR))
	case showBrief:
		return fmt.Sprintf("%v%v",
			fmtMsg(info.Message, 4, "    ", style.C0, C1),
			fmtMsg(info.Cause, 4, "    ", style.C1, cR))
	case showEssential:
		return fmt.Sprintf("%v%v",
			fmtMsg(info.Message, 0, "", style.C0, C1),
			fmtMsg(info.Cause, 0, "", style.C1, cR))
	case showNone:
		return ""
	}
}

func fmtFooter(info *LogInfo, style *LevelStyle) string {
	cR := misc.Color()
	switch style.Footer {
	default:
		fallthrough
	case showVerbose:
		n := 79 - int(1+4*(levelTooHigh-levelTooLow-1))
		return fmt.Sprintf(
			"%s%v[%v]%s\n",
			style.C2, strings.Repeat("-", n), fmtVerboseState(info.State), cR)
	case showBrief:
		n := (79 - int(2+(levelTooHigh-levelTooLow-1))) / 8
		return fmt.Sprintf("%s%v[%v]%s\n",
			style.C2, strings.Repeat("\t", n), fmtBriefState(info.State), cR)
	case showEssential:
		return fmt.Sprintf("%s\n", cR)
	case showNone:
		return ""
	}
}

func (o *Logger) Write(p []byte) (n int, err error) {
	if o.writer != nil {
		return o.writer.Write(p)
	}
	fmt.Println(string(p))
	return len(p), nil
}
