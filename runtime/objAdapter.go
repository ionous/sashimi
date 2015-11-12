package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

// ObjectAdapter wraps Instances for user script callbacks.
// WARNING: for users to test object equality, the ObjectAdapter must be comparable;
// it can't implement the interface as a pointer, and it cant have any cached values.
type ObjectAdapter struct {
	game *Game // for console, Go(), and relations
	gobj api.Instance
}

// String helps debugging.
func (oa ObjectAdapter) String() string {
	return oa.gobj.GetId().String()
}

// Id uniquely identifies the object.
func (oa ObjectAdapter) Id() ident.Id {
	return oa.gobj.GetId()
}

func (oa ObjectAdapter) Remove() {
	panic("not implemented")
}

// Exists always returns true for ObjectAdapter; see also NullObject which always returns false.
func (oa ObjectAdapter) Exists() bool {
	return true
}

// FromClass returns true when the object is compatible with ( based on ) the named class. ( parent or other ancestor )
func (oa ObjectAdapter) FromClass(class string) (okay bool) {
	clsid := StripStringId(class)
	return oa.game.ModelApi.AreCompatible(oa.gobj.GetParentClass().GetId(), clsid)
}

// Is this object in the passed state?
func (oa ObjectAdapter) Is(state string) (ret bool) {
	choice := MakeStringId(state)
	if prop, ok := oa.gobj.GetPropertyByChoice(choice); !ok {
		oa.log("%s.Is(%s): no such choice.", oa, state)
	} else {
		currChoice := prop.GetValue().GetState()
		ret = currChoice == choice
	}
	return ret
}

// IsNow changes the state of an object.
func (oa ObjectAdapter) IsNow(state string) {
	newChoice := MakeStringId(state)
	if prop, ok := oa.gobj.GetPropertyByChoice(newChoice); !ok {
		oa.log("%s.IsNow(%s): no such choice.", oa, state)
	} else if e := prop.GetValue().SetState(newChoice); e != nil {
		oa.log("%s.IsNow(%s): error setting value: %s.", oa, state, e)
	}
}

// Num value of the named property.
func (oa ObjectAdapter) Num(prop string) (ret float32) {
	pid := MakeStringId(prop)
	if prop, ok := oa.gobj.GetProperty(pid); !ok {
		oa.log("%s.Num(%s): no such property", oa, prop)
	} else if t := prop.GetType(); t != api.NumProperty {
		oa.log("%s.Num(%s): property type(%d) is not a number.", oa, prop, t)
	} else {
		ret = prop.GetValue().GetNum()
	}
	return ret
}

// SetNum changes the value of an existing number property.
func (oa ObjectAdapter) SetNum(prop string, value float32) {
	pid := MakeStringId(prop)
	if prop, ok := oa.gobj.GetProperty(pid); !ok {
		oa.log("%s.SetNum(%s): no such property", oa, prop)
	} else if t := prop.GetType(); t != api.NumProperty {
		oa.log("%s.SetNum(%s): property type(%d) is not a number.", oa, prop, t)
	} else if e := prop.GetValue().SetNum(value); e != nil {
		oa.log("%s.SetNum(%s): error setting value: %s.", oa, prop, e)
	}
}

// Text value of the named property ( expanding any templated text. )
// ( interestingly, inform seems to error when trying to store or manipulate templated text. )
func (oa ObjectAdapter) Text(prop string) (ret string) {
	pid := MakeStringId(prop)
	if prop, ok := oa.gobj.GetProperty(pid); !ok {
		oa.log("%s.Text(%s): no such property", oa, prop)
	} else if t := prop.GetType(); t != api.TextProperty {
		oa.log("%s.Text(%s): property type(%d) is not text.", oa, prop, t)
	} else {
		ret = prop.GetValue().GetText()
	}
	return ret
}

