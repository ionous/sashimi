package runtime

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"strings"
)

//
// Event data for event handlers and actions.
//
type RuntimeAction struct {
	game   *Game
	action *M.ActionInfo
	objs   []*GameObject
	values map[string]TemplateValues
	after  []CallbackPair
}

// queue for running after the default actions
func (act *RuntimeAction) runAfterDefaults(cb CallbackPair) {
	act.after = append(act.after, cb)
}

// Run the passed game callback.
// This creates the game event adapter, sets up the necessary template context, etc.
// Returns the results of the callback.
func (act *RuntimeAction) runCallback(cb CallbackPair, cls *M.ClassInfo) bool {
	// FIX: it might be cooler to make act some sort of API...
	// then we could have the callback object: callback.run( IPlay, Data ) maybe...
	act.game.log.Println("calling:", act.action, cb)

	adapter := NewGameAdapter(act.game)
	adapter.data = act
	adapter.hint = cls
	templateValueStack.pushValues(act.values)
	defer templateValueStack.pop()
	cb.call(adapter)
	return !adapter.cancelled
}

//
// Default actions occur after event processing assuming that they have not been cancelled.
//
func (act *RuntimeAction) runDefaultActions() {
	//act.game.log.Println("default action:", act.action)
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

// fundByParamName: source, target, or context
func (act *RuntimeAction) findByParamName(name string) (ret G.IObject, okay bool) {
	for index, src := range []string{"action.Source", "action.Target", "action.Context"} {
		if strings.EqualFold(name, src) {
			ret, okay = act.getObject(index)
			break
		}
	}
	return ret, okay
}

// findByExactClass; true if found
func (act *RuntimeAction) findByClass(cls *M.ClassInfo) (ret G.IObject, okay bool) {
	// these are the classes originally named in the action declaration; not the sub-classes of the event target. ie. s.The("actors", Can("crawl"), not s.The("babies", When("crawling")
	if obj, ok := act.findByExactClass(cls); ok {
		ret, okay = obj, true
	} else {
		// when all else fails try compatible classes one by one.
		ret, okay = act.findBySimilarClass(cls)
	}
	return ret, okay
}

// findByExactClass; true if found
func (act *RuntimeAction) findByExactClass(cls *M.ClassInfo) (ret G.IObject, okay bool) {
	for i, nounClass := range act.action.NounTypes {
		// if DebugGet {
		//  act.game.Println("exact", i, nounClass)
		// }
		if cls == nounClass {
			ret, okay = act.getObject(i)
			break
		}
	}
	return
}

// findBySimilarClass; true if found
func (act *RuntimeAction) findBySimilarClass(cls *M.ClassInfo) (ret G.IObject, okay bool) {
	for i, nounClass := range act.action.NounTypes {
		if similar := cls.CompatibleWith(nounClass.Id); similar {
			ret, okay = act.getObject(i)
			break
		}
	}
	return
}

// getObject returns the index object; true if the index was in range.
func (act *RuntimeAction) getObject(i int) (ret G.IObject, okay bool) {
	if i >= 0 && i < len(act.objs) {
		ret, okay = NewObjectAdapter(act.game, act.objs[i]), true
	}
	return
}
