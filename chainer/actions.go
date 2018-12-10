package chainer

type (
	// TAction -
	TAction struct {
		name        string
		keychain    string
		description string
		handler     TActionHandler
	}

	// TActionHandler -
	TActionHandler func(keychain string) bool

	// IAction -
	IAction interface {
		Name() string
		Keychain() string
		Description() string
		Handler() TActionHandler
		BinaryKey() ([]TBinaryKey, error)
	}

	// TBuilder -
	TBuilder struct{}
)

// NewAction -
func NewAction(name, keychain, desc string, handler TActionHandler) IAction {
	return &TAction{
		name:        name,
		keychain:    keychain,
		description: desc,
		handler:     handler,
	}
}

// Keychain -
func (o *TAction) Keychain() string { return o.keychain }

// Name -
func (o *TAction) Name() string { return o.name }

// Description -
func (o *TAction) Description() string { return o.description }

// Handler -
func (o *TAction) Handler() TActionHandler { return o.handler }

// BinaryKey -
func (o *TAction) BinaryKey() ([]TBinaryKey, error) { return keychainToBinary(o.keychain) }
