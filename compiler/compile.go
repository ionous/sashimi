package compiler

import (
	i "github.com/ionous/sashimi/compiler/internal"
	M "github.com/ionous/sashimi/compiler/model"
	X "github.com/ionous/sashimi/compiler/xmodel"
	S "github.com/ionous/sashimi/source"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"io"
	"log"
	"reflect"
)

type Config struct {
	Output io.Writer
}

type converter struct {
	m *M.Model
	x *X.Model
}

func newConverter(x *X.Model) converter {
	m := &M.Model{
		Actions:        make(M.Actions),
		Classes:        make(M.Classes),
		Enumerations:   make(M.Enumerations),
		Events:         make(M.Events),
		Instances:      make(M.Instances),
		Aliases:        make(M.Aliases),
		ParserActions:  make([]M.ParserAction, 0, len(x.ParserActions)),
		Relations:      make(M.Relations),
		SingleToPlural: make(M.SingleToPlural),
	}
	return converter{m, x}
}

// Compile script statements into a "model", a form usable by the runtime.
func (cfg Config) Compile(src S.Statements) (ret *M.Model, err error) {
	names := i.NewNameSource()
	rel := i.NewRelativeFactory(names.NewScope(nil))
	log := log.New(cfg.Output, "compiling: ", log.Lshortfile)
	ctx := &i.Compiler{
		src, names.NewScope(nil),
		i.NewClassFactory(names, names.NewScope("enums"), rel),
		i.NewInstanceFactory(names, log),
		rel,
		log,
	}
	if x, e := ctx.Compile(); e != nil {
		err = e
	} else {
		c := newConverter(x)

		// actionId -> callbackId
		actions := make(map[ident.Id][]M.CallbackModel)
		log.Println("converting handlers", len(x.ActionHandlers))
		for _, handler := range x.ActionHandlers {
			act, callback, useCapture := handler.Action, handler.Callback, handler.UseCapture()
			arr := actions[act]
			// FIX: for now treating target as bubble,
			// really the compiler should hand off a sorted flat list based on three separate groups; target growing in the same direction as after, but distinctly in the middle of things.
			cm := M.CallbackModel{
				ExecuteBlock: callback,
			}
			if !useCapture {
				arr = append(arr, cm)
			} else {
				// prepend:
				arr = append([]M.CallbackModel{cm}, arr...)
			}
			actions[act] = arr
		}
		log.Println("converting actions", len(x.Actions))
		for id, a := range x.Actions {
			c.m.Actions[id] = &M.ActionModel{
				Id:      a.Id,
				Name:    a.ActionName,
				EventId: a.EventId,
				NounTypes: func() (ret []ident.Id) {
					for _, class := range a.NounTypes {
						ret = append(ret, class.Id)
					}
					return
				}(),
				DefaultActions: actions[a.Id],
			}
		}
		log.Println("converting classes", len(x.Classes))
		for id, srcClass := range x.Classes {
			// if srcClass.Constraints.Len() > 0 {
			// 	panic("constraints not implemented")
			// }
			c.m.Classes[id] = &M.ClassModel{
				Id:       srcClass.Id,
				Parents:  parents(srcClass, nil),
				Plural:   srcClass.Plural,
				Singular: srcClass.Singular,
				Properties: func() (props []M.PropertyModel) {
					for _, prop := range srcClass.Properties {
						props = append(props, func() M.PropertyModel {
							ret := propertyBase(prop)
							//
							switch p := prop.(type) {
							case X.PointerProperty:
								ret.Relates = p.Class

							case X.EnumProperty:
								ret.Id = c.makeEnum(&p)

							case X.RelativeProperty:
								ret.Relation = p.Relation
								ret.Relates = p.Relates
							}
							return ret
						}())
					}
					return props
				}(),
			}
		}

		type MapEventCallbacks map[ident.Id]M.EventModelCallbacks
		capture, bubble := make(MapEventCallbacks), make(MapEventCallbacks)
		for _, l := range x.EventListeners {
			e, cb := l.Event, l.ListenerCallback
			var callbacks MapEventCallbacks
			if cb.UseCapture() {
				callbacks = capture
			} else {
				callbacks = bubble
			}
			cm := M.CallbackModel{
				ExecuteBlock: cb.Callback,
			}
			// append
			var arr = callbacks[e]
			arr = append(arr, M.ListenerModel{
				Instance: cb.Instance,
				Class:    cb.Class,
				Callback: cm,
				Options:  eventOptions(cb.Options),
			})
			callbacks[e] = arr
		}
		log.Println("converting events", len(x.Events))
		for id, evt := range x.Events {
			c.m.Events[id] = &M.EventModel{
				Id:   evt.Id,
				Name: evt.EventName,
				// Type: makes perfect sense: the parameters associated with the action.
				//ActionId:  evt.ActionId,
				// FIX: a one to one action/event ratio isnt desirable
				// actions should raise an event, but different actions should be able to raise the same event; events shouldnt know from whence they came
				Capture: capture[id],
				Bubble:  bubble[id],
			}
		}
		log.Println("converting instances", len(x.Instances))
		for id, inst := range x.Instances {
			c.m.Instances[id] = &M.InstanceModel{
				Id:    inst.Id,
				Class: inst.Class.Id,
				Name:  inst.Name,
				Values: func() M.Values {
					ret := make(M.Values)
					for k, v := range inst.Values {
						// the compiler originally stored ints for enums
						// but storing choices (ids) is easier to debug
						// even at the cost of a little more storage
						if i, ok := v.(int); ok {
							if p, ok := inst.Class.GetProperty(k); !ok {
								panic("couldnt find property for integer enum")
							} else {
								enum := p.(X.EnumProperty)
								if c, e := enum.IndexToChoice(i); e != nil {
									panic(e)
								} else {
									v = c
								}
							}
						}
						ret[k] = v
					}
					return ret
				}(),
			}
		}
		for k, v := range x.NounNames {
			names := make(M.RankedStringIds, len(v))
			for i, n := range v {
				names[i] = n
			}
			c.m.Aliases[k] = names
		}
		for _, a := range x.ParserActions {
			c.m.ParserActions = append(c.m.ParserActions, M.ParserAction{
				a.Action,
				a.Commands,
			})
		}
		for id, rel := range x.Relations {
			if rel.Source.Property.Empty() || rel.Dest.Property.Empty() {
				panic(errutil.New("unsupported relation", id, rel.Style, rel.Source.Property, rel.Dest.Property))
			}
			c.m.Relations[id] = &M.RelationModel{
				Id:     rel.Id,
				Name:   rel.Name,
				Source: rel.Source.Property,
				Target: rel.Dest.Property,
				Style:  M.RelationStyle(rel.Style),
			}

			switch rel.Style {
			case X.ManyToMany:
				panic(errutil.New("unsupported relation", id, rel.Style, rel.Source.Property, rel.Dest.Property))
			case X.OneToOne:
				c.flattenTable(rel.Dest.Property)
				c.flattenTable(rel.Source.Property)
			case X.OneToMany:
				c.flattenTable(rel.Dest.Property)
			case X.ManyToOne:
				c.flattenTable(rel.Source.Property)
			}

		}
		for k, v := range x.SingleToPlural {
			c.m.SingleToPlural[k] = v
		}
		ret = c.m
	}
	return
}

