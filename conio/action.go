package conio

import "tui/utils"
import "fmt"

// IAction -
type IAction interface {
	Name() string
	HotKey() string
	Description() string
	Do(ev IEvent) bool
}

// TAbstractAction -
type TAbstractAction struct {
	name        string
	hotKey      string
	description string
}

// TAction -
type (
	TAction struct {
		TAbstractAction
		handler TEventHandler
	}
	// TEventHandler -
	TEventHandler func(ev IEvent) bool
)

// TKeyboardAction -
type (
	TKeyboardAction struct {
		TAbstractAction
		handler func(ev TKeyboardEvent) bool
	}
	// TKeyboardHandler -
	TKeyboardHandler func(ev TKeyboardEvent) bool
)

// Name -
func (act *TAbstractAction) Name() string {
	return act.name
}

// HotKey -
func (act *TAbstractAction) HotKey() string {
	return act.hotKey
}

// Description -
func (act *TAbstractAction) Description() string {
	return act.description
}

// Handler -
func (act *TAbstractAction) Handler() TEventHandler {
	return nil
}

// Do -
func (act *TAbstractAction) Do(ev IEvent) bool {
	utils.Assert(false, "abstract class call")
	return false
}

// Do -
func (act *TAction) Do(ev IEvent) bool {
	if act.handler == nil {
		return false
	}
	return act.handler(ev)
}

// Do -
func (act *TKeyboardAction) Do(ev IEvent) bool {
	if act.handler == nil {
		return false
	}
	nev, ok := ev.(*TKeyboardEvent)
	utils.Assert(ok, fmt.Sprintf("incompatible event <%T> in <%T> method", ev, act))
	return act.handler(*nev)
}

// NewAction -
func NewAction(handler func(ev IEvent) bool) *TAction {
	action := &TAction{}
	action.handler = handler
	return action
}

// NewKeyboardAction -
func NewKeyboardAction(handler func(ev TKeyboardEvent) bool) *TKeyboardAction {
	action := &TKeyboardAction{}
	action.handler = handler
	return action
}
