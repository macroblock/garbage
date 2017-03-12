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
	kbd := conio.NewKeyboard()
	utils.Assert(kbd != nil, "keyboard init failed")
	defer kbd.Close()
	scr := conio.NewScreen()
	utils.Assert(kbd != nil, "keyboard init failed")
	defer scr.Close()
	//counter := 0
	var key rune
	for key != 27 {
		scr.Clear('+', conio.ColorBlack, conio.ColorWhite)
		scr.FillRect(1, 1, 3, 3, '#', conio.ColorYellow, conio.ColorBlue)
		scr.DrawString(2, 2, "Тестовая string", conio.ColorGreen, conio.ColorDefault)
		scr.DrawString(3, 5, fmt.Sprintf("key: '%c' code: %d", key, key), conio.ColorYellow, conio.ColorDefault)
		scr.Flush()
		key = kbd.ReadKey()
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