// set that one item as a value.
func (c converter) flattenTable(srcProp ident.Id) (err error) {
	for k, v := range c.x.Instances {
		// for all instances of type source
		if p, ok := v.Class.GetProperty(srcProp); ok {
			// 	find the one item they point to
			// 		via the table and their relative property is rev
			if rel, ok := p.(X.RelativeProperty); !ok {
				err = errutil.New("not a relative property?", k, srcProp, reflect.TypeOf(p))
				break
			} else if rel.IsMany {
				err = errutil.New("relative property is many; want one", k, srcProp)
				break
			} else if table, ok := c.x.Tables[rel.Relation]; !ok {
				err = errutil.New("missing table?", rel.Relation)
				break
			} else if lst := table.List(k, rel.IsRev); len(lst) > 1 {
				err = errutil.New("expected at most one item", k, srcProp, "got", lst)
				break
			} else {
				// note: we always set a value, even if its empty.
				// metal would normally panic on getZero for relation values
				// but, because we have a blank value, not an nil value -- getZero wont get called.
				var other ident.Id
				if len(lst) != 0 {
					other = lst[0]
				}
				if dst, ok := c.m.Instances[k]; !ok {
					err = errutil.New("couldnt find target instance", k)
					break
				} else if old, ok := dst.Values[srcProp]; ok {
					err = errutil.New("value being set twice", k, srcProp, "was", old, "now", other)
					break
				} else {
					dst.Values[srcProp] = other
				}
			}
		}
	}
	return
}

