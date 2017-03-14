package conio

import (
	"garbage/utils"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// TKeyboardEvent -
type TKeyboardEvent struct {
	Key rune
}

// TWindowEvent -
type TWindowEvent struct {
	Width  int
	Height int
}

var (
	isReadEventsStarted = false
	iBuff               chan interface{}
)

// TEventStream -
type TEventStream struct {
}

// NewEventStream -
func NewEventStream() *TEventStream {
	utils.Assert(termbox.IsInit, "conio is not initialized correctly")
	utils.Assert(iBuff == nil, "only one eventStream instance can be present")
	evs := &TEventStream{}
	iBuff = make(chan interface{}, 32)
	evs.Start()
	return evs
}

// Close -
func (evs *TEventStream) Close() {
	utils.Assert(iBuff != nil, "eventStream is not initialized")
	close(iBuff)
	evs.Stop()
	iBuff = nil
}

// Stop -
func (evs *TEventStream) Stop() {
	utils.Assert(iBuff != nil, "eventStream not initialized")
	stopReadEvents()
}

// Start -
func (evs *TEventStream) Start() {
	utils.Assert(iBuff != nil, "eventStream is not initialized")
	startReadEvents()
}

// ReadEvent - reads a key. Blocking
func (evs *TEventStream) ReadEvent() interface{} {
	utils.Assert(iBuff != nil, "eventStream is not initialized")
	return <-iBuff
}

// HasEvent - returns true if key was pressed
func (evs *TEventStream) HasEvent() bool {
	utils.Assert(iBuff != nil, "eventStream is not initialized")
	//time.Sleep(1)
	return len(iBuff) > 0
}

func startReadEvents() {
	utils.Assert(!isReadEventsStarted, "readEvents loop is already started")
	go readEvents()
}

func stopReadEvents() {
	utils.Assert(isReadEventsStarted, "readEvents loop is not started yet")
	for len(iBuff) > 0 {
		<-iBuff
	}
	termbox.Interrupt()
	for isReadEventsStarted {
		time.Sleep(1)
	}
}

func readEvents() {
	isReadEventsStarted = true
loop:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventInterrupt:
			break loop

		case termbox.EventResize:
			iBuff <- TWindowEvent{ev.Width, ev.Height}

		case termbox.EventKey:
			if ev.Ch == 0 {
				iBuff <- TKeyboardEvent{rune(ev.Key)}
			} else {
				iBuff <- TKeyboardEvent{ev.Ch}
			}
		} // end switch
	}
	isReadEventsStarted = false
}
