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
		Message    appearance
		Cause      appearance
	}
)

const (
	showNone    appearance = iota
	showBrief   appearance = iota
	showVerbose appearance = iota
)

var (
	defaultStyle = LevelStyle{
		C0:      misc.Color(),
		C1:      misc.Color(),
		C2:      misc.Color(),
		Header:  showVerbose,
		Title:   showVerbose,
		Message: showVerbose,
		Cause:   showVerbose,
	}

	levelToColor = []*LevelStyle{
		{ // panic
			C0:      misc.Color(cBold, cBgRed, cWhite),
			C1:      misc.Color(cBgRed, cWhite),
			C2:      misc.Color(cBold, cBlack),
			Header:  showVerbose,
			Title:   showVerbose,
			Message: showVerbose,
			Cause:   showVerbose,
		},
		{ // critical
			C0:      misc.Color(cBold, cBgRed, cWhite),
			C1:      misc.Color(cBgRed, cWhite),
			C2:      misc.Color(cBold, cBlack),
			Header:  showVerbose,
			Title:   showVerbose,
			Message: showVerbose,
			Cause:   showVerbose,
		},
		{ // error
			C0:      misc.Color(cBold, cRed),
			C1:      misc.Color(cRed),
			C2:      misc.Color(cBold, cBlack),
			Header:  showVerbose,
			Title:   showNone,
			Message: showVerbose,
			Cause:   showVerbose,
		},
		{ // warning
			C0:      misc.Color(cBold, cYellow),
			C1:      misc.Color(cYellow),
			C2:      misc.Color(cBold, cBlack),
			Header:  showVerbose,
			Title:   showNone,
			Message: showVerbose,
			Cause:   showVerbose,
		},
		{ // notice
			C0:      misc.Color(cBold, cGreen),
			C1:      misc.Color(cGreen),
			C2:      misc.Color(cBold, cBlack),
			Header:  showVerbose,
			Title:   showNone,
			Message: showVerbose,
			Cause:   showVerbose,
		},
		{ // info
			C0:      misc.Color(cWhite),
			C1:      misc.Color(cWhite),
			C2:      misc.Color(cBold, cBlack),
			Header:  showNone,
			Title:   showNone,
			Message: showBrief,
			Cause:   showBrief,
		},
		{ // debug
			C0:      misc.Color(cBold, cBlack),
			C1:      misc.Color(cBold, cBlack),
			C2:      misc.Color(cBold, cBlack),
			Header:  showVerbose,
			Title:   showVerbose,
			Message: showVerbose,
			Cause:   showVerbose,
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
