package zlog

import (
	"io"
	"time"
)

// Logger -
type Logger struct {
	name         string
	writer       io.Writer
	styler       Styler
	filter       Filter
	moduleFilter []string
	format       string
}

// Styler -
type Styler func(key rune, params *FormatParams) (string, bool)

// FormatParams -
type FormatParams struct {
	Format     string
	Time       time.Time
	LogLevel   Level
	Text       string
	Error      error
	HasError   bool
	State      Filter
	ModuleName string
	FileName   string
	FuncName   string
	LineNumber int
}

// Filter -
func (o *Logger) Filter() Filter {
	return o.filter
}
