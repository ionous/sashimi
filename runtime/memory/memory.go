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
	// objects ordered by index for linear travseral
	_actions   []*M.ActionModel
	_events    []*M.EventModel
	_classes   []*M.ClassModel
	_instances []*M.InstanceModel
	_relations []*M.RelationModel
	// flat cache of properties for each class
	_properties  PropertyCache
	objectValues ObjectValue
}

// actionId -> callbackId
type DefaultActions map[ident.Id][]ident.Id

// indexed by eventId
type EventCallbacks map[ident.Id][]M.ListenerModel

// array of properties ( for flat cache ) ordered for linear traversal.
// children appear first; redudent properties are *not* removed.
type PropertyRef *M.PropertyModel
type PropertyList []PropertyRef
type PropertyCache map[ident.Id]PropertyList

// single/plural properties need some fixing
var singular = ident.MakeId("singular")
var plural = ident.MakeId("plural")

func merge(base, more PropertyList) (ret PropertyList) {
	ret = base
	for i, _ := range more {
		prop := more[i]
		ret = append(ret, prop)
	}
	return ret
}

func NewMemoryModel(m *M.Model, v ObjectValue, t table.Tables) *MemoryModel {
	return &MemoryModel{
		Model:        m,
		objectValues: v,
		_properties:  make(PropertyCache),
	}
}

func (mdl *MemoryModel) NumAction() int {
	return len(mdl.Actions)
}

func (mdl *MemoryModel) ActionNum(i int) api.Action {
	if mdl._actions == nil {
		mdl._actions = make([]*M.ActionModel, 0, len(mdl.Actions))
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

func (mdl *MemoryModel) NumEvent() int {
	return len(mdl.Events)
}

func (mdl *MemoryModel) EventNum(i int) api.Event {
	if mdl._events == nil {
		mdl._events = make([]*M.EventModel, 0, len(mdl.Events))
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
		mdl._classes = make([]*M.ClassModel, 0, len(mdl.Classes))
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
		mdl._instances = make([]*M.InstanceModel, 0, len(mdl.Instances))
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
		mdl._relations = make([]*M.RelationModel, 0, len(mdl.Relations))
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
		if c.Id == parent {
			okay = true
		} else {
			for _, pid := range c.Parents {
				if pid == parent {
					okay = true
					break
				}
			}
		}
	}
	return
}

// hrmmm...
func (mdl MemoryModel) MatchNounName(n string, f func(ident.Id) bool) (int, bool) {
	return mdl.Aliases.Try(n, f)
}

func (mdl *MemoryModel) makeInstance(n *M.InstanceModel) api.Instance {
	return instInfo{mdl, n}
}

func (mdl MemoryModel) getObjects(src, rel ident.Id, isRev bool) []ident.Id {
	table := mdl.getTable(rel)
	return table.List(src, isRev)
}

func (mdl *MemoryModel) getPropertyList(cls *M.ClassModel) (ret PropertyList) {
	if props, ok := mdl._properties[cls.Id]; ok {
		ret = props
	} else {
		ret = merge(ret, mdl.makePropertyList(cls))
		for _, pid := range cls.Parents {
			parent := mdl.Classes[pid]
			ret = merge(ret, mdl.makePropertyList(parent))
		}
		mdl._properties[cls.Id] = ret
	}
	return
}

func (mdl *MemoryModel) makePropertyList(cls *M.ClassModel) (ret PropertyList) {
	for i, _ := range cls.Properties {
		ret = append(ret, &cls.Properties[i])
	}
	return
}

func (mdl MemoryModel) getTable(rel ident.Id) (ret *table.Table) {
	if table, ok := mdl.Tables[rel]; !ok {
		panic(fmt.Sprintf("internal error, no table found for relation %s", rel))
	} else {
		ret = table
	}
	return
}

func (mdl MemoryModel) getZero(prop *M.PropertyModel) (ret interface{}) {
	switch prop.Type {
	case M.NumProperty:
		if !prop.IsMany {
			ret = float32(0)
		} else {
			ret = []interface{}{}
		}
	case M.TextProperty:
		if !prop.IsMany {
			ret = ""
		} else {
			ret = []interface{}{}
		}
	case M.EnumProperty:
		//enum := mdl.Enumerations[prop.Id]
		ret = 1 //enum.ChoiceToIndex(enum.Best())
	case M.PointerProperty:
		if !prop.IsMany {
			ret = ident.Empty()
		} else {
			ret = []interface{}{}
		}

	default:
		panic(fmt.Errorf("GetZero not supported for property %s type %v", prop.Id, prop.Type))
	}
	return ret
}

// returns error if not compatible.
func (mdl MemoryModel) canAppend(dst, src ident.Id, rel *M.PropertyModel) (err error) {
	if other, ok := mdl.Instances[dst]; !ok {
		err = fmt.Errorf("no such instance '%s'", dst)
	} else if !mdl.AreCompatible(other.Class, rel.Relates) {
		err = fmt.Errorf("%s not compatible with %v in relation %v", other, rel.Relates, rel.Relation)
	}
	return err
}

func (mdl MemoryModel) appendObject(dst, src ident.Id, rel *M.PropertyModel) {
	if rel.IsRev {
		dst, src = src, dst
	}
	table := mdl.getTable(rel.Relation)
	table.Add(src, dst)
}
