package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/macroblock/garbage/tpool/tpool"
)

const (
	total    = 20
	maxDelay = 2000
)

var data []int
var jobChan chan tpool.IJob

func init() {
	data = make([]int, total)
	for i := range data {
		data[i] = rand.Intn(maxDelay)
	}
	jobChan = make(chan tpool.IJob)
	fmt.Printf("1 in %v\n", jobChan)
	go func() {
		for _, v := range data {
			jobChan <- &TJob{delay: v}
		}
		// fmt.Printf("2 in %v\n", dataChan)
		close(jobChan)
	}()
}

// TJob -
type TJob struct {
	delay int
}

// Run -
func (o *TJob) Run(quit <-chan interface{}) error {
	fmt.Printf("   delay %v\n", o.delay)
	time.Sleep(time.Duration(o.delay) * time.Millisecond)
	return nil
}

// OnChangeState -
func (o *TJob) OnChangeState(id int, state int) {
	switch state {
	case tpool.StateStartUp:
	case tpool.StateReady:
	}
}

func main() {
	fmt.Printf("start\n")

	pool := tpool.New(jobChan)
	pool.AddThreads(1)
	fmt.Printf("in %v\n", jobChan)

	pool.Wait()

	fmt.Printf("done.\n")
}
