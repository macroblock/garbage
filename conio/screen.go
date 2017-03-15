package conio

import (
	"garbage/utils"

	termbox "github.com/nsf/termbox-go"
)

var screenInstance *TScreen

// TAlignment -
type TAlignment int

// AlignLeft -
const (
	AlignLeft = TAlignment(iota)
	AlignRight
	AlignCenter
)

// TScreen -
type TScreen struct {
	fgc, bgc         TColor
	cursorX, cursorY int
	isCursorVisible  bool
	alignment        TAlignment
}

// NewScreen -
func NewScreen() *TScreen {
	utils.Assert(termbox.IsInit, "conio is not initialized correctly")
	utils.Assert(screenInstance == nil, "only one screen instance can be present")
	screenInstance = &TScreen{}
	screenInstance.ShowCursor(true)
	screenInstance.MoveCursor(0, 0)
	screenInstance.SetColor(ColorDefault, ColorDefault)
	screenInstance.SetAlignment(AlignLeft)
	return screenInstance
}

// Close -
func (scr *TScreen) Close() {
	utils.Assert(screenInstance != nil, "screen is not initialized")
	screenInstance = nil
}

// FgColor -
func (scr *TScreen) FgColor() TColor {
	return scr.fgc
}

// BgColor -
func (scr *TScreen) BgColor() TColor {
	return scr.bgc
}

// SetColor -
func (scr *TScreen) SetColor(fg, bg TColor) {
	scr.fgc = fg
	scr.bgc = bg
}

// InvertColor -
func (scr *TScreen) InvertColor() {
	scr.fgc, scr.bgc = scr.bgc, scr.fgc
}

// SetFgColor -
func (scr *TScreen) SetFgColor(fg TColor) {
	scr.fgc = fg
}

// SetBgColor -
func (scr *TScreen) SetBgColor(bg TColor) {
	scr.bgc = bg
}

// SetAlignment -
func (scr *TScreen) SetAlignment(alignment TAlignment) {
	scr.alignment = alignment
}

// DrawRune -
func (scr *TScreen) DrawRune(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, termbox.Attribute(scr.fgc.color), termbox.Attribute(scr.bgc.color))
}

func (scr *TScreen) drawRunes(x, y int, runes []rune) {
	for i, ch := range runes {
		termbox.SetCell(x+i, y, ch, termbox.Attribute(scr.fgc.color), termbox.Attribute(scr.bgc.color))
	}
}

// DrawString -
func (scr *TScreen) DrawString(x, y int, str string) {
	runes := []rune(str)
	scr.drawRunes(x, y, runes)
}

// DrawAlignedString -
func (scr *TScreen) DrawAlignedString(x, y, w int, str string) {
	if w <= 0 {
		return
	}
	runes := []rune(str)
	len := len(runes)
	var i1, offs int
	i2 := utils.Min(w, len)
	ellPos := x + w - 1
	switch scr.alignment {
	case AlignRight:
		i1 = utils.Max(0, len-w)
		i2 = len
		offs = utils.Max(0, w-len)
		ellPos = x
	case AlignCenter:
		offs = utils.Max(0, (w-len)/2)
	} // end switch
	scr.drawRunes(x+offs, y, runes[i1:i2])
	if w < len {
		scr.DrawRune(ellPos, y, '…')
	}
}

// FillRect -
func (scr *TScreen) FillRect(x, y, w, h int, ch rune) {
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			scr.DrawRune(x+i, y+j, ch)
		}
	}
}

// Clear -
func (scr *TScreen) Clear(ch rune, fg, bg TColor) {
	//termbox.Clear(termbox.Attribute(fg.color), termbox.Attribute(bg.color))
	scr.SetColor(fg, bg)
	scr.FillRect(0, 0, scr.Width(), scr.Height(), ch)
}

// DrawBorder -
func (scr *TScreen) DrawBorder(x, y, w, h int, border TBorder) {
	if h > 0 {
		scr.FillRect(x+1, y, w-2, 1, border.H1)
		scr.FillRect(x+1, y+h-1, w-2, 1, border.H2)
	}
	if w > 0 {
		scr.FillRect(x, y+1, 1, h-2, border.V1)
		scr.FillRect(x+w-1, y+1, 1, h-2, border.V2)
	}
	if w > 1 && h > 1 {
		scr.DrawRune(x, y, border.LU)
		scr.DrawRune(x+w-1, y, border.RU)
		scr.DrawRune(x, y+h-1, border.LD)
		scr.DrawRune(x+w-1, y+h-1, border.RD)
	}
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

func moveCursor(x, y int) {
	w, h := termbox.Size()
	if x >= 0 && x < w && y >= 0 && y < h {
		termbox.SetCursor(x, y)
	}
}

// MoveCursor -
func (scr *TScreen) MoveCursor(x, y int) {
	scr.cursorX, scr.cursorY = x, y
	if scr.isCursorVisible {
		moveCursor(x, y)
	}
}

// ShowCursor -
func (scr *TScreen) ShowCursor(enable bool) {
	scr.isCursorVisible = enable
	if enable {
		moveCursor(scr.cursorX, scr.cursorY)
	} else {
		termbox.SetCursor(-1, -1)
	}
}
