package conio

import (
	"garbage/utils"

	termbox "github.com/nsf/termbox-go"
)

var screenInstance *TScreen

// TScreen -
type TScreen struct {
	cursorX, cursorY int
	isCursorVisible  bool
}

// NewScreen -
func NewScreen() *TScreen {
	utils.Assert(termbox.IsInit, "conio is not initialized correctly")
	utils.Assert(screenInstance == nil, "only one screen instance can be present")
	screenInstance = &TScreen{}
	screenInstance.ShowCursor(true)
	screenInstance.MoveCursor(0, 0)
	return screenInstance
}

// Close -
func (scr *TScreen) Close() {
	utils.Assert(screenInstance != nil, "screen is not initialized")
	screenInstance = nil
}

// DrawString -
func (scr *TScreen) DrawString(x, y int, s string, fg, bg TColor) {
	i := 0
	for _, ch := range s {
		termbox.SetCell(x+i, y, ch, termbox.Attribute(fg.color), termbox.Attribute(bg.color))
		i++
	}
}

// FillRect -
func (scr *TScreen) FillRect(x, y, w, h int, ch rune, fg, bg TColor) {
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			termbox.SetCell(x+i, y+j, ch, termbox.Attribute(fg.color), termbox.Attribute(bg.color))
		}
	}
}

// Clear -
func (scr *TScreen) Clear(ch rune, fg, bg TColor) {
	w, h := termbox.Size()
	screenInstance.FillRect(0, 0, w, h, ch, fg, bg)
}

// Flush -
func (scr *TScreen) Flush() {
	termbox.Flush()
}

// MoveCursor -
func (scr *TScreen) MoveCursor(x, y int) {
	scr.cursorX, scr.cursorY = x, y
	if scr.isCursorVisible {
		termbox.SetCursor(x, y)
	}
}

// ShowCursor -
func (scr *TScreen) ShowCursor(enable bool) {
	scr.isCursorVisible = enable
	if enable {
		termbox.SetCursor(scr.cursorX, scr.cursorY)
	} else {
		termbox.SetCursor(-1, -1)
	}
}
