package msg

import (
	"fmt"

	"github.com/macroblock/garbage/sdf/rake"
)

type (
	// TFn -
	TFn struct {
		Message
		Name string
		Fn   func() error
	}
	// TSetTitle -
	TSetTitle struct {
		Message
		title string
	}

	// TSetVisiblity -
	TSetVisiblity struct {
		Message
		on bool
	}
)

// Fn -
func Fn(name string, fn func() error) rake.IMessage { return &TFn{Name: name, Fn: fn} }

// SetTitle -
func SetTitle(title string) rake.IMessage { return &TSetTitle{title: title} }

// SetVisiblity -
func SetVisiblity(on bool) rake.IMessage { return &TSetVisiblity{on: on} }

// Show -
func Show() rake.IMessage { return &TSetVisiblity{on: true} }

// Hide -
func Hide() rake.IMessage { return &TSetVisiblity{on: false} }

func (o *TFn) String() string           { return fmt.Sprintf("function: %q", o.Name) }
func (o *TSetTitle) String() string     { return fmt.Sprintf("set title: %q", o.title) }
func (o *TSetVisiblity) String() string { return fmt.Sprintf("set visiblity: %v", o.on) }
