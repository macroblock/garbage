package zlog

import "strings"

// Level -
type Level int

// General loglevel flags
const (
	levelTooLow Level = -1 + iota
	LevelPanic
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
	levelTooHigh
)

var levelToStr = []string{"PANIC", "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO", "DEBUG", "UNSUPPORTED"}
var filterToStr = []string{"PNC", "CRT", "ERR", "WRN", "NTC", "INF", "DBG"}

// Filter -
type Filter uint

// General loglevel filters
const (
	FilterNone Filter = 1<<iota - 1
	FilterPanic
	FilterCritical
	FilterError
	FilterWarning
	FilterNotice
	FilterInfo
	FilterDebug
	FilterAll // Filter  = 1<<uint(levelTooHigh) - 1
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

// Include -
func (o Filter) Include(f Filter) Filter { return o | f }

// Exclude -
func (o Filter) Exclude(f Filter) Filter { return o &^ f }

// Intersect -
func (o Filter) Intersect(f Filter) Filter { return o & f }

// String -
func (o Filter) String() string {
	if o == FilterNone {
		return ""
	}
	sl := []string{}
	i := 0
	x := o & FilterAll
	for x != 0 {
		if x&1 != 0 {
			sl = append(sl, filterToStr[i])
		}
		i++
		x = x >> 1
	}
	if len(sl) == 0 {
		return levelToStr[levelTooHigh]
	}
	return strings.Join(sl, "|")
}
