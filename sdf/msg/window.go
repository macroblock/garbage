package msg

import (
	"fmt"

	"github.com/macroblock/garbage/sdf/rake"
)

type (
	// TNewSysWindow -
	TNewSysWindow struct {
		Message
		Result Result
		Name   string
		Fn     func() (rake.ISysWindow, error)
	}

	// TCloseSysWindow -
	TCloseSysWindow struct {
		Message
		Window rake.ISysWindow
	}
)

// NewSysWindow -
func NewSysWindow(name string, fn func() (rake.ISysWindow, error), result Result) rake.IMessage {
	return &TNewSysWindow{Result: result, Name: name, Fn: fn}
}

// CloseSysWindow -
func CloseSysWindow(window rake.ISysWindow) rake.IMessage {
	return &TCloseSysWindow{Window: window}
}

func (o *TNewSysWindow) String() string { return fmt.Sprintf("new system window: %q", o.Name) }
func (o *TCloseSysWindow) String() string {
	return fmt.Sprintf("close system window: %q", o.Window.Name())
}
