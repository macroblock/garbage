package conio

import (
	"strconv"
	"time"
)

// IEvent -
type IEvent interface {
	Time() time.Time
	Type() string
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
	ch  rune
	key int
}

// TWindowEvent -
type TWindowEvent struct {
	TEvent
	width  int
	height int
}

// Type -
func (ev *TKeyboardEvent) Type() string {
	return "keyboard"
}

// String -
func (ev *TKeyboardEvent) String() string {
	return ev.time.Format("15:04:05.000") + ": " + ev.Type() + " - " + strconv.Itoa(ev.key) + " " + string(ev.ch)
}

// Rune -
func (ev *TKeyboardEvent) Rune() rune {
	return ev.ch
}

// Key -
func (ev *TKeyboardEvent) Key() int {
	return ev.key
}

// Type -
func (ev *TWindowEvent) Type() string {
	return "window"
}

// String -
func (ev *TWindowEvent) String() string {
	return ev.time.Format("15:04:05.000") + ": " + ev.Type() + " - " + strconv.Itoa(ev.width) + "," + strconv.Itoa(ev.height)
}

// Size -
func (ev *TWindowEvent) Size() (int, int) {
	return ev.width, ev.height
}
