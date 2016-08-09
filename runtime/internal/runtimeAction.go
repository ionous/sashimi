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
	action    meta.Action
	objs      []meta.Instance
	after     []QueuedCallback // FIX: WHY DO WE COPY THIS!?!
	cancelled bool
}

func NewRuntimeAction(act meta.Action, objects []meta.Instance) *RuntimeAction {
	return &RuntimeAction{action: act, objs: objects}
}

func (act *RuntimeAction) GetTarget() (ret meta.Instance) {
	if len(act.objs) > 0 {
		ret = act.objs[0]
	}
	return
}

func (act *RuntimeAction) GetContext() (ret meta.Instance) {
	if len(act.objs) > 1 {
		ret = act.objs[1]
	}
	return
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

// findByName:
func (act *RuntimeAction) findByName(m meta.Model, name string, hint ident.Id) (ret meta.Instance, okay bool) {
	if obj, ok := act.findByParamName(name); ok {
		okay, ret = true, obj
	} else if obj, ok := act.findByClassName(m, name, hint); ok {
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
func (act *RuntimeAction) findByClassName(m meta.Model, name string, hint ident.Id) (ret meta.Instance, okay bool) {
	clsid := MakeStringId(m.Pluralize(lang.StripArticle(name)))
	if clsid == hint {
		ret, okay = act.getObject(0)
	} else {
		ret, okay = act.findByClass(m, clsid)
	}
	return
}

// findByExactClass; true if found
func (act *RuntimeAction) findByClass(m meta.Model, id ident.Id) (ret meta.Instance, okay bool) {
	// these are the classes originally named in the action declaration; not the sub-classes of the event target. ie. s.The("actors", Can("crawl"), not s.The("babies", When("crawling")
	if obj, ok := act.findByExactClass(m, id); ok {
		ret, okay = obj, true
	} else {
		// when all else fails try compatible classes one by one.
		ret, okay = act.findBySimilarClass(m, id)
	}
	return ret, okay
}

// findByExactClass; true if found
func (act *RuntimeAction) findByExactClass(_ meta.Model, id ident.Id) (ret meta.Instance, okay bool) {
	for i, nounClass := range act.action.GetNouns() {
		if same := id == nounClass; same {
			ret, okay = act.getObject(i)
			break
		}
	}
	return
}

// findBySimilarClass; true if found
func (act *RuntimeAction) findBySimilarClass(m meta.Model, id ident.Id) (ret meta.Instance, okay bool) {
	for i, nounClass := range act.action.GetNouns() {
		if similar := m.AreCompatible(id, nounClass); similar {
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
