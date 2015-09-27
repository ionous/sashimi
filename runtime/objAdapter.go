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
// String helps debugging.
//
func (oa ObjectAdapter) String() string {
	return oa.gobj.Id().String()
}

//
// Id uniquely identifies the object.
//
func (oa ObjectAdapter) Id() ident.Id {
	return oa.gobj.Id()
}

//
func (oa ObjectAdapter) Remove() {
	delete(oa.game.Objects, oa.gobj.Id())
}

//
// Name of the object.
//
// func (oa ObjectAdapter) Name() string {
// 	return oa.gobj.inst.Name()
// }

//
// Exists always returns true for ObjectAdapter; see also NullObject which always returns false.
//
func (oa ObjectAdapter) Exists() bool {
	return true
}

//
// Class returns true when this object is compatible with ( based on ) the named class. ( parent or other ancestor )
//
func (oa ObjectAdapter) Class(class string) (okay bool) {
	if cls, ok := oa.game.Model.Classes.FindClassBySingular(class); ok {
		okay = oa.gobj.Class().CompatibleWith(cls.Id())
	}
	return okay
}

//
// Is this object in the passed state?
//
func (oa ObjectAdapter) Is(state string) (ret bool) {
	if prop, index, ok := oa.gobj.Class().PropertyByChoice(state); !ok {
		oa.logError(fmt.Errorf("is: no such choice '%s'.'%s'", oa, state))
	} else {
		testChoice, _ := prop.IndexToChoice(index)
		currChoice := oa.gobj.Value(prop.Id())
		ret = currChoice == testChoice
	}
	return ret
}

//
// IsNow changes the state of an object.
//
func (oa ObjectAdapter) IsNow(state string) {
	if prop, index, ok := oa.gobj.Class().PropertyByChoice(state); !ok {
		oa.logError(fmt.Errorf("IsNow: no such choice '%s'.'%s'", oa, state))
	} else {
		// get the current choice from the implied property slot
		if currChoice, ok := oa.gobj.Value(prop.Id()).(ident.Id); !ok {
			err := TypeMismatch(oa.gobj.Id().String(), prop.Id().String())
			oa.logError(err)
		} else {
			newChoice, _ := prop.IndexToChoice(index)
			if currChoice != newChoice {
				oa.gobj.removeDirect(currChoice)        // delete the old choice's boolean,
				oa.gobj.setDirect(newChoice, true)      // and set the new
				oa.gobj.setDirect(prop.Id(), newChoice) // // set the property slot to the new choice
				oa.game.Properties.Notify(oa.gobj.Id(), prop.Id(), currChoice, newChoice)
			}
		}
	}
}

//
// Num value of the named property.
//
func (oa ObjectAdapter) Num(prop string) (ret float32) {
	id := M.MakeStringId(prop)
	if val, ok := oa.gobj.Value(id).(float32); !ok {
		oa.logError(TypeMismatch(prop, "get num"))
	} else {
		ret = val
	}
	return ret
}

//
// SetNum changes the value of an existing number property.
//
func (oa ObjectAdapter) SetNum(prop string, value float32) {
	id := M.MakeStringId(prop)
	if old, ok := oa.gobj.SetValue(id, value); !ok {
		oa.logError(TypeMismatch(prop, "set num"))
	} else {
		oa.game.Properties.Notify(oa.gobj.Id(), id, old, value)
	}
}

//
// Text value of the named property ( expanding any templated text. )
// ( interestingly, inform seems to error when trying to store or manipulate templated text. )
//
func (oa ObjectAdapter) Text(prop string) (ret string) {
	id := M.MakeStringId(prop)
	// is oa text stored as a template?
	if temp, ok := oa.gobj.temps[id.String()]; ok {
		if s, e := runTemplate(temp, oa.gobj.vals); e != nil {
			oa.logError(e)
		} else {
			ret = s
		}
	} else if val, ok := oa.gobj.Value(id).(string); !ok {
		oa.logError(TypeMismatch(prop, "get text"))
	} else {
		ret = val
	}
	return ret
}

//
// SetText changes the value of an existing text property.
//
func (oa ObjectAdapter) SetText(prop string, text string) {
	id := M.MakeStringId(prop)
	if e := oa.gobj.temps.New(id.String(), text); e != nil {
		oa.logError(e)
	} else if old, ok := oa.gobj.SetValue(id, text); !ok {
		oa.logError(TypeMismatch(prop, "set text"))
	} else {
		oa.game.Properties.Notify(oa.gobj.Id(), id, old, text)
	}
}

