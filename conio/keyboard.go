package conio

import (
	"garbage/utils"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var (
	isReadEventsStarted = false
	iBuff               chan rune
)

// TKeyboard -
type TKeyboard struct {
}

// NewKeyboard -
func NewKeyboard() *TKeyboard {
	utils.Assert(termbox.IsInit, "conio is not initialized correctly")
	utils.Assert(iBuff == nil, "only one keyboard instance can be present")
	kbd := &TKeyboard{}
	iBuff = make(chan rune, 32)
	kbd.Start()
	return kbd
}

// Close -
func (kbd *TKeyboard) Close() {
	utils.Assert(iBuff != nil, "keyboard is not initialized")
	close(iBuff)
	kbd.Stop()
	iBuff = nil
}

// Stop -
func (kbd *TKeyboard) Stop() {
	utils.Assert(iBuff != nil, "keyboard not initialized")
	stopReadEvents()
}

// Start -
func (kbd *TKeyboard) Start() {
	utils.Assert(iBuff != nil, "keyboard is not initialized")
	startReadEvents()
}

// ReadKey - reads a key. Blocking
func (kbd *TKeyboard) ReadKey() rune {
	utils.Assert(iBuff != nil, "keyboard is not initialized")
	return <-iBuff
}

// KeyPressed - returns true if key was pressed
func (kbd *TKeyboard) KeyPressed() bool {
	utils.Assert(iBuff != nil, "keyboard is not initialized")
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
		case termbox.EventKey:
			{
				if ev.Ch == 0 {
					iBuff <- rune(ev.Key)
				} else {
					iBuff <- ev.Ch
				}
			}
		} // end switch
	}
	isReadEventsStarted = false
}
