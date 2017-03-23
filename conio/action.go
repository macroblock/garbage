package conio

import "tui/utils"
import "fmt"

// IAction -
type IAction interface {
	Name() string
	EventKey() string
	Description() string
	Do(ev IEvent) bool
}

// TAbstractAction -
type TAbstractAction struct {
	name        string
	eventKey    string
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
		handler TKeyboardHandler
	}
	// TKeyboardHandler -
	TKeyboardHandler func(ev TKeyboardEvent) bool
)

// Name -
func (act *TAbstractAction) Name() string {
	return act.name
}

// EventKey -
func (act *TAbstractAction) EventKey() string {
	return act.eventKey
}

// Description -
func (act *TAbstractAction) Description() string {
	return act.description
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
func NewAction(name, eventKey, descr string, handler TEventHandler) *TAction {
	act := &TAction{}
	act.name = name
	act.eventKey = eventKey
	act.description = descr
	act.handler = handler
	ActionMap.Add(act)
	return act
}

// NewKeyboardAction -
func NewKeyboardAction(name, eventKey, descr string, handler TKeyboardHandler) *TKeyboardAction {
	act := &TKeyboardAction{}
	act.name = name
	act.eventKey = eventKey
	act.description = descr
	act.handler = handler
	ActionMap.Add(act)
	return act
}