//
// Object returns a related object.
//
func (oa ObjectAdapter) Object(prop string) (ret G.IObject) {
	// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
	var res ident.Id
	if p, ok := oa.gobj.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		case *M.PointerProperty:
			if val, ok := oa.gobj.Value(p.Id()).(ident.Id); ok {
				res = val
			}
		case *M.RelativeProperty:
			// TBD: can the relative property changes automatically reflect into the value table
			// ex. on event?
			if rel, ok := oa.gobj.Value(p.Id()).(RelativeValue); ok {
				if p.ToMany() {
					oa.logError(fmt.Errorf("object requested, but relation is list"))
				} else {
					list := rel.List()
					if len(list) != 0 {
						res = ident.Id(list[0])
					}
				}
			}
		}
	}
	return NewObjectAdapter(oa.game, oa.game.Objects[res])
}

//
// Set changes an object relationship.
//
func (oa ObjectAdapter) Set(prop string, object G.IObject) {
	if p, ok := oa.gobj.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		default:
			oa.logError(TypeMismatch(oa.String(), prop))
		case *M.PointerProperty:
			set := false
			if other, ok := object.(ObjectAdapter); !ok {
				oa.gobj.SetValue(p.Id(), ident.Id(""))
				set = true
			} else {
				oa.gobj.SetValue(p.Id(), other.gobj.Id())
				set = true
			}
			if !set {
				oa.logError(fmt.Errorf("couldnt set value for prop %s", prop))
			}
		case *M.RelativeProperty:
			if rel, ok := oa.gobj.Value(p.Id()).(RelativeValue); ok {

				// if the referenced object doesnt exist, we take it to mean they are clearing the reference.
				if other, ok := object.(ObjectAdapter); !ok {
					if removed, e := rel.ClearReference(); e != nil {
						oa.logError(e)
					} else {
						oa.game.Properties.Notify(oa.gobj.Id(), p.Id(), removed, ident.Empty())
					}
				} else {
					// FIX? the impedence b/t IObject and Reference is annoying.
					other := other.gobj.Id()
					if ref, ok := oa.game.Model.Instances[other]; !ok {
						oa.logError(fmt.Errorf("Set: couldnt find object names %s", other))
					} else if removed, e := rel.SetReference(ref); e != nil {
						oa.logError(e)
					} else {
						// removed is probably a single object
						oa.game.Properties.Notify(oa.gobj.Id(), p.Id(), removed, other)
					}
				}
			}
		}
	}
}

//
// ObjectList returns a list of related objects.
//
func (oa ObjectAdapter) ObjectList(prop string) (ret []G.IObject) {
	if p, ok := oa.gobj.Class().FindProperty(prop); ok {
		switch p := p.(type) {
		default:
			oa.logError(TypeMismatch(oa.String(), prop))
		case *M.RelativeProperty:
			if rel, ok := oa.gobj.Value(p.Id()).(RelativeValue); ok {
				list := rel.List()
				ret = make([]G.IObject, len(list))
				for i, objId := range list {
					ret[i] = NewObjectAdapter(oa.game, oa.game.Objects[objId])
				}
			}
		}
	}
	return ret
}

//
// Says provides this object with a voice.
//
func (oa ObjectAdapter) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, "\n")
	oa.game.output.ActorSays(oa.gobj, lines)
}

//
// Go sends all the events associated with the named action,
// and runs the default action if appropriate.
// @see also: Game.ProcessEventQueue
//
func (oa ObjectAdapter) Go(act string, objects ...G.IObject) {
	if action, ok := oa.game.Model.Actions.FindActionByName(act); !ok {
		e := fmt.Errorf("unknown action for Go %s", act)
		oa.logError(e)
	} else {
		// FIX, ugly: we need the props, even tho we already have the objects...
		nouns := make([]ident.Id, len(objects)+1)
		nouns[0] = oa.Id()
		for i, o := range objects {
			nouns[i+1] = o.Id()
		}
		if act, e := oa.game.newRuntimeAction(action, nouns...); e != nil {
			oa.logError(e)
		} else {
			tgt := ObjectTarget{oa.game, oa.gobj}
			msg := &E.Message{Name: action.Event(), Data: act}
			if e := oa.game.frame.SendMessage(tgt, msg); e != nil {
				oa.logError(e)
			}
		}
	}
}

//
//
//
func (oa ObjectAdapter) logError(err error) (hadError bool) {
	if err != nil {
		oa.game.log.Output(4, fmt.Sprint("!!!Error:", err.Error()))
		hadError = true
	}
	return hadError
}
