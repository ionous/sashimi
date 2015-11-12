package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
)

type MemoryModel struct {
	*M.Model
	_actions     []*M.ActionInfo
	_events      []*M.EventInfo
	_classes     []*M.ClassInfo
	_instances   []*M.InstanceInfo
	_relations   []*M.Relation
	_properties  map[ident.Id]PropertyList
	objectValues ObjectValue
	tables       table.Tables
	actions      DefaultActions
	//capture, bubble EventCallbacks
}

// actionId -> callbackId
type DefaultActions map[ident.Id][]ident.Id

// indexed by eventId
//type EventCallbacks map[ident.Id][]M.ListenerCallback

// array of properties ( for flat cache )
type PropertyList []M.IProperty

func NewMemoryModel(m *M.Model, v ObjectValue, t table.Tables) *MemoryModel {
	actions := make(DefaultActions)

	for _, handler := range m.ActionHandlers {
		act, callback, useCapture := handler.Action, handler.Callback, handler.UseCapture()
		arr := actions[act]
		// FIX: for now treating target as bubble,
		// really the compiler should hand off a sorted flat list based on three separate groups
		// target growing in the same direction as after, but distinctly in the middle of things.
		if !useCapture {
			arr = append(arr, callback)
		} else {
			// prepend:
			arr = append([]ident.Id{callback}, arr...)
		}
		actions[act] = arr
	}

	// capture, bubble := make(EventCallbacks), make(EventCallbacks)
	// for _, l := range m.EventListeners {
	// 	e, cb := l.Event, l.ListenerCallback
	// 	var callbacks EventCallbacks
	// 	if cb.UseCapture() {
	// 		callbacks = capture
	// 	} else {
	// 		callbacks = bubble
	// 	}
	// 	// append
	// 	var arr = callbacks[e]
	// 	arr = append(arr, cb)
	// 	callbacks[e] = arr
	// }

	flatCache := make(map[ident.Id]PropertyList)
	return &MemoryModel{
		Model:        m,
		objectValues: v,
		tables:       t,
		_properties:  flatCache,
		actions:      actions,
		// capture:      capture,
		// bubble:       bubble,
	}
}

func (mdl *MemoryModel) NumAction() int {
	return len(mdl.Actions)
}
func (mdl *MemoryModel) ActionNum(i int) api.Action {
	if mdl._actions == nil {
		mdl._actions = make([]*M.ActionInfo, 0, len(mdl.Actions))
		for _, v := range mdl.Actions {
			mdl._actions = append(mdl._actions, v)
		}
	}
	// panics on out of range
	a := mdl._actions[i]
	return actionInfo{mdl, a}
}

func (mdl *MemoryModel) GetAction(id ident.Id) (ret api.Action, okay bool) {
	if a, ok := mdl.Actions[id]; ok {
		ret, okay = actionInfo{mdl, a}, true
	}
	return
}
func (mdl *MemoryModel) GetDefaultCallbacks(id ident.Id) (ret api.ActionCallbacks, okay bool) {
	if a, ok := mdl.GetAction(id); ok {
		ret, okay = a.(actionInfo).GetCallbacks()
	}
	return
}

func (mdl *MemoryModel) NumEvent() int {
	return len(mdl.Events)
}

func (mdl *MemoryModel) EventNum(i int) api.Event {
	if mdl._events == nil {
		mdl._events = make([]*M.EventInfo, 0, len(mdl.Events))
		for _, v := range mdl.Events {
			mdl._events = append(mdl._events, v)
		}
	}
	// panics on out of range
	a := mdl._events[i]
	return eventInfo{mdl, a}
}

func (mdl *MemoryModel) GetEvent(id ident.Id) (ret api.Event, okay bool) {
	if a, ok := mdl.Events[id]; ok {
		ret, okay = eventInfo{mdl, a}, true
	}
	return
}

func (mdl MemoryModel) NumClass() int {
	return len(mdl.Classes)
}
func (mdl *MemoryModel) ClassNum(i int) api.Class {
	if mdl._classes == nil {
		mdl._classes = make([]*M.ClassInfo, 0, len(mdl.Classes))
		for _, v := range mdl.Classes {
			mdl._classes = append(mdl._classes, v)
		}
	}
	// panics on out of range
	c := mdl._classes[i]
	return classInfo{mdl, c}
}

