package internal

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// GameObject wraps Instances for user script s.
// WARNING: for users to test object equality, the GameObject must be comparable;
// it can't implement the interface as a pointer, and it cant have any cached values.
type GameObject struct {
	*GameEventAdapter // for console, Go(), and relations
	gobj              meta.Instance
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
	if _, found := oa.Model.GetClass(clsid); !found {
		oa.Println("FromClass: no such class found", clsid)
	}
	return oa.Model.AreCompatible(oa.gobj.GetParentClass(), clsid)
}

func (oa GameObject) ParentRelation() (ret G.IObject, rel string) {
	if parent, prop, ok := oa.LookupParent(oa.gobj); ok {
		ret, rel = oa.NewGameObject(parent), prop.GetName()
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
		ret = nullValue{}
	} else if p.GetType()&meta.ArrayProperty != 0 {
		oa.log("Get(%s): property is array", prop)
		ret = nullValue{}
	} else {
		ret = gameValue{oa.GameEventAdapter,
			NewPath(p.GetId()), p.GetType(), p.GetValue()}
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
		ret = gameList{oa.GameEventAdapter, NewPath(p.GetId()), p.GetType(), p.GetValues()}
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
				ret[i] = oa.NewGameObjectFromId(objId)
			}
		}
	}
	return ret
}

// Says provides this object with a voice.
func (oa GameObject) Says(text string) {
	// FIX: share some template love with GameEventAdapter.Say()
	lines := strings.Split(text, lang.NewLine)
	oa.Output.ActorSays(oa.gobj, lines)
}

// Go sends all the events associated with the named action,
// and runs the default action if appropriate.
// ex. g.The("player").Go("show to", "the alien boy", "the ring")
func (oa GameObject) Go(run string, objects ...G.IObject) (ret G.IPromise) {
	if c, e := oa.queueNamedAction(run, objects); e != nil {
		oa.log("Go(%s) with %v: error preparing action: %s", run, objects, e)
		ret = NilPromise{}
	} else {
		ret = c
	}
	return
}

type NilPromise struct{}

func (NilPromise) Then(G.Callback) {}

// FIX: other variants of this exist in runtime.Game
func (oa GameObject) queueNamedAction(action string, objects []G.IObject) (ret G.IPromise, err error) {
	// FUTURE: ast introspection to find whether the action exists..
	actionId := MakeStringId(action)
	if act, ok := oa.Model.GetAction(actionId); !ok {
		err = fmt.Errorf("couldnt find action %s", action)
	} else {
		// FIX, ugly: we need the ids, even tho we already have the objects..
		// FIXIXFIFX: if that's the case -- just use the raw strings externally.
		nouns := make([]ident.Id, len(objects)+1)
		nouns[0] = oa.Id()
		for i, o := range objects {
			nouns[i+1] = o.Id()
		}
		// this verifies that the objects exist
		if data, e := oa.NewRuntimeAction(act, nouns...); e != nil {
			err = e
		} else {
			future := &QueuedAction{data: data}
			oa.Queue.QueueFuture(future)
			// NOTE: the next callbacks get *our* context, not the context of the action.
			ret = NewPendingChain(oa.GameEventAdapter, future)
		}
	}
	return
}

func (oa GameObject) log(format string, v ...interface{}) {
	suffix := fmt.Sprintf(format, v...)
	prefix := oa.gobj.GetId().String()
	oa.Println(prefix, suffix)
}
