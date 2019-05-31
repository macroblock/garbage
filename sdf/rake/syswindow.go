package rake

// ISysWindow -
type ISysWindow interface {
	SendMessage(...IMessage)
	Name() string
	Close()
}
