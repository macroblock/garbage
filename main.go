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
	conio.NewKeyboardAction("NextItem", "j", "", func(ev conio.TKeyboardEvent) bool {
		index++
		index = utils.Min(index, numItems-1)
		return true
	})
	conio.NewKeyboardAction("PrevItem", "k", "", func(ev conio.TKeyboardEvent) bool {
		index--
		index = utils.Max(index, 0)
		return true
	})
	conio.NewKeyboardAction("ShowActionWindow", "a", "", func(ev conio.TKeyboardEvent) bool {
		conio.ActionMap.SetMode("actions")
		showActions = true
		actionNames = conio.ActionMap.Names()
		sort.Strings(actionNames)
		return true
	})
	conio.NewKeyboardAction("HideActionWindow", "actions/a", "", func(ev conio.TKeyboardEvent) bool {
		conio.ActionMap.SetMode("")
		showActions = false
		return true
	})
	conio.NewKeyboardAction("ActWinNextItem", "actions/j", "", func(ev conio.TKeyboardEvent) bool {
		actionIndex++
		actionIndex = utils.Min(actionIndex, len(actionNames)-1)
		return true
	})
	conio.NewKeyboardAction("ActWinPrevItem", "actions/k", "", func(ev conio.TKeyboardEvent) bool {
		actionIndex--
		actionIndex = utils.Max(actionIndex, 0)
		return true
	})

	conio.ActionMap.Apply()
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
	scr.DrawBorder(width-logWidth-2, 0, logWidth+2, height)
	scr.FillRect(width-logWidth-1, 1, logWidth, height-2, ' ')

	scr.SetAlignment(conio.AlignCenter)
	scr.DrawAlignedString(width-logWidth-1, 0, logWidth, "[ Event log ]")
	scr.SetAlignment(conio.AlignLeft)
	for i := len(evLog) - 1; i >= 0 && i-len(evLog)+height-1 >= 1; i-- {
		scr.DrawAlignedString(width-logWidth-1, i-len(evLog)+height-1, logWidth, evLog[i])
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

	if showActions {
		posX += 10
		posY += 5
		w := 50
		scr.SetColor(conio.ColorWhite, conio.ColorBlack)
		scr.DrawBorder(posX-1, posY-1, w+2, len(actionNames)+2)
		scr.FillRect(posX, posY, w, len(actionNames), ' ')
		scr.SetAlignment(conio.AlignCenter)
		scr.DrawAlignedString(posX, posY-1, w, "[ Actions ]")

		scr.SetAlignment(conio.AlignLeft)
		for i, name := range actionNames {
			scr.SetColor(conio.ColorWhite, conio.ColorBlack)
			if i == actionIndex {
				scr.InvertColor()
				scr.FillRect(posX, posY+i, w, 1, ' ')
			}
			scr.DrawAlignedString(posX+offsX, posY+i, w, name)
		}
		posX -= 10
		posY -= 5
	}

	scr.Flush()
}

func handleEvent(ev conio.IEvent) {
	evLog = append(evLog, fmt.Sprintf("%s %T", ev.String(), ev))
	conio.HandleEvent(ev)
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
