package msg

type (
	// Message -
	Message struct {
		Sender   interface{}
		Receiver interface{}
	}
)

// SetSender -
func (o *Message) SetSender(sender interface{}) { o.Sender = sender }

// SetReceiver -
func (o *Message) SetReceiver(receiver interface{}) { o.Receiver = receiver }

// GetSender -
func (o *Message) GetSender() interface{} { return o.Sender }

// GetReceiver -
func (o *Message) GetReceiver() interface{} { return o.Receiver }
