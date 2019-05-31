package application

import (
	"fmt"

	"github.com/macroblock/garbage/sdf/channel"
	"github.com/macroblock/garbage/sdf/msg"
	"github.com/macroblock/garbage/sdf/rake"
	"github.com/veandco/go-sdl2/sdl"
)

type (
	// SysWindow -
	SysWindow struct {
		window   *sdl.Window
		renderer *sdl.Renderer
		cmdch    chan<- interface{}

		id   int
		name string
	}
)

// NewSysWindow -
func NewSysWindow(cmdMuxer *channel.Muxer) (rake.ISysWindow, error) {

	type retval struct {
		val rake.ISysWindow
		err error
	}

	retch := make(chan retval)
	cmdch := make(chan interface{})

	fn := func() error {
		flags := uint32(0)
		window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 480, flags)
		if err != nil {
			retch <- retval{nil, err}
			return err
		}

		flags = sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC
		renderer, err := sdl.CreateRenderer(window, -1, flags)
		_ = renderer
		if err != nil {
			window.Destroy()
			retch <- retval{nil, err}
			return err
		}
		id, err := window.GetID()
		if err != nil {
			renderer.Destroy()
			window.Destroy()
			retch <- retval{nil, err}
			return err
		}

		o := &SysWindow{}
		o.cmdch = cmdch
		o.window = window
		o.renderer = renderer
		o.id = int(id)
		retch <- retval{o, err}
		close(retch)
		return nil
	}

	cmdMuxer.Add(cmdch)
	cmdch <- fn
	ret := <-retch
	return ret.val, ret.err
}

// ID -
func (o *SysWindow) ID() int {
	return o.id
}

// Name -
func (o *SysWindow) Name() string {
	return o.name
}

// SendMessage -
func (o *SysWindow) SendMessage(messages ...rake.IMessage) {

}

// Close -
func (o *SysWindow) Close() {
	o.cmdch <- msg.Fn(fmt.Sprintf("close system window %q", o.name), func() error {
		err := o.renderer.Destroy()
		er2 := o.window.Destroy()
		if err == nil {
			err = er2
		}
		close(o.cmdch)
		return err
	})
}

func newSysWindow() (rake.ISysWindow, error) {
	flags := uint32(0)
	window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 480, flags)
	if err != nil {
		return nil, err
	}

	flags = sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC
	renderer, err := sdl.CreateRenderer(window, -1, flags)
	_ = renderer
	if err != nil {
		window.Destroy()
		return nil, err
	}
	id, err := window.GetID()
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		return nil, err
	}

	ret := &SysWindow{}
	ret.cmdch = make(chan interface{})
	ret.window = window
	ret.renderer = renderer
	ret.id = int(id)
	return ret, nil
}
