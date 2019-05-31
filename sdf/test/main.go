package main

import (
	"fmt"

	"github.com/macroblock/garbage/sdf/sdf"
)

func main() {
	app := sdf.Application()
	app.Run()
	app.NewWindow("test window")

	for err := range app.ErrorChannel() {
		fmt.Println("error: ", err)
	}
}