func (mdl *MemoryModel) GetClass(id ident.Id) (ret api.Class, okay bool) {
	if c, ok := mdl.Classes[id]; ok {
		ret, okay = classInfo{mdl, c}, true
	}
	return
}

func (mdl MemoryModel) NumInstance() int {
	return len(mdl.Instances)
}

func (mdl *MemoryModel) InstanceNum(i int) api.Instance {
	if mdl._instances == nil {
		mdl._instances = make([]*M.InstanceInfo, 0, len(mdl.Instances))
		for _, v := range mdl.Instances {
			mdl._instances = append(mdl._instances, v)
		}
	}
	// panics on out of range
	n := mdl._instances[i]
	return mdl.makeInstance(n)
}

func (mdl *MemoryModel) GetInstance(id ident.Id) (ret api.Instance, okay bool) {
	if n, ok := mdl.Instances[id]; ok {
		ret, okay = mdl.makeInstance(n), true
	}
	return
}

func (mdl MemoryModel) NumRelation() int {
	return len(mdl.Relations)
}

func (mdl *MemoryModel) RelationNum(i int) api.Relation {
	if mdl._relations == nil {
		mdl._relations = make([]*M.Relation, 0, len(mdl.Relations))
		for _, v := range mdl.Relations {
			mdl._relations = append(mdl._relations, v)
		}
	}
	// panics on out of range
	r := mdl._relations[i]
	return relInfo{mdl, r}
}

func (mdl *MemoryModel) GetRelation(id ident.Id) (ret api.Relation, okay bool) {
	if r, ok := mdl.Relations[id]; ok {
		ret, okay = relInfo{mdl, r}, true
	}
	return
}

func (mdl *MemoryModel) ParserActionNum(i int) api.ParserAction {
	// panics on out of range
	p := mdl.ParserActions[i]
	return api.ParserAction{p.Action, p.Commands}
}

func (mdl MemoryModel) NumParserAction() int {
	return len(mdl.ParserActions)
}

func (mdl MemoryModel) Pluralize(single string) (plural string) {
	if res, ok := mdl.SingleToPlural[single]; ok {
		plural = res
	} else {
		plural = lang.Pluralize(single)
	}
	return
}

func (mdl MemoryModel) AreCompatible(child, parent ident.Id) (okay bool) {
	if c, ok := mdl.Classes[child]; ok {
		okay = c.CompatibleWith(parent)
	}
	return
}

// hrmmm...
func (mdl MemoryModel) MatchNounName(n string, f func(ident.Id) bool) (int, bool) {
	return mdl.NounNames.Try(n, f)
}

func (mdl *MemoryModel) makeInstance(n *M.InstanceInfo) api.Instance {
	cls := classInfo{mdl, n.Class}
	return instInfo{mdl, n, cls}
}

func (mdl MemoryModel) getObjects(src, rel ident.Id, isRev bool) []ident.Id {
	table := mdl.getTable(rel)
	return table.List(src, isRev)
}

func (mdl *MemoryModel) getPropertyList(cls *M.ClassInfo) (ret PropertyList) {
	if props, ok := mdl._properties[cls.Id]; ok {
		ret = props
	} else {
		props := cls.AllProperties()
		ret = make([]M.IProperty, 0, len(props))
		for _, v := range props {
			ret = append(ret, v)
		}
		mdl._properties[cls.Id] = ret
	}
	return
}

func (mdl MemoryModel) getTable(rel ident.Id) (ret *table.Table) {
	if table, ok := mdl.tables[rel]; !ok {
		panic(fmt.Sprintf("internal error, no table found for relation %s", rel))
	} else {
		ret = table
	}
	return
}

// returns error if not compatible.
func (mdl MemoryModel) canAppend(dst, src ident.Id, rel M.RelativeProperty) (err error) {
	if other, ok := mdl.Instances[dst]; !ok {
		err = fmt.Errorf("no such instance '%s'", dst)
	} else if !mdl.AreCompatible(other.Class.Id, rel.Relates) {
		err = fmt.Errorf("%s not compatible with %v in relation %v", other, rel.Relates, rel.Relation)
	}
	return err
}

func (mdl MemoryModel) appendObject(dst, src ident.Id, rel M.RelativeProperty) {
	if rel.IsRev {
		dst, src = src, dst
	}
	table := mdl.getTable(rel.Relation)
	table.Add(src, dst)
}
