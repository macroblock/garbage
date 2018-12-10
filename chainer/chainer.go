package chainer

import "fmt"

type (
	// TChainer -
	TChainer struct {
		tEntry
		current *tEntry
		history string
		actions map[string]IAction
	}

	tEntry struct {
		entries map[TBinaryKey]*tEntry
		action  IAction
	}

	// TBinaryKey -
	TBinaryKey int64
)

// String -
func (o TBinaryKey) String() string {
	return string(rune(o))
}

// String -
func (o tEntry) String() string {
	if o.entries == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", o.entries)
}

// NewChainer -
func NewChainer() *TChainer {
	return &TChainer{
		actions: map[string]IAction{},
	}
}

// Add -
func (o *TChainer) Add(actions ...IAction) *TChainer {
	for _, action := range actions {
		name := action.Name()
		a := o.actions[name]
		log.Warningf(a != nil, "TChainer.Add(): action %q has been overwritten", name)
		o.actions[name] = action
	}
	return o
}

// Remove -
func (o *TChainer) Remove(names ...string) {
	for _, name := range names {
		_, exists := o.actions[name]
		log.Warningf(!exists, "TChainer.Delete(): action %q not found", name)
		delete(o.actions, name)
	}
}

// Get -
func (o *TChainer) Get(name string) IAction {
	return o.actions[name]
}

// Extract -
func (o *TChainer) Extract(name string) IAction {
	action := o.Get(name)
	if action != nil {
		o.Remove(name)
	}
	return action
}

// Apply -
func (o *TChainer) Apply() {
	o.entries = map[TBinaryKey]*tEntry{}
	o.current = nil
	o.history = ""
	for _, action := range o.actions {
		cur := &o.tEntry
		data, err := keychainToBinary(action.Keychain())
		log.Errorf(err, "TChainer.Apply(): compile keychain %q", action.Keychain())
		keychain := ""
		for _, key := range data {
			keychain += key.String()
			entry := cur.entries[key]
			if entry != nil {
				cur = entry
				continue
			}
			// entry = &tEntry{entries: map[tBinaryKey]*tEntry{}}
			entry = &tEntry{}
			if cur.entries == nil {
				cur.entries = map[TBinaryKey]*tEntry{}
			}
			cur.entries[key] = entry
			cur = entry
		}
		log.Noticef("%v %v = %v", action.Name(), action.Keychain(), keychain)
		if cur.action != nil {
			log.Warningf(true, "TChainer.Add(): action %q with keychain %q has been overwritten by action %q", cur.action.Name(), cur.action.Keychain(), action.Name())
		}
		cur.action = action
	}
	fmt.Println(*o)
}

// Handle -
func (o *TChainer) Handle(key TBinaryKey) error {
	if o.current == nil {
		o.current = &o.tEntry
	}
	o.current = o.current.entries[key]
	if o.current == nil {
		o.history = o.history[:0]
		o.current = o.tEntry.entries[key]
		if o.current == nil {
			return nil
		}
	}
	o.history += key.String()
	if o.current.action == nil {
		return nil
	}
	handler := o.current.action.Handler()
	if handler == nil {
		return nil
	}
	ok := handler(o.history)
	if ok {
		// o.history = o.history[:0]
	}
	return nil
}
