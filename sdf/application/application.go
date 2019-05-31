package application

import (
	"fmt"
	"runtime"

	"github.com/macroblock/garbage/sdf/msg"

	"github.com/macroblock/garbage/sdf/channel"
	"github.com/macroblock/garbage/sdf/rake"
	"github.com/veandco/go-sdl2/sdl"
)

type (
	// Application -
	Application struct {
		err      error
		exit     bool
		errch    chan error
		cmdch    chan interface{}
		cmdMuxer *channel.Muxer
		windows  []rake.ISysWindow
	}
)

// New -
func New() *Application {
	return &Application{}
}

// Close -
func (o *Application) Close() {
}

// NewWindow -
func (o *Application) NewWindow(name string) (rake.ISysWindow, error) {
	// win, err := NewSysWindow(o.cmdMuxer)
	// if err != nil {
	// 	o.setError(err)
	// 	return nil
	// }
	// o.windows = append(o.windows, win)
	// return win
	result := msg.NewResult()
	fmt.Println("new")
	o.cmdch <- msg.NewSysWindow(name, newSysWindow, result)
	fmt.Println("get")
	val, err := result.Get()
	fmt.Println("end")
	return val.(rake.ISysWindow), err
}

func (o *Application) cleanUp() {
	for _, v := range o.windows {
		v.Close()
	}
	o.windows = nil

	o.cmdMuxer.Close()
	close(o.cmdch)
	o.cmdch = nil

	close(o.errch)
	o.errch = nil
	o.err = nil
	o.exit = false
}

// Run -
func (o *Application) Run() {

	if o.errch != nil {
		panic("an attempt to run an application thread twice")
	}
	o.errch = make(chan error)
	o.cmdch = make(chan interface{})
	o.cmdMuxer = channel.NewMuxer(o.cmdch)

	go func() {
		runtime.LockOSThread()

		err := sdl.Init(sdl.INIT_EVERYTHING)
		defer sdl.Quit()
		defer o.cleanUp()
		o.setError(err)

		// lastUpdate := time.Now()
		// lastRender := time.Now()
		for o.err == nil && !o.exit {
			// fixedTime = time.Since(programStart)
			o.HandleEvents()
			// sdf.deltaUpdate = time.Since(lastUpdate)
			// lastUpdate = time.Now()
			o.Update()
			// sdf.deltaRender = time.Since(lastRender)
			// lastRender = time.Now()
			o.Render()
			// sdf.renderer.Present()
			// sdl.Delay(1000 / 60)
			// o.setError(fmt.Errorf("test error"))
			sdl.Delay(100)
		}
	}() // go func()
}

// ErrorChannel -
func (o *Application) ErrorChannel() <-chan error {
	if o.errch == nil {
		panic("error channel is nil")
	}
	return o.errch
}

// HandleEvents -
func (o *Application) HandleEvents() {
	o.processInput()
	o.processCommands()
}

// Update -
func (o *Application) Update() {
}

// Render -
func (o *Application) Render() {
}
