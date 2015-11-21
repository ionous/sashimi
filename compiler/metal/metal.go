package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
)

type GenericValue interface{}

type Metal struct {
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

func NewMetal(m *M.Model, v ObjectValue) *Metal {
	return &Metal{
		Model:        m,
		objectValues: v,
		_properties:  make(PropertyCache),
	}
}

func (mdl *Metal) NumAction() int {
	return len(mdl.Actions)
}

func (mdl *Metal) ActionNum(i int) meta.Action {
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

func (mdl *Metal) GetAction(id ident.Id) (ret meta.Action, okay bool) {
	if a, ok := mdl.Actions[id]; ok {
		ret, okay = actionInfo{mdl, a}, true
	}
	return
}

func (mdl *Metal) NumEvent() int {
	return len(mdl.Events)
}

func (mdl *Metal) EventNum(i int) meta.Event {
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

func (mdl *Metal) GetEvent(id ident.Id) (ret meta.Event, okay bool) {
	if a, ok := mdl.Events[id]; ok {
		ret, okay = eventInfo{mdl, a}, true
	}
	return
}

func (mdl Metal) NumClass() int {
	return len(mdl.Classes)
}

func (mdl *Metal) ClassNum(i int) meta.Class {
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

func (mdl *Metal) GetClass(id ident.Id) (ret meta.Class, okay bool) {
	if c, ok := mdl.Classes[id]; ok {
		ret, okay = classInfo{mdl, c}, true
	}
	return
}

func (mdl Metal) NumInstance() int {
	return len(mdl.Instances)
}

func (mdl *Metal) InstanceNum(i int) meta.Instance {
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

func (mdl *Metal) GetInstance(id ident.Id) (ret meta.Instance, okay bool) {
	if n, ok := mdl.Instances[id]; ok {
		ret, okay = mdl.makeInstance(n), true
	}
	return
}

func (mdl Metal) NumRelation() int {
	return len(mdl.Relations)
}

func (mdl *Metal) RelationNum(i int) meta.Relation {
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

func (mdl *Metal) GetRelation(id ident.Id) (ret meta.Relation, okay bool) {
	if r, ok := mdl.Relations[id]; ok {
		ret, okay = relInfo{mdl, r}, true
	}
	return
}

func (mdl *Metal) ParserActionNum(i int) meta.ParserAction {
	// panics on out of range
	p := mdl.ParserActions[i]
	return meta.ParserAction{p.Action, p.Commands}
}

func (mdl Metal) NumParserAction() int {
	return len(mdl.ParserActions)
}

func (mdl Metal) Pluralize(single string) (plural string) {
	if res, ok := mdl.SingleToPlural[single]; ok {
		plural = res
	} else {
		plural = lang.Pluralize(single)
	}
	return
}

func (mdl Metal) AreCompatible(child, parent ident.Id) (okay bool) {
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
func (mdl Metal) MatchNounName(n string, f func(ident.Id) bool) (int, bool) {
	return mdl.Aliases.Try(n, f)
}

func (mdl *Metal) makeInstance(n *M.InstanceModel) meta.Instance {
	return instInfo{mdl, n}
}

func (mdl *Metal) getPropertyList(cls *M.ClassModel) (ret PropertyList) {
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

func (mdl *Metal) makePropertyList(cls *M.ClassModel) (ret PropertyList) {
	for i, _ := range cls.Properties {
		ret = append(ret, &cls.Properties[i])
	}
	return
}

func (mdl Metal) getZero(prop *M.PropertyModel) (ret interface{}) {
	switch prop.Type {
	case M.NumProperty:
		if !prop.IsMany {
			ret = float32(0)
		} else {
			ret = []float32{}
		}
	case M.TextProperty:
		if !prop.IsMany {
			ret = ""
		} else {
			ret = []string{}
		}
	case M.EnumProperty:
		enum := mdl.Enumerations[prop.Id]
		ret = enum.Best()
	case M.PointerProperty:
		if !prop.IsMany {
			ret = ident.Empty()
		} else {
			ret = []ident.Id{}
		}
	default:
		panic(fmt.Errorf("GetZero not supported for property %s type %v", prop.Id, prop.Type))
	}
	return ret
}
