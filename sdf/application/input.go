package application

import (
	"fmt"

	"github.com/macroblock/garbage/sdf/msg"
	"github.com/veandco/go-sdl2/sdl"
)

func (o *Application) processInput() {
	for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
		switch ev := ev.(type) {
		case *sdl.QuitEvent:
			o.exit = true
		case *sdl.KeyboardEvent:
			pressed := (ev.Type == sdl.KEYDOWN)
			if pressed {
				o.setError(fmt.Errorf("exit"))
				if ev.Keysym.Scancode == sdl.K_q {
					o.exit = true
				}
			}
		} // switch ev := event.(type)
	} // for PollEvent
	// time.Sleep(0)
	// o.setError(sdl.GetError())
}

func (o *Application) processCommands() {
	select {
	default:
	case cmd, ok := <-o.cmdch:
		if !ok {
			o.setError(fmt.Errorf("command channel is closed"))
			return
		}

		fmt.Printf("message %q %T", cmd, cmd)
		// o.errch <- fmt.Errorf("message: %q", cmd)

		err := error(nil)
		switch cmd := cmd.(type) {
		// case func() error:
		// 	o.setError(cmd())
		case *msg.TFn:
			err = cmd.Fn()
		case *msg.TNewSysWindow:
			fmt.Println("fn")
			val, err := cmd.Fn()
			fmt.Println("set")
			cmd.Result.Set(val, err)
			fmt.Println("set end")
		// case msg.TCloseSysWindow:
		// 	err = cmd.Window.CloseFn()(cmd.Window)
		default:
			o.setError(fmt.Errorf("unsupported command type %T", cmd))
		}
		if err != nil {
			o.setError(err)
		}
	}
}
