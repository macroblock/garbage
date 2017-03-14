package conio

import (
	"garbage/utils"

	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

var screenInstance *TScreen

// TAlignment -
type TAlignment int

const (
	AlignLeft = TAlignment(iota)
	AlignRight
	AlignCenter
)

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

// DrawStringAlign -
func (scr *TScreen) DrawAlignedString(x, y, w int, s string, alignment TAlignment, fg, bg TColor) {
	len := utf8.RuneCountInString(s)
	if len <= w {
		switch alignment {
		case AlignLeft:
		case AlignRight:
			x = x + w - len
		case AlignCenter:
			x = x + (w-len)/2
		} // end case
		scr.DrawString(x, y, s, fg, bg)
	} else {
		idx := 0
		ellPos := x + w - 1
		_ = ellPos
		switch alignment {
		case AlignLeft:
			len = w
		case AlignRight:
			idx = len - w
			len = w
			ellPos = x
		case AlignCenter:
			len = w
		} // end case
		if len > 0 {
			i := 0
			for _, ch := range s {
				if i >= idx+len {
					break
				}
				if i >= idx {
					termbox.SetCell(x+i-idx, y, ch, termbox.Attribute(fg.color), termbox.Attribute(bg.color))
				}
				i++
			}
			termbox.SetCell(ellPos, y, '…', termbox.Attribute(fg.color), termbox.Attribute(bg.color))
		}
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
	//termbox.Clear(termbox.Attribute(fg.color), termbox.Attribute(bg.color))
	screenInstance.FillRect(0, 0, scr.Width(), scr.Height(), ch, fg, bg)
}

// Width -
func (scr *TScreen) Width() int {
	w, _ := termbox.Size()
	return w
}

// Height -
func (scr *TScreen) Height() int {
	_, h := termbox.Size()
	return h
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
