package rake

// IApplication -
type IApplication interface {
	NewWindow(name string) (ISysWindow, error)
	Close()
	Run()
	ErrorChannel() <-chan error
}
