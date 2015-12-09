package internal

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

var _ = fmt.Println

// RuntimeAction contains data for event handlers and actions.
type RuntimeAction struct {
	game      *Game // FIX: YUCK!
	action    meta.Action
	objs      []meta.Instance
	after     []QueuedCallback // FIX: WHY DO WE COPY THIS!?!
	cancelled bool
}

func NewRuntimeAction(g *Game, act meta.Action, objects []meta.Instance) *RuntimeAction {
	return &RuntimeAction{game: g, action: act, objs: objects}
}

func (act *RuntimeAction) GetTarget() meta.Instance {
	return act.objs[0]
}

// each action can have a chain of default actions
type QueuedCallback struct {
	src  ident.Id
	call G.Callback
}

// FIX: change callbacks to include a source file/line location
func (q QueuedCallback) String() string {
	return fmt.Sprint(q.call)
}

func (act *RuntimeAction) Cancelled() bool {
	return act.cancelled
}

// queue for running after the default actions
func (act *RuntimeAction) runAfterDefaults(cb QueuedCallback) {
	act.after = append(act.after, cb)
}

// Default actions occur after event processing assuming that they have not been cancelled.
func (act *RuntimeAction) runDefaultActions() (err error) {
	if callbacks, ok := act.action.GetCallbacks(); ok {
		for i := 0; i < callbacks.NumCallback(); i++ {
			cb := callbacks.CallbackNum(i)
			play := act.game.newPlay(act, ident.Empty())

			if found, ok := act.game.LookupCallback(cb); !ok {
				err = fmt.Errorf("internal error, couldnt find callback %s", cb)
				panic(err.Error())
				break
			} else {
				found(play)
			}
		}

		for _, after := range act.after {
			play := act.game.newPlay(act, ident.Empty())
			after.call(play)
		}
	}
	return
}

// findByName:
func (act *RuntimeAction) findByName(name string, hint ident.Id) (ret meta.Instance, okay bool) {
	if obj, ok := act.findByParamName(name); ok {
		okay, ret = true, obj
	} else if obj, ok := act.findByClassName(name, hint); ok {
		okay, ret = true, obj
	}
	return
}

// findByParamName: source, target, or context
func (act *RuntimeAction) findByParamName(name string) (ret meta.Instance, okay bool) {
	for index, src := range []string{"action.Source", "action.Target", "action.Context"} {
		if strings.EqualFold(name, src) {
			ret, okay = act.getObject(index)
			break
		}
	}
	return
}

// findByClassName:
func (act *RuntimeAction) findByClassName(name string, hint ident.Id) (ret meta.Instance, okay bool) {
	clsid := MakeStringId(act.game.Model.Pluralize(lang.StripArticle(name)))
	if clsid == hint {
		ret, okay = act.getObject(0)
	} else {
		ret, okay = act.findByClass(clsid)
	}
	return
}

// findByExactClass; true if found
func (act *RuntimeAction) findByClass(id ident.Id) (ret meta.Instance, okay bool) {
	// these are the classes originally named in the action declaration; not the sub-classes of the event target. ie. s.The("actors", Can("crawl"), not s.The("babies", When("crawling")
	if obj, ok := act.findByExactClass(id); ok {
		ret, okay = obj, true
	} else {
		// when all else fails try compatible classes one by one.
		ret, okay = act.findBySimilarClass(id)
	}
	return ret, okay
}

// findByExactClass; true if found
func (act *RuntimeAction) findByExactClass(id ident.Id) (ret meta.Instance, okay bool) {
	for i, nounClass := range act.action.GetNouns() {
		if same := id == nounClass; same {
			ret, okay = act.getObject(i)
			break
		}
	}
	return
}

// findBySimilarClass; true if found
func (act *RuntimeAction) findBySimilarClass(id ident.Id) (ret meta.Instance, okay bool) {
	for i, nounClass := range act.action.GetNouns() {
		if similar := act.game.Model.AreCompatible(id, nounClass); similar {
			ret, okay = act.getObject(i)
			break
		}
	}
	return
}

// getObject returns the index object; true if the index was in range.
func (act *RuntimeAction) getObject(i int) (ret meta.Instance, okay bool) {
	if i >= 0 && i < len(act.objs) {
		ret, okay = act.objs[i], true
	}
	return
}
