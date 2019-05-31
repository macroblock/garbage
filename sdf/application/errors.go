package application

func (o *Application) setError(err error) {
	if err == nil {
		return
	}
	if o.err == nil {
		o.err = err
	}
	o.errch <- err
}
