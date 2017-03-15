package conio

import termbox "github.com/nsf/termbox-go"

// Screen -
func Screen() *TScreen {
	return screenInstance // !!! TODO - multiple screens support
}

// EventStream -
func EventStream() *TEventStream {
	return eventStreamInstance // !!! TODO - multiple eventStream support
}

// Init - initializes conio
func Init() error {
	err := termbox.Init()
	return err
}

// Close - closes conio
func Close() {
	if iBuff != nil {
		close(iBuff)
	}
	if isReadEventsStarted {
		stopReadEvents()
	}
	if iBuff != nil {
		iBuff = nil
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Close()
}

func init() {
	initBorder()
}
