package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
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

func (oa ObjectAdapter) GetGameObject() *GameObject {
	return oa.gobj
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
func (oa ObjectAdapter) FromClass(class string) (okay bool) {
	clsid := StripStringId(class)
	return oa.game.ModelApi.AreCompatible(oa.gobj.cls.GetId(), clsid)
}

//
// Is this object in the passed state?
//
func (oa ObjectAdapter) Is(state string) (ret bool) {
	choice := MakeStringId(state)
	if prop, ok := oa.gobj.cls.GetPropertyByChoice(choice); !ok {
		oa.logError(fmt.Errorf("is: no such choice '%s'.'%s'", oa, state))
	} else {
		currChoice := oa.gobj.Value(prop.GetId())
		ret = currChoice == choice
	}
	return ret
}

//
// IsNow changes the state of an object.
//
func (oa ObjectAdapter) IsNow(state string) {
	newChoice := MakeStringId(state)
	if prop, ok := oa.gobj.cls.GetPropertyByChoice(newChoice); !ok {
		oa.logError(fmt.Errorf("IsNow: no such choice '%s'.'%s'", oa, state))
	} else {
		// get the current choice from the implied property slot
		if currChoice, ok := oa.gobj.Value(prop.GetId()).(ident.Id); !ok {
			err := TypeMismatch(oa.gobj.Id().String(), prop.GetId().String())
			oa.logError(err)
		} else {
			if currChoice != newChoice {
				oa.gobj.removeDirect(currChoice)           // delete the old choice's boolean,
				oa.gobj.setDirect(newChoice, true)         // and set the new
				oa.gobj.setDirect(prop.GetId(), newChoice) // // set the property slot to the new choice
				oa.game.Properties.VisitWatchers(func(ch PropertyChange) {
					ch.StateChange(oa.gobj, prop.GetId(), currChoice, newChoice)
				})
			}
		}
	}
}

//
// Num value of the named property.
//
func (oa ObjectAdapter) Num(prop string) (ret float32) {
	id := MakeStringId(prop)
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
	id := MakeStringId(prop)
	if old, ok := oa.gobj.SetValue(id, value); !ok {
		oa.logError(TypeMismatch(prop, "set num"))
	} else {
		old := old.(float32)
		oa.game.Properties.VisitWatchers(func(ch PropertyChange) {
			ch.NumChange(oa.gobj, id, old, value)
		})
	}
}

//
// Text value of the named property ( expanding any templated text. )
// ( interestingly, inform seems to error when trying to store or manipulate templated text. )
//
func (oa ObjectAdapter) Text(prop string) (ret string) {
	id := MakeStringId(prop)
	if val, ok := oa.gobj.Value(id).(string); !ok {
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
	id := MakeStringId(prop)
	if old, ok := oa.gobj.SetValue(id, text); !ok {
		oa.logError(TypeMismatch(prop, "set text"))
	} else {
		old := old.(string)
		oa.game.Properties.VisitWatchers(func(ch PropertyChange) {
			ch.TextChange(oa.gobj, id, old, text)
		})
	}
}

//
// Object returns a related object.
//
func (oa ObjectAdapter) Object(prop string) (ret G.IObject) {
	// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
	var res ident.Id
	id := MakeStringId(prop)
	if p, ok := oa.gobj.cls.GetProperty(id); ok {
		if t := p.GetType(); t == api.ObjectProperty {
			if val, ok := oa.gobj.vals[p.GetId().String()]; ok {
				res = val.(ident.Id)
			}
		}
	}
	return NewObjectAdapter(oa.game, oa.game.Objects[res])
}

//
// Set changes an object relationship.
//
func (oa ObjectAdapter) Set(prop string, object G.IObject) {
	propId := MakeStringId(prop)
	if p, ok := oa.gobj.cls.GetProperty(propId); ok {
		switch t := p.GetType(); t {
		default:
			oa.logError(TypeMismatch(oa.String(), prop))

		case api.ObjectProperty:
			var id ident.Id
			if other, ok := object.(ObjectAdapter); ok {
				id = other.gobj.Id()
			}
			oa.gobj.SetValue(propId, id)
			//
			i, _ := oa.game.ModelApi.GetInstance(oa.gobj.Id())
			p, _ := i.GetProperty(propId)
			p.GetValue().SetObject(id)

		case api.ObjectProperty | api.ArrayProperty:
			// STORE FIX: this should come from the object
			i, _ := oa.game.ModelApi.GetInstance(oa.gobj.Id())
			p, _ := i.GetProperty(propId)
			values := p.GetValues()
			if other, ok := object.(ObjectAdapter); !ok {
				values.ClearValues()
			} else {
				if e := values.AppendObject(other.gobj.Id()); e != nil {
					oa.logError(e)
				}
			}
		}
	}
}

//
// ObjectList returns a list of related objects.
//
func (oa ObjectAdapter) ObjectList(prop string) (ret []G.IObject) {
	if p, ok := oa.gobj.cls.GetProperty(MakeStringId(prop)); ok {
		switch p.GetType() {
		default:
			oa.logError(TypeMismatch(oa.String(), prop))

		case api.ObjectProperty | api.ArrayProperty:
			vals := p.GetValues()
			l := vals.NumValue()
			ret = make([]G.IObject, l)
			for i := 0; i < l; i++ {
				objId := vals.ValueNum(i).GetObject()
				ret[i] = NewObjectAdapter(oa.game, oa.game.Objects[objId])
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
	actionId := MakeStringId(act)
	if action, ok := oa.game.ModelApi.GetAction(actionId); !ok {
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
			msg := &E.Message{Name: action.GetEvent().GetEventName(), Data: act}
			if e := oa.game.frame.SendMessage(tgt, msg); e != nil {
				oa.logError(e)
			}
		}
	}
}

//
func (oa ObjectAdapter) logError(err error) (hadError bool) {
	if err != nil {
		oa.game.log.Output(4, fmt.Sprint("!!!Error:", err.Error()))
		hadError = true
	}
	return hadError
}