func eventOptions(opts X.ListenerOptions) (ret M.ListenerOptions) {
	for _, flag := range []struct {
		x X.ListenerOptions
		m M.ListenerOptions
	}{
		{X.EventCapture, 0}, // we split the list into capture and bubble: dont need this tracked
		{X.EventTargetOnly, M.EventTargetOnly},
		{X.EventQueueAfter, M.EventQueueAfter},
		{X.EventPreventDefault, M.EventPreventDefault},
	} {
		if opts&flag.x != 0 {
			ret |= flag.m
			opts &= ^flag.x
		}
	}
	if opts != 0 {
		panic(errutil.New("uncopied event options remain", opts))
	}
	return
}

func propertyBase(prop X.IProperty) (ret M.PropertyModel) {
	ret.Id = prop.GetId()
	ret.Type = propertyType(prop)
	ret.Name = prop.GetName()
	ret.IsMany = propertyIsMany(prop)
	return ret
}

func propertyType(prop X.IProperty) (ret M.PropertyType) {
	switch prop.(type) {
	case X.NumProperty:
		ret = M.NumProperty
	case X.TextProperty:
		ret = M.TextProperty
	case X.EnumProperty:
		ret = M.EnumProperty
	case X.PointerProperty, X.RelativeProperty:
		ret = M.PointerProperty
	default:
		panic("unknown x-model property type")
	}
	return
}

func propertyIsMany(prop X.IProperty) (ret bool) {
	switch p := prop.(type) {
	case X.NumProperty:
		ret = p.IsMany
	case X.TextProperty:
		ret = p.IsMany
	case X.EnumProperty:
		ret = false
	case X.PointerProperty:
		ret = p.IsMany
	case X.RelativeProperty:
		ret = p.IsMany
	default:
		panic("unknown x-model property type")
	}
	return
}

func parents(cls *X.ClassInfo, list []ident.Id) (ret []ident.Id) {
	if p := cls.Parent; p == nil {
		ret = list
	} else {
		ret = parents(p, append(list, p.Id))
	}
	return
}

func (c converter) makeEnum(src *X.EnumProperty) ident.Id {
	eid := src.Id
	c.m.Enumerations[eid] = &M.EnumModel{
		Choices: func() []ident.Id {
			ret := make([]ident.Id, len(src.Values))
			for i, v := range src.Values {
				ret[i] = v.Id
			}
			return ret
		}()}
	return eid
}

// func (m converter) constrain(clsid ident.Id, cons *X.ConstraintSet) {
// 	cls := c.m.Classes[clsid]
// 	if parent := cls.Parent(); !parent.Empty() && cons.Parent != nil {
// 		constrain(parent, classes, cons.Parent)
// 	}
// 	// changes constraints, which can appear at any level of a class below
// 	// where a property was defined, into new properties.
// 	// it's probably better to have this
// 	for propId, c := range cons.Map {
// 		switch ex := c.(type) {
// 		case *X.EnumConstraint:
// 			cls.Properties.EnumConstraints[pid] = M.EnumConstraint{
// 				Only: ex.Only,
// 				Never: func() (ret []ident.Id) {
// 					for n, _ := range ex.Never {
// 						ret = append(ret, n)
// 					}
// 					return
// 				}(),
// 				Usual:        ex.Usual,
// 				UsuallyLocal: ex.UsuallyLocal,
// 			}
// 		default:
// 			panic("unknown constraint type")
// 		}
// 	}
//}

type MemoryResult struct {
	Model *M.Model
}

func Compile(out io.Writer, src S.Statements) (res MemoryResult, err error) {
	cfg := Config{out}
	if m, e := cfg.Compile(src); e != nil {
		err = e
	} else {
		res = MemoryResult{m}
	}
	return
}