// SetText changes the value of an existing text property.
func (oa ObjectAdapter) SetText(prop string, text string) {
	pid := MakeStringId(prop)
	if prop, ok := oa.gobj.GetProperty(pid); !ok {
		oa.log("%s.SetText(%s): no such property", oa, prop)
	} else if t := prop.GetType(); t != api.TextProperty {
		oa.log("%s.SetText(%s): property type(%d) is not text.", oa, prop, t)
	} else if e := prop.GetValue().SetText(text); e != nil {
		oa.log("%s.SetText(%s): error setting value: %s.", oa, prop, e)
	}
}

// Object returns a related object.
func (oa ObjectAdapter) Object(prop string) (ret G.IObject) {
	// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
	var res ident.Id
	pid := MakeStringId(prop)
	if p, ok := oa.gobj.GetProperty(pid); ok {
		if t := p.GetType(); t != api.ObjectProperty {
			oa.log("%s.Object(%s): property type(%d) is not an object.", oa, prop, t)
		} else {
			res = p.GetValue().GetObject()
		}
	}
	return NewObjectAdapterFromId(oa.game, res)
}

// Set changes an object relationship.
func (oa ObjectAdapter) Set(prop string, object G.IObject) {
	pid := MakeStringId(prop)
	if p, ok := oa.gobj.GetProperty(pid); ok {
		switch t := p.GetType(); t {
		default:
			oa.log("%s.Set(%s): property type(%d) is not an object.", oa, prop, t)

		case api.ObjectProperty:
			var id ident.Id
			if other, ok := object.(ObjectAdapter); ok {
				id = other.gobj.GetId()
			}
			if e := p.GetValue().SetObject(id); e != nil {
				oa.log("%s.Set(%s): error setting value: %s.", oa, prop, e)
			}

		case api.ObjectProperty | api.ArrayProperty:
			values := p.GetValues()
			if other, ok := object.(ObjectAdapter); !ok {
				values.ClearValues()
			} else {
				if e := values.AppendObject(other.gobj.GetId()); e != nil {
					oa.log("%s.Set(%s): error appending value: %s.", oa, prop, e)
				}
			}
		}
	}
}

// ObjectList returns a list of related objects.
func (oa ObjectAdapter) ObjectList(prop string) (ret []G.IObject) {
	if p, ok := oa.gobj.GetProperty(MakeStringId(prop)); ok {
		switch t := p.GetType(); t {
		default:
			oa.log("%s.ObjectList(%s): invalid type(%d).", oa, prop, t)

		case api.ObjectProperty | api.ArrayProperty:
			vals := p.GetValues()
			numobjects := vals.NumValue()
			ret = make([]G.IObject, numobjects)
			for i := 0; i < numobjects; i++ {
				objId := vals.ValueNum(i).GetObject()
				ret[i] = NewObjectAdapterFromId(oa.game, objId)
			}
		}
	}
	return ret
}

// Says provides this object with a voice.
func (oa ObjectAdapter) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, "\n")
	oa.game.output.ActorSays(oa.gobj, lines)
}

// Go sends all the events associated with the named action,
// and runs the default action if appropriate.
// @see also: Game.ProcessEventQueue
func (oa ObjectAdapter) Go(run string, objects ...G.IObject) {
	actionId := MakeStringId(run)
	if action, ok := oa.game.ModelApi.GetAction(actionId); !ok {
		oa.log("%s.Go(%s): no such action", oa, run)
	} else {
		// FIX, ugly: we need the props, even tho we already have the objects...
		nouns := make([]ident.Id, len(objects)+1)
		nouns[0] = oa.Id()
		for i, o := range objects {
			nouns[i+1] = o.Id()
		}
		if act, e := oa.game.newRuntimeAction(action, nouns...); e != nil {
			oa.log("%s.Go(%s) with %v: error running action: %s", oa, run, objects, e)
		} else {
			tgt := ObjectTarget{oa.game, oa.gobj}
			msg := &E.Message{Id: action.GetEvent().GetId(), Data: act}
			if e := oa.game.frame.SendMessage(tgt, msg); e != nil {
				oa.log("%s.Go(%s): error sending message: %s", oa, run, e)
			}
		}
	}
}

//
func (oa ObjectAdapter) log(format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	oa.game.log.Output(4, msg)
}
