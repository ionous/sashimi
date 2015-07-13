package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

//
// ObjectAdapter wraps GameObject(s) for user script callbacks.
// WARNING: for users to test object equality, the ObjectAdapter must be comparable;
// it can't implement the interface as a pointer, and it cant have any cached values.
//
type ObjectAdapter struct {
	game *Game // for console, Go(), and relations
	gobj *GameObject
}

//
// NewObjectAdapter gives the passed game object the IObject interface.
//
func NewObjectAdapter(game *Game, gobj *GameObject) ObjectAdapter {
	return ObjectAdapter{game, gobj}
}

//
// String helps debugging.
//
func (adapt ObjectAdapter) String() string {
	return adapt.gobj.Id().String()
}

//
// Id uniquely identifies the object.
//
func (adapt ObjectAdapter) Id() ident.Id {
	return adapt.gobj.Id()
}

//
// Name of the object.
//
// func (adapt ObjectAdapter) Name() string {
// 	return adapt.gobj.inst.Name()
// }

//
// Exists always returns true for ObjectAdapter; see also NullObject which always returns false.
//
func (adapt ObjectAdapter) Exists() bool {
	return true
}

//
// Class returns true when this object is compatible with ( based on ) the named class. ( parent or other ancestor )
//
func (adapt ObjectAdapter) Class(class string) (okay bool) {
	if cls, ok := adapt.game.Model.Classes.FindClassBySingular(class); ok {
		okay = adapt.gobj.Class().CompatibleWith(cls.Id())
	}
	return okay
}

//
// Is this object in the passed state?
//
func (adapt ObjectAdapter) Is(state string) (ret bool) {
	if prop, index, ok := adapt.gobj.Class().PropertyByChoice(state); !ok {
		adapt.logError(fmt.Errorf("is: no such choice '%s'.'%s'", adapt, state))
	} else {
		testChoice, _ := prop.IndexToChoice(index)
		currChoice := adapt.gobj.Value(prop.Id())
		ret = currChoice == testChoice
	}
	return ret
}

//
// SetIs changes the state of an object.
//
func (adapt ObjectAdapter) SetIs(state string) {
	if prop, index, ok := adapt.gobj.Class().PropertyByChoice(state); !ok {
		adapt.logError(fmt.Errorf("SetIs: no such choice '%s'.'%s'", adapt, state))
	} else {
		// get the current choice from the implied property slot
		if currChoice, ok := adapt.gobj.Value(prop.Id()).(ident.Id); !ok {
			err := TypeMismatch(adapt.gobj.Id().String(), prop.Id().String())
			adapt.logError(err)
		} else {
			newChoice, _ := prop.IndexToChoice(index)
			if currChoice != newChoice {
				adapt.gobj.removeDirect(currChoice)        // delete the old choice's boolean,
				adapt.gobj.setDirect(newChoice, true)      // and set the new
				adapt.gobj.setDirect(prop.Id(), newChoice) // // set the property slot to the new choice
				adapt.game.Properties.Notify(adapt.gobj.Id(), prop.Id(), currChoice, newChoice)
			}
		}
	}
}

//
// Num value of the named property.
//
func (adapt ObjectAdapter) Num(prop string) (ret float32) {
	id := M.MakeStringId(prop)
	if val, ok := adapt.gobj.Value(id).(float32); !ok {
		adapt.logError(TypeMismatch(prop, "get num"))
	} else {
		ret = val
	}
	return ret
}

//
// SetNum changes the value of an existing number property.
//
func (adapt ObjectAdapter) SetNum(prop string, value float32) {
	id := M.MakeStringId(prop)
	if old, ok := adapt.gobj.SetValue(id, value); !ok {
		adapt.logError(TypeMismatch(prop, "set num"))
	} else {
		adapt.game.Properties.Notify(adapt.gobj.Id(), id, old, value)
	}
}

//
// Text value of the named property ( expanding any templated text. )
// ( interestingly, inform seems to error when trying to store or manipulate templated text. )
//
func (adapt ObjectAdapter) Text(prop string) (ret string) {
	id := M.MakeStringId(prop)
	// is adapt text stored as a template?
	if temp, ok := adapt.gobj.temps[id.String()]; ok {
		if s, e := runTemplate(temp, adapt.gobj.vals); e != nil {
			adapt.logError(e)
		} else {
			ret = s
		}
	} else if val, ok := adapt.gobj.Value(id).(string); !ok {
		adapt.logError(TypeMismatch(prop, "get text"))
	} else {
		ret = val
	}
	return ret
}

//
// SetText changes the value of an existing text property.
//
func (adapt ObjectAdapter) SetText(prop string, text string) {
	id := M.MakeStringId(prop)
	if e := adapt.gobj.temps.New(id.String(), text); e != nil {
		adapt.logError(e)
	} else if old, ok := adapt.gobj.SetValue(id, text); !ok {
		adapt.logError(TypeMismatch(prop, "set text"))
	} else {
		adapt.game.Properties.Notify(adapt.gobj.Id(), id, old, text)
	}
}

