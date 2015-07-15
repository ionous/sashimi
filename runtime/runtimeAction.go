package runtime

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
)

//
// Event data for event handlers and actions.
//
type RuntimeAction struct {
	game   *Game
	action *M.ActionInfo
	objs   []*GameObject
	values map[string]TemplateValues
	after  []G.Callback
}

// queue for running after the default actions
func (act *RuntimeAction) runAfterDefaults(cb G.Callback) {
	act.after = append(act.after, cb)
}

// Run the passed game callback.
// This creates the game event adapter, sets up the necessary template context, etc.
// Returns the results of the callback.
func (act *RuntimeAction) runCallback(cb G.Callback, cls *M.ClassInfo) bool {
	// FIX: it might be cooler to make act some sort of API...
	// then we could have the callback object: callback.run( IPlay, Data ) maybe...
	adapter := NewGameAdapter(act.game)
	adapter.data = act
	adapter.hint = cls
	templateValueStack.pushValues(act.values)
	defer templateValueStack.pop()
	cb(adapter)
	return !adapter.cancelled
}

//
// Default actions occur after event processing assuming that they have not been cancelled.
//
func (act *RuntimeAction) runDefaultActions() {
	// FIX: assign defaults at initialization?
	// it'd be even better if act didn't need game --
	// the main reason it does it to share code b/t Go() and the ProcessEventLoop
	if actions, existed := act.game.defaultActions[act.action]; existed {
		for _, cb := range actions {
			act.runCallback(cb, nil)
		}
	}
	for _, after := range act.after {
		act.runCallback(after, nil)
	}
	act.game.SystemActions.Run(act.action, act.objs)
}
