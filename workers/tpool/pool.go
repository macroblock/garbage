package tpool

import (
	"fmt"
	"sync"
)

// TPool -
type TPool struct {
	mt        sync.Mutex
	wg        sync.WaitGroup
	input     chan IJob
	threads   []<-chan interface{}
	fnOnEmpty func()
}

// New -
func New(in chan IJob) *TPool {
	o := &TPool{}
	o.mt = sync.Mutex{}
	o.wg = sync.WaitGroup{}
	o.input = in
	return o
}

// Wait -
func (o *TPool) Wait() {
	o.wg.Wait()
}

// AddThreads -
func (o *TPool) AddThreads(num uint) {
	for i := uint(0); i < num; i++ {
		o.AddThread()
		fmt.Printf("thread #%v added\n", i)
	}
}

func (o *TPool) reuseThread() int {
	o.mt.Lock()
	defer o.mt.Unlock()

	id := -1
	for i := range o.threads {
		if o.threads[i] == nil {
			id = i
			break
		}
	}
	if id < 0 {
		id = len(o.threads)
		o.threads = append(o.threads, nil)
	}
	o.threads[id] = make(chan interface{})
	return id
}

// AddThread -
func (o *TPool) AddThread() {
	o.wg.Add(1)
	id := o.reuseThread()
	quit := o.threads[id]

	go func() {
		defer func() {
			o.threads[id] = nil
			o.wg.Done()
		}()

		for {
			select {
			case job, ok := <-o.input:
				if !ok {
					fmt.Printf("id[%v] input channel was unexpectedly closed\n", id)
					return
				}
				fmt.Printf("#%v started\n", id)
				err := job.Run(quit)
				_ = err
				fmt.Printf("#%v done\n", id)

			case <-quit:
				fmt.Printf("id[%v] was terminated\n", id)
				return
				// default:
				// 	o.mt.Lock()
				// 	if o.fnOnEmpty != nil {
				// 		o.fnOnEmpty()
				// 	}
				// 	o.mt.Unlock()
			}
		}
	}()
}

// OnEmty -
func (o *TPool) OnEmty(fn func()) {
	o.mt.Lock()
	defer o.mt.Unlock()
	o.fnOnEmpty = fn
}
