package zlog

// Level -
type Level int

// General loglevel flags
const (
	levelTooLow Level = -1 + iota
	LevelPanic
	LevelCritical
	LevelError
	LevelWarning
	LevelReset
	LevelNotice
	LevelInfo
	LevelDebug
	levelTooHigh
)

var levelToStr = []string{"PNC", "CRT", "ERR", "WRN", "RESET", "NTC", "INF", "DBG", "UNSUPPORTED"}

// Filter -
type Filter uint

// General loglevel filters
const (
	All Filter = 1<<uint(levelTooHigh) - 1
)

// Only -
func (o Level) Only() Filter { return 1 << uint(o) }

// Below -
func (o Level) Below() Filter { return o.Only() - 1 }

// Above -
func (o Level) Above() Filter { return levelTooHigh.Below() &^ o.OrLower() }

// OrLower -
func (o Level) OrLower() Filter { return o.Only()<<1 - 1 }

// OrHigher -
func (o Level) OrHigher() Filter { return levelTooHigh.Below() &^ o.Below() }

// In -
func (o Level) In(f Filter) bool { return f&o.Only() != 0 }

// NotIn -
func (o Level) NotIn(f Filter) bool { return f&o.Only() == 0 }

// String -
func (o Level) String() string {
	if o <= levelTooLow || o >= levelTooHigh {
		o = levelTooHigh
	}
	return levelToStr[o]
}
