package conio

import termbox "github.com/nsf/termbox-go"

// TColor -
type TColor struct {
	color int32
}

// Colors -
var (
	ColorDefault  = TColor{int32(termbox.ColorDefault)}
	ColorBlack    = TColor{int32(termbox.ColorBlack)}
	ColorRed      = TColor{int32(termbox.ColorRed)}
	ColorGreen    = TColor{int32(termbox.ColorGreen)}
	ColorYellow   = TColor{int32(termbox.ColorYellow)}
	ColorBlue     = TColor{int32(termbox.ColorBlue)}
	ColorMagenta  = TColor{int32(termbox.ColorMagenta)}
	ColorCyan     = TColor{int32(termbox.ColorCyan)}
	ColorWhite    = TColor{int32(termbox.ColorWhite)}
	ColorDarkGray = TColor{int32(termbox.ColorBlack | termbox.AttrBold)}
)
