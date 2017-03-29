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
	logWidth    = 50
	evLog       []string
	showActions = false
	actionNames []string
	actionIndex = 0
)

func initialize() {
	borderNames = conio.BorderMap.Names()
	numItems = len(borderNames)
	sort.Strings(borderNames)

	evLog = make([]string, 0)

	conio.Screen().Flush()
	width = conio.Screen().Width()
	height = conio.Screen().Height()

	conio.Screen().ShowCursor(false)
	conio.Screen().EnableShadow(true)

	conio.NewKeyboardAction("Exit", "q", "", func(ev conio.TKeyboardEvent) bool {
		canClose = true
		return true
	})
	conio.NewKeyboardAction("Next item", "j", "", func(ev conio.TKeyboardEvent) bool {
		index++
		index = utils.Min(index, numItems-1)
		return true
	})
	conio.NewKeyboardAction("Prev item", "k", "", func(ev conio.TKeyboardEvent) bool {
		index--
		index = utils.Max(index, 0)
		return true
	})
	conio.NewKeyboardAction("actions/Show", "a", "", func(ev conio.TKeyboardEvent) bool {
		conio.ActionMap.SetMode("actions")
		showActions = true
		actionNames = conio.ActionMap.Names()
		sort.Strings(actionNames)
		return true
	})
	conio.NewKeyboardAction("actions/Hide", "actions/a", "", func(ev conio.TKeyboardEvent) bool {
		conio.ActionMap.SetMode("")
		showActions = false
		return true
	})
	conio.NewKeyboardAction("actions/Next item", "actions/j", "", func(ev conio.TKeyboardEvent) bool {
		actionIndex++
		actionIndex = utils.Min(actionIndex, len(actionNames)-1)
		return true
	})
	conio.NewKeyboardAction("actions/Prev item", "actions/k", "", func(ev conio.TKeyboardEvent) bool {
		actionIndex--
		actionIndex = utils.Max(actionIndex, 0)
		return true
	})

	conio.ActionMap.Apply()
}

func drawWindow(x, y, w, h int, title string, draw func(x, y, w, h int)) {
	scr := conio.Screen()
	scr.DrawBorder(x, y, w, h)
	scr.FillRect(x+1, y+1, w-2, h-2, ' ')
	scr.DrawAlignedString(x+1, y, w-2, title)
	draw(x+1, y+1, w-2, h-2)
}

func draw() {
	scr := conio.Screen()
	winFg := conio.ColorWhite
	winBg := conio.ColorRed
	logFg := conio.ColorWhite
	logBg := conio.ColorDarkGray

	scr.Clear('░', conio.ColorWhite, conio.ColorBlack)
	scr.SetColor(conio.ColorYellow, conio.ColorDefault)
	scr.DrawString(1, 0, fmt.Sprintf("char: '%c' %d", ch, ch))
	scr.DrawString(1, 1, fmt.Sprintf("key: '%c' %d", key, key))
	scr.DrawString(1, 2, fmt.Sprintf("w: '%d' h: %d", width, height))

	scr.SelectBorder(borderNames[index])

	scr.SetColor(logFg, logBg)
	drawWindow(width-logWidth-1, 0, logWidth, height, "[ Event Log ]", func(x, y, w, h int) {
		for i := 0; i < utils.Min(len(evLog)-1, h); i++ {
			scr.DrawAlignedString(x, y+h-1-i, w, evLog[len(evLog)-1-i])
		}
	})

	scr.SetColor(winFg, winBg)
	drawWindow(posX, posY, width-logWidth+2, len(borderNames)+2, "[ Select border type ]", func(x, y, w, h int) {
		for i, name := range borderNames {
			if i == index {
				scr.InvertColor()
				scr.FillRect(x, y+i, w, 1, ' ')
				scr.DrawAlignedString(x+offsX, y+i, w-offsX, name)
				scr.InvertColor()
				continue
			}
			scr.DrawAlignedString(x+offsX, y+i, w-offsX, name)
		}
	})

	if showActions {
		scr.SetColor(conio.ColorWhite, conio.ColorBlack)
		drawWindow(posX+10, posY+5, 50, len(actionNames)+2, "[ Actions ]", func(x, y, w, h int) {
			for i, name := range actionNames {
				if i == actionIndex {
					scr.InvertColor()
					scr.FillRect(x, y+i, w, 1, ' ')
					scr.DrawAlignedString(x+offsX, y+i, w-offsX, name)
					scr.InvertColor()
					continue
				}
				scr.DrawAlignedString(x+offsX, y+i, w-offsX, name)
			}
		})
	}

	scr.Flush()
}

func handleEvent(ev conio.IEvent) {
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
		evLog = append(evLog, fmt.Sprintf("%s %T", ev.String(), ev))
		conio.HandleEvent(ev)
		handleEvent(ev)
	}

}
