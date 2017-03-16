package main

import (
	"fmt"
	"garbage/conio"
	"garbage/utils"
	"sort"
)

var (
	key           int
	ch            rune
	width, height int
	canClose      bool
)

var (
	posX        = 5
	posY        = 5
	offsX       = 1
	index       = 0
	numItems    = -1
	borderNames []string
	logWidth    = 40
	evLog       []string
)

func initialize() {
	borderNames = conio.BorderMap.Names()
	numItems = len(borderNames)

	evLog = make([]string, 0)

	sort.Strings(borderNames)

	conio.Screen().Flush()
	width = conio.Screen().Width()
	height = conio.Screen().Height()

	conio.Screen().ShowCursor(false)
}

func draw() {
	scr := conio.Screen()
	winFg := conio.ColorWhite
	winBg := conio.ColorRed
	logFg := conio.ColorDarkGray
	logBg := conio.ColorDarkGray

	scr.Clear(' ', conio.ColorBlue, conio.ColorBlack)
	scr.SetColor(conio.ColorYellow, conio.ColorDefault)
	scr.DrawString(1, 0, fmt.Sprintf("char: '%c' %d", ch, ch))
	scr.DrawString(1, 1, fmt.Sprintf(" key: '%c' %d", key, key))
	scr.DrawString(1, 2, fmt.Sprintf("w: '%d' h: %d", width, height))

	scr.SelectBorder(borderNames[index])
	scr.SetColor(logFg, logBg)
	scr.DrawBorder(width-logWidth-2, 0, logWidth+2, height)
	scr.FillRect(width-logWidth-1, 1, logWidth, height-2, ' ')

	for i := len(evLog) - 1; i >= 0 && i-len(evLog)+height-1 >= 1; i-- {
		scr.DrawAlignedString(width-logWidth, i-len(evLog)+height-1, logWidth, evLog[i])
	}

	scr.SetColor(winFg, winBg)
	scr.DrawBorder(posX-1, posY-1, width-logWidth+2, len(borderNames)+2)
	scr.FillRect(posX, posY, width-logWidth, len(borderNames), ' ')

	scr.SetAlignment(conio.AlignCenter)
	scr.DrawAlignedString(posX, posY-1, width-logWidth, "[ Select border type ]")

	scr.SetAlignment(conio.AlignLeft)
	for i, name := range borderNames {
		scr.SetColor(winFg, winBg)
		if i == index {
			scr.InvertColor()
			scr.FillRect(posX, posY+i, width-logWidth, 1, ' ')
		}
		scr.DrawAlignedString(posX+offsX, posY+i, width-logWidth-offsX, name)
	}
	//scr.SetAlignment(conio.AlignRight)
	//scr.DrawAlignedString(posX, posY+len(borderNames), width-posX*2, "[ Select border type ]")
	//scr.DrawAlignedString(posX, posY+index, width-posX*2, " => ")
	//scr.MoveCursor(-10, -10)
	scr.Flush()
}

func handleEvent(ev conio.IEvent) {
	evLog = append(evLog, ev.String())
	if kbdEvent, ok := ev.(*conio.TKeyboardEvent); ok {
		key = kbdEvent.Key()
		ch = kbdEvent.Rune()
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
	if winEvent, ok := ev.(*conio.TWindowEvent); ok {
		width, height = winEvent.Size()
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

	initialize()
	for !canClose {
		draw()
		ev := evs.ReadEvent()
		handleEvent(ev)
	}

}
