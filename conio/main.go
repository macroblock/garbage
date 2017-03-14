package conio

import termbox "github.com/nsf/termbox-go"

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
