package internal

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

// GameObject wraps Instances for user script callbacks.
// WARNING: for users to test object equality, the GameObject must be comparable;
// it can't implement the interface as a pointer, and it cant have any cached values.
type GameObject struct {
	game *Game // for console, Go(), and relations
	gobj meta.Instance
}

// String helps debugging.
func (oa GameObject) String() string {
	return oa.gobj.GetId().String()
}

// Id uniquely identifies the object.
func (oa GameObject) Id() ident.Id {
	return oa.gobj.GetId()
}

// Exists always returns true for GameObject; see also NullObject which always returns false.
func (oa GameObject) Exists() bool {
	return true
}

// FromClass returns true when the object is compatible with ( based on ) the named class. ( parent or other ancestor )
func (oa GameObject) FromClass(class string) (okay bool) {
	clsid := StripStringId(class)
	return oa.game.Model.AreCompatible(oa.gobj.GetParentClass().GetId(), clsid)
}

func (oa GameObject) ParentRelation() (ret G.IObject, rel string) {
	if parent, prop, ok := oa.game.LookupParent(oa.game.Model, oa.gobj); ok {
		ret, rel = NewGameObject(oa.game, parent), prop.GetName()
	} else {
		ret = NullObjectSource(NewPath(oa.Id()).Add("parent"), 1)
	}
	return ret, rel
}

// Is this object in the passed state?
func (oa GameObject) Is(state string) (ret bool) {
	choice := MakeStringId(state)
	if prop, ok := oa.gobj.GetPropertyByChoice(choice); !ok {
		oa.log("Is(%s): no such choice.", state)
	} else {
		currChoice := prop.GetValue().GetState()
		ret = currChoice == choice
	}
	return ret
}

// IsNow changes the state of an object.
func (oa GameObject) IsNow(state string) {
	newChoice := MakeStringId(state)
	if prop, ok := oa.gobj.GetPropertyByChoice(newChoice); !ok {
		oa.log("IsNow(%s): no such choice.", state)
	} else if e := prop.GetValue().SetState(newChoice); e != nil {
		oa.log("IsNow(%s): error setting value: %s.", state, e)
	}
}

func (oa GameObject) Get(prop string) (ret G.IValue) {
	if p, ok := oa.gobj.FindProperty(prop); !ok {
		oa.log("Get(%s): no such property", prop)
		panic(prop)
		ret = nullValue{}
	} else if p.GetType()&meta.ArrayProperty != 0 {
		oa.log("Get(%s): property is array", prop)
		ret = nullValue{}
	} else {
		ret = gameValue{oa.game, NewPath(p.GetId()), p.GetType(), p.GetValue()}
	}
	return
}

func (oa GameObject) List(prop string) (ret G.IList) {
	if p, ok := oa.gobj.FindProperty(prop); !ok {
		oa.log("List(%s): no such property.", prop)
		ret = nullList{}
	} else if p.GetType()&meta.ArrayProperty == 0 {
		oa.log("List(%s): property is a value, not a list.", prop)
		ret = nullList{}
	} else {
		ret = gameList{oa.game, NewPath(p.GetId()), p.GetType(), p.GetValues()}
	}
	return
}

// Num value of the named property.
func (oa GameObject) Num(prop string) (ret float32) {
	return oa.Get(prop).Num()
}

// SetNum changes the value of an existing number property.
func (oa GameObject) SetNum(prop string, value float32) {
	oa.Get(prop).SetNum(value)
}

// Text value of the named property ( expanding any templated text. )
// ( interestingly, inform seems to error when trying to store or manipulate templated text. )
func (oa GameObject) Text(prop string) (ret string) {
	return oa.Get(prop).Text()
}

// SetText changes the value of an existing text property.
func (oa GameObject) SetText(prop string, text string) {
	oa.Get(prop).SetText(text)
}

// Object returns a related object.
func (oa GameObject) Object(prop string) (ret G.IObject) {
	return oa.Get(prop).Object()
}

// Set changes an object relationship.
func (oa GameObject) Set(prop string, object G.IObject) {
	oa.Get(prop).SetObject(object)
}

// ObjectList returns a list of related objects.
func (oa GameObject) ObjectList(prop string) (ret []G.IObject) {
	if p, ok := oa.gobj.FindProperty(prop); ok {
		switch t := p.GetType(); t {
		default:
			oa.log("ObjectList(%s): invalid type(%d).", prop, t)

		case meta.ObjectProperty | meta.ArrayProperty:
			vals := p.GetValues()
			numobjects := vals.NumValue()
			ret = make([]G.IObject, numobjects)
			for i := 0; i < numobjects; i++ {
				objId := vals.ValueNum(i).GetObject()
				ret[i] = NewGameObjectFromId(oa.game, objId)
			}
		}
	}
	return ret
}

// Says provides this object with a voice.
func (oa GameObject) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, "\n")
	oa.game.Output.ActorSays(oa.gobj, lines)
}

// Go sends all the events associated with the named action,
// and runs the default action if appropriate.
// @see also: Game.ProcessEventQueue
func (oa GameObject) Go(run string, objects ...G.IObject) {
	actionId := MakeStringId(run)
	if action, ok := oa.game.Model.GetAction(actionId); !ok {
		oa.log("Go(%s): no such action", run)
	} else {
		// FIX, ugly: we need the props, even tho we already have the objects...
		nouns := make([]ident.Id, len(objects)+1)
		nouns[0] = oa.Id()
		for i, o := range objects {
			nouns[i+1] = o.Id()
		}
		if act, e := oa.game.NewRuntimeAction(action, nouns...); e != nil {
			oa.log("Go(%s) with %v: error running action: %s", run, objects, e)
		} else {
			tgt := ObjectTarget{oa.game, oa.gobj}
			msg := &E.Message{Id: action.GetEvent().GetId(), Data: act}
			if e := oa.game.SendMessage(tgt, msg); e != nil {
				oa.log("Go(%s): error sending message: %s", run, e)
			}
		}
	}
}

func (oa GameObject) log(format string, v ...interface{}) {
	suffix := fmt.Sprintf(format, v...)
	prefix := oa.gobj.GetId().String()
	oa.game.Println(prefix, suffix)
}
