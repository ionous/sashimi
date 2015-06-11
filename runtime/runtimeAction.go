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
func (this *RuntimeAction) runAfterDefaults(cb G.Callback) {
	this.after = append(this.after, cb)
}

// Run the passed game callback.
// This creates the game event adapter, sets up the necessary template context, etc.
// Returns the results of the callback.
func (this *RuntimeAction) runCallback(cb G.Callback) bool {
	// FIX: it might be cooler to make this some sort of API...
	// then we could have the callback object: callback.run( IPlay, Data ) maybe...
	adapter := &GameEventAdapter{Game: this.game, data: this}
	templateValueStack.pushValues(this.values)
	defer templateValueStack.pop()
	cb(adapter)
	return !adapter.cancelled
}

//
// Default actions occur after event processing assuming that they have not been cancelled.
//
func (this *RuntimeAction) runDefaultActions() {
	// FIX: assign defaults at initialization?
	// it'd be even better if this didn't need game --
	// the main reason it does it to share code b/t Go() and the ProcessEventLoop
	if actions, existed := this.game.defaultActions[this.action]; existed {
		for _, cb := range actions {
			this.runCallback(cb)
		}
	}
	for _, after := range this.after {
		this.runCallback(after)
	}
	this.game.SystemActions.Run(this.action, this.objs)
}
