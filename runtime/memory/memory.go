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
	_actions       []*M.ActionInfo
	_events        []*M.EventInfo
	_classes       []*M.ClassInfo
	_instances     []*M.InstanceInfo
	_relations     []*M.Relation
	_properties    map[ident.Id]PropertyList
	objectValues   ObjectValue
	tables         table.Tables
	defaultActions DefaultActions
}

// indexed by action id-<callbackId
type DefaultActions map[ident.Id][]ident.Id

type CallbackList struct {
	callbacks []ident.Id
}

func (cl CallbackList) NumCallback() int {
	return len(cl.callbacks)
}

func (cl CallbackList) CallbackNum(i int) ident.Id {
	p := cl.callbacks[i]
	return p // CallbackWrapper(p)
}

// type CallbackWrapper CallbackPair

// func (ac CallbackWrapper) GetAction() ident.Id {
// 	return ac.act
// }

// func (ac CallbackWrapper) GetCallback() ident.Id {
// 	return ac.callback
// }

type PropertyList []M.IProperty

// single/plural properties need some fixing
var singular = ident.MakeId("singular")
var plural = ident.MakeId("plural")

type junkProperty struct {
	id  ident.Id
	val string
}

func (p junkProperty) GetId() ident.Id {
	return p.id
}

func (p junkProperty) GetName() string {
	return p.id.String()
}

func (p junkProperty) GetZero(_ M.ConstraintSet) interface{} {
	return p.val
}

