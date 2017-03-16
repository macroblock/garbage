package conio

import (
	"garbage/utils"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// TEventStream -
type TEventStream struct {
}

var (
	isReadEventsStarted = false
	iBuff               chan IEvent
	eventStreamInstance *TEventStream
)

// NewEventStream -
func NewEventStream() *TEventStream {
	utils.Assert(termbox.IsInit, "conio is not initialized correctly")
	utils.Assert(iBuff == nil, "only one eventStream instance can be present")
	eventStreamInstance := &TEventStream{}
	iBuff = make(chan IEvent, 32)
	eventStreamInstance.Start()
	return eventStreamInstance
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
func (evs *TEventStream) ReadEvent() IEvent {
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
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventInterrupt:
			break loop

		case termbox.EventResize:
			ev := TWindowEvent{}
			ev.time = time.Now()
			ev.width = event.Width
			ev.height = event.Height
			iBuff <- &ev

		case termbox.EventKey:
			ev := TKeyboardEvent{}
			ev.time = time.Now()
			ev.ch = event.Ch
			ev.key = int(event.Key)
			iBuff <- &ev
		}
	}
	isReadEventsStarted = false
}
