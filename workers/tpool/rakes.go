package tpool

// -
const (
	StateReady int = iota
	StateStartUp
	StateDone
)

type (
	// IOnChangeState -
	IOnChangeState interface {
		OnChangeState(id int, state int)
	}
	// IJob -
	IJob interface {
		Run(quit <-chan interface{}) error
	}
)
