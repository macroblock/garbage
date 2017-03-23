package conio

type tActionMap struct {
	byName     map[string]IAction
	byEventKey map[string]IAction
	mode       string
}

// ActionMap -
var ActionMap tActionMap

func initActionMap() {
	ActionMap = tActionMap{}
	ActionMap.byName = map[string]IAction{}
	ActionMap.Apply()

}

// Add -
func (am *tActionMap) Add(action IAction) {
	if _, ok := am.byName[action.Name()]; ok {
		//!!! TODO - log warning
	}
	am.byName[action.Name()] = action
}

// Delete -
func (am *tActionMap) Delete(action IAction) {
	if _, ok := am.byName[action.Name()]; !ok {
		//!!! TODO - log warning
	}
	delete(am.byName, action.Name())
}

// Apply -
func (am *tActionMap) Apply() {
	am.byEventKey = map[string]IAction{}
	for _, act := range am.byName {
		if _, ok := am.byEventKey[act.EventKey()]; ok {
			//!!! TODO - log warning
		}
		am.byEventKey[act.EventKey()] = act
	}
}

// Apply -
func (am *tActionMap) Names() []string {
	names := make([]string, 0, len(am.byName))
	for k, act := range am.byName {
		names = append(names, act.EventKey()+" - "+k)
	}
	return names
}

// SetMode -
func (am *tActionMap) SetMode(mode string) {
	am.mode = mode
	if len(am.mode) > 0 {
		am.mode += "/"
	}
}

// HandleEvent -
func HandleEvent(ev IEvent) {
	if act, ok := ActionMap.byEventKey[ActionMap.mode+ev.EventKey()]; ok {
		act.Do(ev)
	}
}
