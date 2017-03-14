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

	scr.ShowCursor(false)
	//counter := 0
	var key rune
	width := scr.Width()
	height := scr.Height()
	for key != 27 {
		scr.Clear(' ', conio.ColorBlue, conio.ColorBlack)
		scr.DrawBorder(0, 0, width, height, conio.Border.Get("Double"), conio.ColorBlue, conio.ColorBlack)
		scr.DrawBorder(9, 9, width-18, height-18, conio.Border.Get("Single"), conio.ColorYellow, conio.ColorBlack)
		scr.DrawString(2, 2, "Тестовая string", conio.ColorGreen, conio.ColorDefault)
		scr.DrawString(3, 5, fmt.Sprintf("key: '%c' code: %d", key, key), conio.ColorYellow, conio.ColorDefault)
		scr.DrawString(3, 7, fmt.Sprintf("w: '%d' h: %d", width, height), conio.ColorYellow, conio.ColorDefault)

		scr.FillRect(10, 12, width-20, 1, ' ', conio.ColorDefault, conio.ColorDefault)
		scr.DrawAlignedString(10, 12, width-20, "\"Тестовая string\"", conio.AlignLeft, conio.ColorDefault, conio.ColorDefault)
		scr.FillRect(10, 14, width-20, 1, ' ', conio.ColorDefault, conio.ColorDefault)
		scr.DrawAlignedString(10, 14, width-20, "\"Тестовая string\"", conio.AlignRight, conio.ColorDefault, conio.ColorDefault)
		scr.FillRect(10, 16, width-20, 1, ' ', conio.ColorDefault, conio.ColorDefault)
		scr.DrawAlignedString(10, 16, width-20, "\"Тестовая string\"", conio.AlignCenter, conio.ColorDefault, conio.ColorDefault)
		//scr.MoveCursor(-10, -10)
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
