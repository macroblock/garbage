package main

import (
	"fmt"
	"tui/conio"
	. "tui/utils"
)

func main() {
	err := conio.Init()
	Assert(err == nil, "conio init failed")
	defer conio.Close()
	kbd := conio.NewKeyboard()
	Assert(kbd != nil, "keyboard init failed")

	counter := 0
	var key rune
	for key != 27 {
		if kbd.KeyPressed() {
			key = kbd.ReadKey()
			fmt.Printf("\n'%c' - %d\n", key, key)
			counter = 0
		}
		counter++
		fmt.Printf("%d\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08", counter)
	}

	kbd2.Close()
	kbd.KeyPressed()
}
