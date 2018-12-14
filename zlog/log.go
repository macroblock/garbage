package zlog

import (
	"fmt"
	"time"
)

type node struct {
	state   Filter
	loggers []*Logger
}

// Log -
type Log struct {
	node *node
	name string
}

// NewLog -
func NewLog(name string) *Log {
	return &Log{name: name, node: &node{}}
}

// Instance -
func Instance(name string) *Log {
	ret := *defaultLog
	ret.name = name
	return &ret
}

// Add -
func (o *Log) Add(loggers ...*Logger) {
	for _, v := range loggers {
		if v != nil {
			o.node.loggers = append(o.node.loggers, v)
		}
	}
}

// Log -
func (o *Log) Log(level Level, resetFilter Filter, err error, text string) {
	if text == "" {
		if err == nil {
			return
		}
		text = err.Error()
		err = nil
	}

	formatParams := FormatParams{
		Time:       time.Now(),
		LogLevel:   level,
		Text:       text,
		Error:      err,
		State:      o.node.state,
		ModuleName: o.name,
	}

	o.node.state &^= resetFilter
	o.node.state |= level.Only()

	for _, logger := range o.node.loggers {
		if level.NotIn(logger.Filter()) {
			continue
		}

		// formatParams.Format = logger.Format()
		// msg := logger.Formatter(formatParams)
		msg := formatMessage(formatParams)
		if _, err := logger.Writer().Write([]byte(msg)); err != nil {
			// TODO: smarter
			fmt.Println(err)
		}
	}
	if level == levelPanic {
		panic(text)
	}
}

func formatMessage(fp *FormatParams) string {
	return fmt.Sprintf("%v %v %v\n    %v\n", fp.Time, fp.LogLevel, fp.Text, fp.Error)
}

// String -
// func (o *Log) String() string {
// 	sl := []string{}
// 	for _, l := range o.node.loggers {
// 		sl = append(sl, l.LevelFilter().String()) //+": "+strings.Join(l.prefixes, ","))
// 	}
// 	return strings.Join(sl, "\n")
// }
