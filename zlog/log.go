package zlog

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Log -
type Log struct {
	state   Filter
	loggers []ILogger
	name    string
}

// NewLog -
func NewLog(name string) *Log {
	return &Log{name: name}
}

// Instance -
func Instance() *Log {
	return defaultLog
}

// Close -
func Close() error {
	fmt.Println("  ")
	return nil
}

// Add -
func (o *Log) Add(loggers ...ILogger) {
	for _, v := range loggers {
		if v != nil {
			o.loggers = append(o.loggers, v)
		}
	}
}

// Log -
func (o *Log) Log(level Level, cause interface{}, message ...interface{}) {
	c, ok := getCause(cause)
	if !ok {
		return
	}
	msg := fmt.Sprint(message...)
	o.log(level, c, msg)
}

// Logf -
func (o *Log) Logf(level Level, cause interface{}, format string, message ...interface{}) {
	c, ok := getCause(cause)
	if !ok {
		return
	}
	msg := fmt.Sprintf(format, message...)
	o.log(level, c, msg)
}

func (o *Log) log(level Level, cause string, msg string) {
	// if len(msg) == 0 && len(cause) == 0 {
	// 	return
	// }

	if len(msg) == 0 {
		msg = cause
		cause = "" //cause[:0]
	}

	info := LogInfo{
		Level:   level,
		Cause:   cause,
		Message: msg,
		Time:    time.Now(),
		State:   o.state,
		LogName: o.name,
	}
	fillCallInfo(&info)

	// o.state &^= reset
	o.state |= level.Only()

	for _, logger := range o.loggers {
		if !logger.Admit(level.Only(), info.PackageName) {
			continue
		}
		msg := logger.Format(&info)

		if _, err := logger.Write([]byte(msg)); err != nil {
			// TODO: smarter
			fmt.Println(err)
		}
	}

	if level == LevelPanic {
		panic(msg)
	}
}

// String -
// func (o *Log) String() string {
// 	sl := []string{}
// 	for _, l := range o.node.loggers {
// 		sl = append(sl, l.LevelFilter().String()) //+": "+strings.Join(l.prefixes, ","))
// 	}
// 	return strings.Join(sl, "\n")
// }

func getCause(cause interface{}) (string, bool) {
	ret := ""
	ok := true
	switch v := cause.(type) {
	default:
		// slices
		if reflect.TypeOf(cause).Kind() == reflect.Slice {
			v := reflect.ValueOf(cause)
			if v.Len() == 0 {
				return "", false
			}
			slice := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				slice = append(slice, fmt.Sprint(v.Index(i)))
			}
			ret = strings.Join(slice, "\n")
			return ret, true
		}
	case nil:
		return "", false
	case bool:
		ok = v
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		if v == 0 {
			return "", false
		}
	case error:
		ret = v.Error()
		if len(ret) == 0 {
			return "", false
		}
	case string:
		if len(v) == 0 {
			return "", false
		}
		ret = v
	}
	return ret, ok
}
