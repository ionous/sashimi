package memory

// FIX: property change events

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
	_actions   []*M.ActionInfo
	_classes   []*M.ClassInfo
	_instances []*M.InstanceInfo
	_relations []*M.Relation
	_events    []*M.ActionInfo
	tables     table.Tables
}

func NewMemoryModel(model *M.Model, tables table.Tables) *MemoryModel {
	return &MemoryModel{Model: model, tables: tables}
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
func (mdl *MemoryModel) NumEvent() int {
	return len(mdl.Events)
}
func (mdl *MemoryModel) EventNum(i int) api.Event {
	if mdl._events == nil {
		mdl._events = make([]*M.ActionInfo, 0, len(mdl.Events))
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
	return &classInfo{mdl, c, nil}
}

func (mdl *MemoryModel) GetClass(id ident.Id) (ret api.Class, okay bool) {
	if c, ok := mdl.Classes[id]; ok {
		ret, okay = &classInfo{mdl, c, nil}, true
	}
	return
}

func (c *classInfo) getProperty(p M.IProperty) api.Property {
	cls := c.ClassInfo
	return &propInfo{c.mdl, c.Id, p, func() interface{} {
		return p.GetZero(cls.Constraints)
	}}
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

func (mdl *MemoryModel) ParserActionNum(i int) api.ParserAction {
	// panics on out of range
	p := mdl.ParserActions[i]
	return api.ParserAction{actionInfo{mdl, p.Action}, p.Commands}
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

func (mdl *MemoryModel) makeInstance(n *M.InstanceInfo) api.Instance {
	cls := &classInfo{mdl, n.Class, nil}
	return instInfo{mdl, n, cls}
}
func (mdl MemoryModel) getObjects(src, rel ident.Id, isRev bool) []ident.Id {
	table := mdl.getTable(rel)
	return table.List(src, isRev)
}

func (mdl MemoryModel) getTable(rel ident.Id) (ret *table.Table) {
	if table, ok := mdl.tables[rel]; !ok {
		panic(fmt.Sprintf("internal error, no table found for relation %s", rel))
	} else {
		ret = table
	}
	return
}

type actionInfo struct {
	mdl *MemoryModel
	*M.ActionInfo
}

func (a actionInfo) GetId() ident.Id {
	return a.Id
}

func (a actionInfo) GetActionName() string {
	return a.ActionName
}

func (a actionInfo) GetEvent() (ret api.Event) {
	id := M.MakeStringId(a.EventName)
	if e, ok := a.mdl.Events[id]; !ok {
		panic(fmt.Sprintf("internal error, no event found for action %s", a.ActionName))
	} else {
		ret = eventInfo{a.mdl, e}
	}
	return
}

func (a actionInfo) GetNouns() api.Nouns {
	ret := make(api.Nouns, len(a.NounTypes))
	for i, c := range a.NounTypes {
		ret[i] = c.Id
	}
	return ret
}

type eventInfo struct {
	mdl *MemoryModel
	*M.ActionInfo
}

func (e eventInfo) GetId() ident.Id {
	return e.Id
}

func (e eventInfo) GetEventName() string {
	return e.EventName
}

type classInfo struct {
	mdl *MemoryModel
	*M.ClassInfo
	_properties []M.IProperty
}

func (c classInfo) GetId() ident.Id {
	return c.Id
}

func (c classInfo) GetParentClass() (ret api.Class) {
	if p := c.Parent; p != nil {
		ret = &classInfo{c.mdl, p, nil}
	}
	return
}

func (c classInfo) GetOriginalName() string {
	return c.Plural
}

func (c *classInfo) NumProperty() int {
	props := c.props()
	return len(props)
}

func (c *classInfo) PropertyNum(i int) api.Property {
	p := c.propertyNum(i)
	return c.getProperty(p)
}

func (c *classInfo) propertyNum(i int) M.IProperty {
	props := c.props()
	// panics on out of range
	return props[i]
}

func (c *classInfo) GetProperty(id ident.Id) (ret api.Property, okay bool) {
	if p, ok := c.PropertyById(id); ok {
		ret, okay = c.getProperty(p), true
	}
	return
}

func (c *classInfo) GetPropertyByChoice(id ident.Id) (ret api.Property, okay bool) {
	if p, _, ok := c.PropertyByChoiceId(id); ok {
		ret, okay = c.getProperty(p), true
	}
	return
}

func (c *classInfo) props() []M.IProperty {
	if c._properties == nil {
		props := c.AllProperties()
		p := make([]M.IProperty, 0, len(props))
		for _, v := range props {
			p = append(p, v)
		}
		c._properties = p
	}
	return c._properties
}

type instInfo struct {
	mdl *MemoryModel
	*M.InstanceInfo
	class *classInfo
}

func (n instInfo) GetId() ident.Id {
	return n.Id
}
func (n instInfo) GetParentClass() api.Class {
	return &classInfo{n.mdl, n.Class, nil}
}

func (n instInfo) GetOriginalName() string {
	return n.Name
}

func (n instInfo) NumProperty() int {
	return n.class.NumProperty()
}

func (n instInfo) PropertyNum(i int) (ret api.Property) {
	p := n.class.propertyNum(i)
	return n.getProperty(p)
}

func (n instInfo) GetProperty(id ident.Id) (ret api.Property, okay bool) {
	if p, ok := n.Class.PropertyById(id); ok {
		ret, okay = n.getProperty(p), true
	}
	return
}

func (n instInfo) GetPropertyByChoice(id ident.Id) (ret api.Property, okay bool) {
	if p, _, ok := n.Class.PropertyByChoiceId(id); ok {
		ret, okay = n.getProperty(p), true
	}
	return
}

func (n instInfo) getProperty(p M.IProperty) api.Property {
	inst := n.InstanceInfo
	return &propInfo{n.mdl, n.Id, p, func() (ret interface{}) {
		if v, ok := inst.Values[p.GetId()]; ok {
			ret = v
		} else {
			ret = p.GetZero(inst.Class.Constraints)
		}
		return
	}}
}

type propInfo struct {
	mdl  *MemoryModel
	src  ident.Id
	prop M.IProperty
	val  func() interface{} // hooray virtual funcions.
}

func (p propInfo) String() string {
	return fmt.Sprintf("%s.%s", p.src, p.prop.GetId())
}

func (p propInfo) GetId() ident.Id {
	return p.prop.GetId()
}

func (p propInfo) GetType() api.PropertyType {
	switch p := p.prop.(type) {
	case M.NumProperty:
		return api.NumProperty
	case M.TextProperty:
		return api.TextProperty
	case M.EnumProperty:
		return api.StateProperty
	case M.PointerProperty:
		return api.ObjectProperty
	case M.RelativeProperty:
		if p.IsMany {
			return api.ObjectProperty | api.ArrayProperty
		} else {
			return api.ObjectProperty
		}
	default:
		panic("unknown property type")
	}
}

func (p propInfo) GetValue() api.Value {
	switch m := p.prop.(type) {
	case M.NumProperty:
		return numValue{panicValue(p)}
	case M.TextProperty:
		return textValue{panicValue(p)}
	case M.EnumProperty:
		return enumValue{panicValue(p)}
	case M.PointerProperty:
		return pointerValue{panicValue(p)}
	case M.RelativeProperty:
		if !m.IsMany {
			return singleValue(p)
		}
	default:
		panic("unknown property type")
	}
	panic("invalid property type")
}

func (p propInfo) GetValues() api.Values {
	switch m := p.prop.(type) {
	case M.NumProperty:
	case M.TextProperty:
	case M.EnumProperty:
	case M.PointerProperty:
	case M.RelativeProperty:
		if m.IsMany {
			return manyValue(p)
		}
	default:
		panic("unknown property type")
	}
	panic("invalid property type")
}

type panicValue propInfo

func (p panicValue) GetNum() float32 {
	panic(fmt.Errorf("get num not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetNum(float32) {
	panic(fmt.Errorf("set num not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetText() string {
	panic(fmt.Errorf("get text not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetText(string) {
	panic(fmt.Errorf("set text not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetState() ident.Id {
	panic(fmt.Errorf("get state not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetState(ident.Id) {
	panic(fmt.Errorf("set state not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetObject() ident.Id {
	panic(fmt.Errorf("get object not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetObject(ident.Id) error {
	panic(fmt.Errorf("set object not supported for property %v", p.prop.GetId()))
}

type numValue struct{ panicValue }

func (p numValue) GetNum() float32 {
	return p.val().(float32)
}

type textValue struct{ panicValue }

func (p textValue) GetText() string {
	return p.val().(string)
}

type enumValue struct{ panicValue }

func (p enumValue) GetState() (ret ident.Id) {
	if idx, ok := p.val().(int); !ok {
		panic(fmt.Sprintf("%v %T", p.val(), p.val()))
	} else {
		enum := p.prop.(M.EnumProperty)
		c, _ := enum.IndexToChoice(idx)
		ret = c
	}
	return ret
}

type pointerValue struct{ panicValue }

func (p pointerValue) GetObject() ident.Id {
	return p.val().(ident.Id)
}

func (p pointerValue) SetObject(ident.Id) error {
	// FIX: ignoring, instance sets are doubled up to make relative roperties work
	return nil
}

type objectReadValue struct {
	panicValue
	currentVal ident.Id
}

func (p objectReadValue) GetObject() (ret ident.Id) {
	return p.currentVal
}

type objectWriteValue struct {
	objectReadValue
}

// the one side of a many-to-one, one-to-one, or one-to-many relation.
func singleValue(p propInfo) api.Value {
	rel := p.prop.(M.RelativeProperty)
	objs := p.mdl.getObjects(p.src, rel.Relation, rel.IsRev)
	var v ident.Id
	if len(objs) > 0 {
		v = objs[0]
	}
	return objectWriteValue{objectReadValue{panicValue(p), v}}
}

func (p objectWriteValue) SetObject(id ident.Id) (err error) {
	p.mdl.clearValues(p.src, p.prop.(M.RelativeProperty))
	if e := p.mdl.appendObject(id, p.src, p.prop.(M.RelativeProperty)); e != nil {
		err = e
	} else {
		p.currentVal = id
	}
	return err
}

type objectList struct {
	panicValue
	objs []ident.Id
}

func manyValue(p propInfo) api.Values {
	rel := p.prop.(M.RelativeProperty)
	objs := p.mdl.getObjects(p.src, rel.Relation, rel.IsRev)
	return objectList{panicValue(p), objs}
}

func (p objectList) NumValue() int {
	return len(p.objs)
}

func (p objectList) ValueNum(i int) api.Value {
	return objectReadValue{p.panicValue, p.objs[i]}
}

// the triggers for the property watchers will have to be at this level since the concept of relation is now at this level.... only needs to happend for relation types
//	prev, next := oa.game.Objects[prev], oa.game.Objects[next]
// oa.game.Properties.VisitWatchers(func(ch PropertyChange) {
// 	ch.ReferenceChange(oa.gobj, p.GetId(), rel.Relates, prev, next)
// })
//

func (p objectList) ClearValues() {
	p.mdl.clearValues(p.src, p.prop.(M.RelativeProperty))
	p.objs = nil
}

func (objectList) AppendNum(float32) {
	panic("not implemented")
}
func (p objectList) AppendText(string) {
	panic("not implemented")
}

func (p objectList) AppendObject(id ident.Id) (err error) {
	if e := p.mdl.appendObject(id, p.src, p.prop.(M.RelativeProperty)); e != nil {
		err = e
	} else {
		p.objs = append(p.objs, id)
	}
	return
}

func (mdl MemoryModel) clearValues(src ident.Id, rel M.RelativeProperty) {
	table := mdl.getTable(rel.Relation)
	isRev := rel.IsRev
	table.Remove(func(x, y ident.Id) bool {
		var test, other ident.Id
		if isRev {
			test, other = y, x
		} else {
			test, other = x, y
		}
		_ = other
		return src == test
	})
}

func (mdl MemoryModel) appendObject(dst, src ident.Id, rel M.RelativeProperty) (err error) {
	table := mdl.getTable(rel.Relation)
	if other, ok := mdl.Instances[dst]; !ok {
		err = fmt.Errorf("no such instance %s", dst)
	} else if !mdl.AreCompatible(other.Class.Id, rel.Relates) {
		err = fmt.Errorf("%s not compatible with %v in relation %v", other, rel.Relates, rel.Relation)
	} else {
		if rel.IsRev {
			dst, src = src, dst
		}
		table.Add(src, dst)
	}
	return err
}

// func (r *relInfo) RemoveRelative(src, dst ident.Id) {
// 	table := p.mdl.getTable(rel.Relation)
// 	isRev := r.getRelative(src).IsRev
// 	table.Remove(func(x, y ident.Id) bool {
// 		return (!isRev && dst == x) || (isRev && src == x)
// 	})
// }

// // FIX: im not a huge fan of the property search, its only needed for the auto-inversion:
// // FIX: im not a huge fan of the auto-inversion, can this be solved from the client side?
// func (r *relInfo) getRelative(src ident.Id) (srcProp *M.RelativeProperty) {
// 	if src == r._relative.src {
// 		srcProp = r._relative.prop
// 	} else {
// 		inst := r.mdl.Instances[src]
// 		inst.Class.Visit(func(cls *M.ClassInfo) (finished bool) {
// 			for _, p := range cls.Properties {
// 				if rel, ok := p.(M.RelativeProperty); ok {
// 					if rel.Relation == r.Id {
// 						srcProp = &rel
// 						finished = true
// 						break
// 					}
// 				}
// 			}
// 			return
// 		})
// 		r._relative.src, r._relative.prop = src, srcProp
// 	}
// 	// outside the cache test to handle empty src id
// 	if srcProp == nil {
// 		panic(fmt.Sprintf("'%v' is not related by '%v'", src, r.Relation))
// 	}
// 	return
// }
