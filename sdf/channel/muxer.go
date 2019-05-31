package channel

import (
	"sync"
)

type (
	// Muxer -
	Muxer struct {
		wg     sync.WaitGroup
		master chan interface{}
	}
)

// NewMuxer -
func NewMuxer(master chan interface{}) *Muxer {
	o := &Muxer{master: master}
	return o
}

// Close -
func (o *Muxer) Close() {
	o.wg.Wait()
}

// Add -
func (o *Muxer) Add(ch <-chan interface{}) {
	o.wg.Add(1)
	go func(ch <-chan interface{}) {
		// Pump it.
		for x := range ch {
			o.master <- x
		}
		o.wg.Done()
	}(ch)
}
