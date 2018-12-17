package zlog

import "github.com/macroblock/imed/pkg/misc"

const (
	cReset = misc.ColorReset
	cBold  = misc.ColorBold
	cFaint = misc.ColorFaint

	cBlack   = misc.ColorBlack
	cRed     = misc.ColorRed
	cGreen   = misc.ColorGreen
	cYellow  = misc.ColorYellow
	cBlue    = misc.ColorBlue
	cMagenta = misc.ColorMagenta
	cCyan    = misc.ColorCyan
	cWhite   = misc.ColorWhite

	cBgBlack   = misc.ColorBgBlack
	cBgRed     = misc.ColorBgRed
	cBgGreen   = misc.ColorBgGreen
	cBgYellow  = misc.ColorBgYellow
	cBgBlue    = misc.ColorBgBlue
	cBgMagenta = misc.ColorBgMagenta
	cBgCyan    = misc.ColorBgCyan
	cBgWhite   = misc.ColorBgWhite
)

type (
	appearance int

	// LevelStyle -
	LevelStyle struct {
		C0, C1, C2 string
		Header     appearance
		Title      appearance
		Body       appearance
		Footer     appearance
	}
)

const (
	showNone appearance = iota
	showEssential
	showBrief
	showVerbose
)

var (
	defaultStyle = LevelStyle{
		C0:     misc.Color(),
		C1:     misc.Color(),
		C2:     misc.Color(),
		Header: showVerbose,
		Body:   showVerbose,
		Footer: showBrief,
	}

	levelToColor = []*LevelStyle{
		{ // panic
			C0:     misc.Color(cBgRed, cWhite, cBold),
			C1:     misc.Color(cBgRed, cWhite),
			C2:     misc.Color(cBlack, cBold),
			Header: showVerbose,
			Body:   showVerbose,
			Footer: showNone,
		},
		{ // critical
			C0:     misc.Color(cBgRed, cWhite, cBold),
			C1:     misc.Color(cBgRed, cWhite),
			C2:     misc.Color(cBlack, cBold),
			Header: showVerbose,
			Body:   showVerbose,
			Footer: showNone,
		},
		{ // error
			C0:     misc.Color(cRed, cBold),
			C1:     misc.Color(cRed),
			C2:     misc.Color(cBlack, cBold),
			Header: showBrief,
			Body:   showVerbose,
			Footer: showNone,
		},
		{ // warning
			C0:     misc.Color(cYellow, cBold),
			C1:     misc.Color(cYellow),
			C2:     misc.Color(cBlack, cBold),
			Header: showBrief,
			Body:   showVerbose,
			Footer: showNone,
		},
		{ // notice
			C0:     misc.Color(cGreen, cBold),
			C1:     misc.Color(cGreen),
			C2:     misc.Color(cBlack, cBold),
			Header: showEssential,
			Body:   showEssential,
			Footer: showNone,
		},
		{ // info
			C0:     misc.Color(cWhite),
			C1:     misc.Color(cWhite),
			C2:     misc.Color(cBlack, cBold),
			Header: showNone,
			Body:   showEssential,
			Footer: showNone,
		},
		{ // debug
			C0:     misc.Color(cBlack, cBold),
			C1:     misc.Color(cBlack, cBold),
			C2:     misc.Color(cBlack, cBold),
			Header: showVerbose,
			Body:   showVerbose,
			Footer: showNone,
		},
	}
)

// GetStyle -
func GetStyle(level Level) *LevelStyle {
	if level <= levelTooLow || level >= levelTooHigh {
		return &defaultStyle
	}
	// return &defaultStyle
	return levelToColor[level]
}
