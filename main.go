package main

import (
	"fmt"
	"garbage/conio"
	"garbage/utils"
	"sort"
)

var (
	key           rune
	width, height int
	canClose      bool
)

var (
	posX        = 5
	posY        = 5
	offsX       = 4
	index       = 0
	numItems    = -1
	borderNames []string
)

func initialize() {
	borderNames = conio.BorderMap.Names()
	numItems = len(borderNames)

	sort.Strings(borderNames)
}

func draw() {
	scr := conio.Screen()
	borderFg := conio.ColorWhite
	borderBg := conio.ColorRed
	captionFg := borderFg
	captionBg := borderBg
	winFg := conio.ColorWhite
	winBg := borderBg
	textFg := winFg
	textBg := winBg

	scr.Clear(' ', conio.ColorBlue, conio.ColorBlack)
	scr.DrawString(1, 1, fmt.Sprintf("key: '%c' code: %d", key, key), conio.ColorYellow, conio.ColorDefault)
	scr.DrawString(1, 2, fmt.Sprintf("w: '%d' h: %d", width, height), conio.ColorYellow, conio.ColorDefault)

	scr.DrawBorder(posX-1, posY-1, width-posX*2+2, len(borderNames)+2, conio.BorderMap.Get(borderNames[index]), borderFg, borderBg)
	scr.FillRect(posX, posY, width-posX*2, len(borderNames), ' ', winFg, winBg)
	scr.DrawAlignedString(posX, posY-1, width-posX*2, "[ Select border type ]", conio.AlignCenter, captionFg, captionBg)

	for i, name := range borderNames {
		scr.DrawAlignedString(posX+offsX, posY+i, width-posX*2-offsX, name, conio.AlignLeft, textFg, textBg)
	}
	scr.DrawAlignedString(posX, posY+index, offsX, " => ", conio.AlignLeft, textFg, textBg)
	//scr.MoveCursor(-10, -10)
	scr.Flush()
}

func handleEvent(ev interface{}) {
	if kbdEvent, ok := ev.(conio.TKeyboardEvent); ok {
		key = kbdEvent.Key
		switch key {
		case conio.KeyEscape:
			canClose = true
		case conio.KeyUp:
			index--
			if index < 0 {
				index = 0
			}
		case conio.KeyDown:
			index++
			if index >= numItems {
				index = numItems - 1
			}
		case conio.KeyLeft, conio.KeyHome, conio.KeyPageUp:
			index = 0
		case conio.KeyRight, conio.KeyEnd, conio.KeyPageDown:
			index = numItems - 1
		} // end case

	}
	if winEvent, ok := ev.(conio.TWindowEvent); ok {
		width = winEvent.Width
		height = winEvent.Height
		conio.Screen().Flush()
	}
}

func main() {
	err := conio.Init()
	utils.Assert(err == nil, "conio init failed")
	defer conio.Close()
	evs := conio.NewEventStream()
	utils.Assert(evs != nil, "eventStream init failed")
	defer evs.Close()
	scr := conio.NewScreen()
	utils.Assert(scr != nil, "screen init failed")
	defer scr.Close()

	scr.ShowCursor(false)

	initialize()
	for !canClose {
		draw()
		ev := evs.ReadEvent()
		handleEvent(ev)
	}

}
