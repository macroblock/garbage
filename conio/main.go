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
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputAlt | termbox.InputMouse)
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
	x := TAction{}
	x.name = "test"
	x.hotKey = "q"
	x.description = "test"
	x.handler = func(ev IEvent) bool {
		return true
	}

	y := TKeyboardAction{}
	y.name = "test kbd"
	y.hotKey = "1"
	y.description = "test kbd"
	y.handler = func(ev TKeyboardEvent) bool { return true }

	_ = x
	_ = y
}
