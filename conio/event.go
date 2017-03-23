package conio

import (
	"strconv"
	"time"
)

// IEvent -
type IEvent interface {
	Time() time.Time
	Type() string
	EventKey() string
	String() string
}

// TEvent -
type TEvent struct {
	time time.Time
}

// Time -
func (ev *TEvent) Time() time.Time {
	return ev.time
}

// TKeyboardEvent -
type TKeyboardEvent struct {
	TEvent
	key int
	mod int
	ch  rune
}

// TWindowEvent -
type TWindowEvent struct {
	TEvent
	width  int
	height int
}

// TMouseEvent -
type TMouseEvent struct {
	TEvent
	key int
	x   int
	y   int
	mod int
}

// EventKey -
func (ev *TKeyboardEvent) EventKey() string {
	return string(ev.ch)
}

// Type -
func (ev *TKeyboardEvent) Type() string {
	return "keyboard"
}

// String -
func (ev *TKeyboardEvent) String() string {
	return ev.time.Format("15:04:05.000") + ": " + ev.Type() + " - " + strconv.Itoa(ev.key) + " " + string(ev.ch) + " " + strconv.Itoa(ev.mod)
}

// Key -
func (ev *TKeyboardEvent) Key() int {
	return ev.key
}

// Mod -
func (ev *TKeyboardEvent) Mod() int {
	return ev.mod
}

// Rune -
func (ev *TKeyboardEvent) Rune() rune {
	return ev.ch
}

// EventKey -
func (ev *TWindowEvent) EventKey() string {
	return "<Resize>"
}

// Type -
func (ev *TWindowEvent) Type() string {
	return "resize"
}

// String -
func (ev *TWindowEvent) String() string {
	return ev.time.Format("15:04:05.000") + ": " + ev.Type() + " - " + strconv.Itoa(ev.width) + "," + strconv.Itoa(ev.height)
}

// Size -
func (ev *TWindowEvent) Size() (int, int) {
	return ev.width, ev.height
}

// EventKey -
func (ev *TMouseEvent) EventKey() string {
	return "<Mouse>"
}

// Type -
func (ev *TMouseEvent) Type() string {
	return "mouse"
}

// String -
func (ev *TMouseEvent) String() string {
	return ev.time.Format("15:04:05.000") + ": " + ev.Type() + " - " + strconv.Itoa(ev.x) + "," + strconv.Itoa(ev.y) + " " + strconv.Itoa(ev.mod) + " " + strconv.Itoa(ev.key)
}

// Key -
func (ev *TMouseEvent) Key() int {
	return ev.key
}

// Mod -
func (ev *TMouseEvent) Mod() int {
	return ev.mod
}

// X -
func (ev *TMouseEvent) X() int {
	return ev.x
}

// Y -
func (ev *TMouseEvent) Y() int {
	return ev.y
}