//
// Object returns a related object.
//
func (adapt ObjectAdapter) Object(prop string) (ret G.IObject) {
	// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
	var res ident.Id
	if p, ok := adapt.gobj.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		case *M.PointerProperty:
			if val, ok := adapt.gobj.Value(p.Id()).(ident.Id); ok {
				res = val
			}
		case *M.RelativeProperty:
			// TBD: can the relative property changes automatically reflect into the value table
			// ex. on event?
			if rel, ok := adapt.gobj.Value(p.Id()).(RelativeValue); ok {
				if p.ToMany() {
					adapt.logError(fmt.Errorf("object requested, but relation is list"))
				} else {
					list := rel.List()
					if len(list) != 0 {
						res = ident.Id(list[0])
					}
				}
			}
		}
	}
	if gobj, ok := adapt.game.Objects[res]; ok {
		ret = ObjectAdapter{adapt.game, gobj}
	} else {
		ret = NullObject{}
	}
	return ret
}

//
// Set changes an object relationship.
//
func (adapt ObjectAdapter) Set(prop string, object G.IObject) {
	if p, ok := adapt.gobj.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		default:
			adapt.logError(TypeMismatch(adapt.String(), prop))
		case *M.PointerProperty:
			set := false
			if other, ok := object.(ObjectAdapter); !ok {
				adapt.gobj.SetValue(p.Id(), ident.Id(""))
				set = true
			} else {
				adapt.gobj.SetValue(p.Id(), other.gobj.Id())
				set = true
			}
			if !set {
				adapt.logError(fmt.Errorf("couldnt set value for prop %s", prop))
			}
		case *M.RelativeProperty:
			if rel, ok := adapt.gobj.Value(p.Id()).(RelativeValue); ok {

				// if the referenced object doesnt exist, we take it to mean they are clearing the reference.
				if other, ok := object.(ObjectAdapter); !ok {
					if removed, e := rel.ClearReference(); e != nil {
						adapt.logError(e)
					} else {
						adapt.game.Properties.Notify(adapt.gobj.Id(), p.Id(), removed, ident.Empty())
					}
				} else {
					// FIX? the impedence b/t IObject and Reference is annoying.
					other := other.gobj.Id()
					if ref, ok := adapt.game.Model.Instances[other]; !ok {
						adapt.logError(fmt.Errorf("Set: couldnt find object names %s", other))
					} else if removed, e := rel.SetReference(ref); e != nil {
						adapt.logError(e)
					} else {
						// removed is probably a single object
						adapt.game.Properties.Notify(adapt.gobj.Id(), p.Id(), removed, other)
					}
				}
			}
		}
	}
}

//
// ObjectList returns a list of related objects.
//
func (adapt ObjectAdapter) ObjectList(prop string) (ret []G.IObject) {
	if p, ok := adapt.gobj.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		default:
			adapt.logError(TypeMismatch(adapt.String(), prop))
		case *M.RelativeProperty:
			if rel, ok := adapt.gobj.Value(p.Id()).(RelativeValue); ok {
				list := rel.List()
				ret = make([]G.IObject, len(list))
				for i, objId := range list {
					if gobj, ok := adapt.game.Objects[objId]; ok {
						ret[i] = ObjectAdapter{adapt.game, gobj}
					} else {
						ret[i] = NullObject{}
					}
				}
			}
		}
	}
	return ret
}

//
// Says provides this object with a voice.
//
func (adapt ObjectAdapter) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, "\n")
	adapt.game.output.ActorSays(adapt.gobj, lines)
}

//
// Go sends all the events associated with the named action,
// and runs the default action if appropriate.
// @see also: Game.ProcessEventQueue
//
func (adapt ObjectAdapter) Go(act string, objects ...G.IObject) {
	if action, ok := adapt.game.Model.Actions.FindActionByName(act); !ok {
		adapt.logError(fmt.Errorf("unknown action for Go %s", act))
	} else {
		// FIX, ugly: we need the props, even tho we already have the objects...
		nouns := make([]string, len(objects)+1)
		nouns[0] = adapt.Id().String()
		for i, o := range objects {
			nouns[i+1] = o.Id().String()
		}
		if act, e := adapt.game.newRuntimeAction(action, nouns...); e != nil {
			adapt.logError(e)
		} else {
			tgt := ObjectTarget{adapt.game, adapt.gobj}
			msg := E.Message{Name: action.Event(), Data: act}
			// see ProcessEventQueue()
			path := E.NewPathTo(tgt)
			//adapt.game.log.Output(3, fmt.Sprintf("go %s %s", prop, path))
			if runDefault, err := msg.Send(path); err != nil {
				adapt.logError(err)
			} else if runDefault {
				act.runDefaultActions()
			}
		}
	}
}

//
//
//
func (adapt ObjectAdapter) logError(err error) (hadError bool) {
	if err != nil {
		adapt.game.log.Output(4, fmt.Sprint("!!!Error:", err.Error()))
		hadError = true
	}
	return hadError
}
