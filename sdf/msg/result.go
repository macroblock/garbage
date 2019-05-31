package msg

type (
	res struct {
		val interface{}
		err error
	}
	// Result -
	Result struct {
		ch chan res
	}
)

// NewResult -
func NewResult() Result {
	return Result{ch: make(chan res)}
}

// Close -
func (o Result) Close() {
	_ = <-o.ch
}

// Set -
func (o Result) Set(val interface{}, err error) {
	o.ch <- res{val, err}
	close(o.ch)
}

// Get -
func (o Result) Get() (interface{}, error) {
	res, ok := <-o.ch
	if !ok {
		panic("attempt to read a close channel")
	}
	return res.val, res.err
}
