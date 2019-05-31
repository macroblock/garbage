package rake

type (
	// IMessage -
	IMessage interface {
		String() string
		SetSender(interface{})
		SetReceiver(interface{})
		GetSender() interface{}
		GetReceiver() interface{}
	}
)