func NewMemoryModel(m *M.Model, v ObjectValue, t table.Tables) *MemoryModel {
	defaultActions := make(DefaultActions)

	// STORE/FIX: arrange action handlers by action id.
	for _, handler := range m.ActionHandlers {
		act, callback, useCapture := handler.Action, handler.Callback, handler.UseCapture()
		arr := defaultActions[act]
		// FIX: for now treating target as bubble,
		// really the compiler should hand off a sorted flat list based on three separate groups
		// target growing in the same direction as after, but distinctly in the middle of things.
		if !useCapture {
			arr = append(arr, callback)
		} else {
			// prepend:
			arr = append([]ident.Id{callback}, arr...)
		}
		defaultActions[act] = arr
	}

	return &MemoryModel{Model: m, objectValues: v, tables: t, _properties: make(map[ident.Id]PropertyList), defaultActions: defaultActions}
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
	if a, ok := mdl.defaultActions[id]; ok {
		ret, okay = CallbackList{a}, true
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

func (c *classInfo) getProperty(p M.IProperty) api.Property {
	return &propBase{
		mdl:      c.mdl,
		src:      c.Id,
		prop:     p,
		getValue: c.getValue,
		setValue: nil}
}

func (c *classInfo) getValue(p M.IProperty) GenericValue {
	return p.GetZero(c.Constraints)
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

func (mdl MemoryModel) AreCompatible(child, parent ident.Id) (okay bool) {
	if c, ok := mdl.Classes[child]; ok {
		okay = c.CompatibleWith(parent)
	}
	return
}

func (mdl MemoryModel) Pluralize(single string) (plural string) {
	if res, ok := mdl.SingleToPlural[single]; ok {
		plural = res
	} else {
		plural = lang.Pluralize(single)
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
	if e, ok := a.mdl.Events[a.EventId]; !ok {
		panic(fmt.Sprintf("internal error, no event found for action %s", a.EventId))
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
	*M.EventInfo
}

func (e eventInfo) GetId() ident.Id {
	return e.Id
}

func (e eventInfo) GetEventName() string {
	return e.EventName
}

func (e eventInfo) GetAction() (ret api.Action) {
	if a, ok := e.mdl.GetAction(e.ActionId); !ok {
		panic(fmt.Sprintf("internal error, no action found for event %s", e.ActionId))
	} else {
		ret = a
	}
	return
}

type classInfo struct {
	mdl *MemoryModel
	*M.ClassInfo
}

func (c classInfo) GetId() ident.Id {
	return c.Id
}

func (c classInfo) GetParentClass() (ret api.Class) {
	if p := c.Parent; p != nil {
		ret = classInfo{c.mdl, p}
	}
	return
}

func (c classInfo) GetOriginalName() string {
	return c.Plural
}

func (c classInfo) NumProperty() int {
	props := c.mdl.getPropertyList(c.ClassInfo)
	return len(props)
}

func (c classInfo) PropertyNum(i int) api.Property {
	p := c.propertyNum(i)
	return c.getProperty(p)
}

func (c classInfo) propertyNum(i int) M.IProperty {
	props := c.mdl.getPropertyList(c.ClassInfo)
	return props[i] // panics on out of range
}

func (c classInfo) GetProperty(id ident.Id) (ret api.Property, okay bool) {
	// hack for singular and plural properties, note: they wont show up in enumeration...
	var prop M.IProperty
	switch id {
	case plural:
		prop, okay = junkProperty{plural, c.Plural}, true
	case singular:
		prop, okay = junkProperty{singular, c.Singular}, true
	default:
		prop, okay = c.PropertyById(id)
	}
	if okay {
		ret = c.getProperty(prop)
	}
	return
}

func (c classInfo) GetPropertyByChoice(id ident.Id) (ret api.Property, okay bool) {
	if p, _, ok := c.PropertyByChoiceId(id); ok {
		ret, okay = c.getProperty(p), true
	}
	return
}

type instInfo struct {
	mdl *MemoryModel
	*M.InstanceInfo
	class classInfo
}

func (n instInfo) GetId() ident.Id {
	return n.Id
}
func (n instInfo) GetParentClass() api.Class {
	return classInfo{n.mdl, n.Class}
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
	return &propBase{
		mdl:      n.mdl,
		src:      n.Id,
		prop:     p,
		getValue: n.getValue,
		setValue: n.setValue}
}

func (n instInfo) getValue(p M.IProperty) (ret GenericValue) {
	// try the object-value interface first
	if v, ok := n.mdl.objectValues.GetValue(n.Id, p.GetId()); ok {
		ret = v
		// fall back to the instance
	} else if v, ok := n.Values[p.GetId()]; ok {
		ret = v
	} else {
		// and from there to class ( chain )
		ret = p.GetZero(n.Class.Constraints)
	}
	return
}

func (n instInfo) setValue(p M.IProperty, v GenericValue) error {
	// STORE FIX: TEST CONSTRAINTS
	return n.mdl.objectValues.SetValue(n.Id, p.GetId(), v)
}

type relInfo struct {
	mdl *MemoryModel
	*M.Relation
}

func (r relInfo) GetId() ident.Id {
	return r.Id
}

type propBase struct {
	mdl  *MemoryModel
	src  ident.Id
	prop M.IProperty
	// life's a little complicated.
	// we have a generic property base ( propBase )
	// an extension to panic on every get and set ( panicValue )
	// and overrides to implement the specific text/num/etc methods ( textValue )
	// the location of values for class and instances differs, so the class and instance pass themselves to their properties, and on to their values.
	getValue func(M.IProperty) GenericValue
	setValue func(M.IProperty, GenericValue) error
}

type propInst struct {
	propBase
}

type propClass struct {
	propBase
}

func (p propBase) String() string {
	return fmt.Sprintf("%s.%s", p.src, p.prop.GetId())
}

func (p propBase) GetId() ident.Id {
	return p.prop.GetId()
}

func (p propBase) GetType() api.PropertyType {
	err := "invalid"
	switch r := p.prop.(type) {
	case M.NumProperty:
		return api.NumProperty
	case M.TextProperty, junkProperty:
		return api.TextProperty
	case M.EnumProperty:
		return api.StateProperty
	case M.PointerProperty:
		return api.ObjectProperty
	case M.RelativeProperty:
		if r.IsMany {
			return api.ObjectProperty | api.ArrayProperty
		} else {
			return api.ObjectProperty
		}
	default:
		err = "unknown"
	}
	panic(fmt.Sprintf("GetType(%s.%s) has %s property type %T", p.src, p.prop.GetId(), err, p.prop))
}

func (p propBase) GetValue() api.Value {
	err := "invalid"
	switch m := p.prop.(type) {
	case M.NumProperty:
		return numValue{panicValue(p)}
	case M.TextProperty, junkProperty:
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
		err = "unknown"
	}
	panic(fmt.Sprintf("GetValue(%s.%s) has %s property type %T", p.src, p.prop.GetId(), err, p.prop))
}

func (p propBase) GetValues() api.Values {
	err := "invalid"
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
		err = "unknown"
	}
	panic(fmt.Sprintf("GetValues(%s.%s) has %s property type %T", p.src, p.prop.GetId(), err, p.prop))
}

func (p propBase) GetRelative() (ret api.Relative, okay bool) {
	switch prop := p.prop.(type) {
	case M.PointerProperty:
	case M.RelativeProperty:

		// get the relation
		relation := p.mdl.Relations[prop.Relation]

		// get the reverse property
		other := relation.GetOther(prop.IsRev)

		ret = api.Relative{
			Relation: prop.Relation,
			Relates:  prop.Relates,
			// FIX: this exists for backwards compatiblity with the client.
			// the reality is, a relation effects a table, there may be multiple views that need updating. either the client could do this by seeing the relation and pulling new data,
			// or we could push all of them. this pushes just one. ( client pulling might be best )
			From:  other.Property,
			IsRev: prop.IsRev,
		}
	default:
		panic(fmt.Sprintf("GetRelative(%s.%s) property does not support relations.", p.src, p.prop.GetId()))
	}
	return
}

// PanicValue implements the Value interface:
// pancing on every get() and set(), and then
// specific property types override the specific methods they need:
// .text for text, num for num, etc.
type panicValue propBase

func (p panicValue) GetNum() float32 {
	panic(fmt.Errorf("get num not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetNum(float32) error {
	panic(fmt.Errorf("set num not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetText() string {
	panic(fmt.Errorf("get text not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetText(string) error {
	panic(fmt.Errorf("set text not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetState() ident.Id {
	panic(fmt.Errorf("get state not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetState(ident.Id) error {
	panic(fmt.Errorf("set state not supported for property %v", p.prop.GetId()))
}
func (p panicValue) GetObject() ident.Id {
	panic(fmt.Errorf("get object not supported for property %v", p.prop.GetId()))
}
func (p panicValue) SetObject(ident.Id) error {
	panic(fmt.Errorf("set object not supported for property %v", p.prop.GetId()))
}

type numValue struct{ panicValue }

func (p numValue) SetNum(f float32) error {
	if p.setValue == nil {
		p.panicValue.SetNum(f)
	}
	return p.setValue(p.prop, f)
}
func (p numValue) GetNum() float32 {
	return p.getValue(p.prop).(float32)
}

type textValue struct{ panicValue }

func (p textValue) GetText() string {
	return p.getValue(p.prop).(string)
}
func (p textValue) SetText(t string) error {
	if p.setValue == nil {
		p.panicValue.SetText(t)
	}
	return p.setValue(p.prop, t)
}

type enumValue struct{ panicValue }

func (p enumValue) GetState() (ret ident.Id) {
	v := p.getValue(p.prop)
	if idx, ok := v.(int); !ok {
		panic(fmt.Sprintf("internal error, couldnt convert state to int '%s.%s' %v(%T)", p.src, p.prop.GetId(), v, v))
	} else {
		enum := p.prop.(M.EnumProperty)
		c, _ := enum.IndexToChoice(idx)
		ret = c
	}
	return
}
func (p enumValue) SetState(c ident.Id) (err error) {
	if p.setValue == nil {
		p.panicValue.SetState(c)
	}
	enum := p.prop.(M.EnumProperty)
	if idx, e := enum.ChoiceToIndex(c); e != nil {
		err = e
	} else {
		p.setValue(p.prop, idx)
	}
	return
}

type pointerValue struct {
	panicValue
}

func (p pointerValue) GetObject() ident.Id {
	return p.getValue(p.prop).(ident.Id)
}
func (p pointerValue) SetObject(o ident.Id) error {
	if p.setValue == nil {
		p.panicValue.SetObject(o)
	}
	return p.setValue(p.prop, o)
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
func singleValue(p propBase) api.Value {
	rel := p.prop.(M.RelativeProperty)
	objs := p.mdl.getObjects(p.src, rel.Relation, rel.IsRev)
	var v ident.Id
	if len(objs) > 0 {
		v = objs[0]
	}
	return objectWriteValue{objectReadValue{panicValue(p), v}}
}

func (p objectWriteValue) SetObject(id ident.Id) (err error) {
	if !id.Empty() {
		err = p.mdl.canAppend(id, p.src, p.prop.(M.RelativeProperty))
	}
	if err == nil {
		p.mdl.clearValues(p.src, p.prop.(M.RelativeProperty))
		p.mdl.appendObject(id, p.src, p.prop.(M.RelativeProperty))
		p.currentVal = id
	}
	return err
}

type objectList struct {
	panicValue
	objs []ident.Id
}

func manyValue(p propBase) api.Values {
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
	if e := p.mdl.canAppend(id, p.src, p.prop.(M.RelativeProperty)); e != nil {
		err = e
	} else {
		p.mdl.appendObject(id, p.src, p.prop.(M.RelativeProperty))
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
