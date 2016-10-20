package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

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
type PropertyEntry struct {
	*M.PropertyModel
	lower string
}
type PropertyList []*PropertyEntry
type PropertyCache map[ident.Id]PropertyList

// single/plural properties need some fixing
const pluralString = "plural"
const singularString = "singular"

var pluralId = ident.MakeId(pluralString)
var singularId = ident.MakeId(singularString)

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
	actions := mdl.actionList()
	return len(actions)
}

func (mdl *Metal) ActionNum(i int) meta.Action {
	actions := mdl.actionList()
	// panics on out of range
	return actionInfo{mdl, actions[i]}
}

func (mdl *Metal) GetAction(id ident.Id) (ret meta.Action, okay bool) {
	if a, ok := mdl.Actions[id]; ok {
		ret, okay = actionInfo{mdl, a}, true
	}
	return
}

func (mdl *Metal) NumEvent() int {
	events := mdl.eventList()
	return len(events)
}

func (mdl *Metal) EventNum(i int) meta.Event {
	// panics on out of range
	events := mdl.eventList()
	return &eventInfo{mdl, events[i]}
}

func (mdl *Metal) GetEvent(id ident.Id) (ret meta.Event, okay bool) {
	if a, ok := mdl.Events[id]; ok {
		ret, okay = &eventInfo{mdl, a}, true
	}
	return
}

func (mdl *Metal) NumClass() int {
	classes := mdl.classList()
	return len(classes)
}

func (mdl *Metal) ClassNum(i int) meta.Class {
	classes := mdl.classList()
	return &classInfo{mdl, classes[i]}
}

func (mdl *Metal) GetClass(id ident.Id) (ret meta.Class, okay bool) {
	if c, ok := mdl.Classes[id]; ok {
		ret, okay = &classInfo{mdl, c}, true
	}
	return
}

func (mdl *Metal) NumInstance() int {
	instances := mdl.instanceList()
	return len(instances)
}

func (mdl *Metal) InstanceNum(i int) meta.Instance {
	// panics on out of range
	instances := mdl.instanceList()
	return mdl.makeInstance(instances[i])
}

func (mdl *Metal) GetInstance(id ident.Id) (ret meta.Instance, okay bool) {
	if n, ok := mdl.Instances[id]; ok {
		ret, okay = mdl.makeInstance(n), true
	}
	return
}

func (mdl *Metal) NumRelation() int {
	relations := mdl.relationList()
	return len(relations)
}

func (mdl *Metal) RelationNum(i int) meta.Relation {
	relations := mdl.relationList()
	// panics on out of range
	r := relations[i]
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

func (mdl *Metal) NumParserAction() int {
	return len(mdl.ParserActions)
}

func (mdl *Metal) Pluralize(single string) (plural string) {
	if res, ok := mdl.SingleToPlural[single]; ok {
		plural = res
	} else {
		plural = lang.Pluralize(single)
	}
	return
}

func (mdl *Metal) AreCompatible(child, parent ident.Id) (okay bool) {
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
func (mdl *Metal) MatchNounName(n string, f func(ident.Id) bool) (int, bool) {
	return mdl.Aliases.Try(n, f)
}

func (mdl *Metal) makeInstance(n *M.InstanceModel) meta.Instance {
	return &instInfo{mdl, n}
}

// FIX: these lists are NG.
// *) id rather not have dynamic memory creation / pinning
// *) it cant respond to new instances
// *) itd be better to have real queries
// NOTE: we use the list caches for len() so gopherjs can avoid Object.keys
func (mdl *Metal) instanceList() []*M.InstanceModel {
	if mdl._instances == nil {
		// commented out re: gopherjs
		// mdl._instances = make([]*M.InstanceModel, 0, len(mdl.Instances))
		for _, v := range mdl.Instances {
			mdl._instances = append(mdl._instances, v)
		}
	}
	return mdl._instances
}
func (mdl *Metal) eventList() []*M.EventModel {
	if mdl._events == nil {
		// commented out re: gopherjs
		// mdl._events = make([]*M.EventModel, 0, len(mdl.Events))
		for _, v := range mdl.Events {
			mdl._events = append(mdl._events, v)
		}
	}
	return mdl._events
}
func (mdl *Metal) classList() []*M.ClassModel {
	if mdl._classes == nil {
		// commented out re: gopherjs
		// mdl._classes = make([]*M.ClassModel, 0, len(mdl.Classes))
		for _, v := range mdl.Classes {
			mdl._classes = append(mdl._classes, v)
		}
	}
	return mdl._classes
}
func (mdl *Metal) relationList() []*M.RelationModel {
	if mdl._relations == nil {
		// commented out re: gopherjs
		// mdl._relations = make([]*M.RelationModel, 0, len(mdl.Relations))
		for _, v := range mdl.Relations {
			mdl._relations = append(mdl._relations, v)
		}
	}
	return mdl._relations
}
func (mdl *Metal) actionList() []*M.ActionModel {
	if mdl._actions == nil {
		// commented out re: gopherjs
		// mdl._actions = make([]*M.ActionModel, 0, len(mdl.Actions))
		for _, v := range mdl.Actions {
			mdl._actions = append(mdl._actions, v)
		}
	}
	return mdl._actions
}

func (mdl *Metal) propertyList(cls *M.ClassModel) (ret PropertyList) {
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
		p := &cls.Properties[i]
		e := &PropertyEntry{p, strings.ToLower(p.Name)}
		ret = append(ret, e)
	}
	return
}

func (mdl *Metal) getZero(prop *M.PropertyModel) (ret interface{}) {
	switch prop.Type {
	case M.NumProperty:
		if !prop.IsMany {
			ret = float64(0)
		} else {
			ret = []float64{}
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
