package main

import (
	"fmt"
	"garbage/conio"
	"garbage/utils"
)

func main() {
	err := conio.Init()
	utils.Assert(err == nil, "conio init failed")
	defer conio.Close()
	evs := conio.NewEventStream()
	utils.Assert(evs != nil, "eventStream init failed")
	defer evs.Close()
	scr := conio.NewScreen()
	utils.Assert(evs != nil, "screen init failed")
	defer scr.Close()
	//counter := 0
	var key rune
	width := scr.Width()
	height := scr.Height()
	for key != 27 {
		scr.Clear('+', conio.ColorBlack, conio.ColorWhite)
		scr.FillRect(1, 1, width-2, height-2, '.', conio.ColorYellow, conio.ColorRed)
		scr.FillRect(10, 10, width-20, height-20, '.', conio.ColorBlue, conio.ColorYellow)
		scr.DrawString(2, 2, "Тестовая string", conio.ColorGreen, conio.ColorDefault)
		scr.DrawString(3, 5, fmt.Sprintf("key: '%c' code: %d", key, key), conio.ColorYellow, conio.ColorDefault)
		scr.DrawString(3, 7, fmt.Sprintf("w: '%d' h: %d", width, height), conio.ColorYellow, conio.ColorDefault)
		scr.Flush()
		ev := evs.ReadEvent()
		if kbdEvent, ok := ev.(conio.TKeyboardEvent); ok {
			key = kbdEvent.Key
		}
		if winEvent, ok := ev.(conio.TWindowEvent); ok {
			width = winEvent.Width
			height = winEvent.Height
			scr.Flush()
		}
		//fmt.Printf("'%c' - %d\n", key, key)
		// 	if kbd.KeyPressed() {
		// 		key = kbd.ReadKey()
		// 		fmt.Printf("\n'%c' - %d\n", key, key)
		// 		counter = 0
		// 	}
		// 	counter++
		// 	fmt.Printf("%d\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08", counter)
	}

}
