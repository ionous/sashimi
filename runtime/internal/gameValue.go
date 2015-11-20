package internal

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FIX -- can two properties be directly compared ==? it sure would be nice.
// would have to change / get rid of log
// replace with value / path.
// possibly replace Game with a shared "context",  possibly the game object adapter
// so that the log can still have its stack/depth unwinding
type gameValue struct {
	game  *Game
	path  PropertyPath
	ptype meta.PropertyType
	value meta.Value
}

func (n gameValue) Num() (ret float32) {
	if n.ptype != meta.NumProperty {
		n.log("Num(): property is not a number.")
	} else {
		ret = n.value.GetNum()
	}
	return
}

func (n gameValue) SetNum(value float32) {
	if n.ptype != meta.NumProperty {
		n.log("SetNum(): property is not a number.")
	} else if e := n.value.SetNum(value); e != nil {
		n.log("SetNum(): error setting value: %v.", e)
	}
}

func (n gameValue) Text() (ret string) {
	if n.ptype != meta.TextProperty {
		n.log("Text(): property is not text.")
	} else {
		ret = n.value.GetText()
	}
	return
}

func (n gameValue) SetText(text string) {
	if n.ptype != meta.TextProperty {
		n.log("SetText(): property is not text.")
	} else if e := n.value.SetText(text); e != nil {
		n.log("SetText(): error setting value: %v.", e)
	}
}

// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
func (n gameValue) Object() G.IObject {
	var res ident.Id
	if n.ptype != meta.ObjectProperty {
		n.log("Object(): property is not an object.")
	} else {
		res = n.value.GetObject()
	}
	return NewGameObjectFromId(n.game, res)
}

func (n gameValue) SetObject(obj G.IObject) {
	if n.ptype != meta.ObjectProperty {
		n.log("SetObject(): property is not an object.")
	} else {
		var id ident.Id
		if obj != nil {
			id = obj.Id()
		}
		if e := n.value.SetObject(id); e != nil {
			n.log("SetObject(): error setting value: %v.", e)
		}
	}
}

func (n gameValue) log(format string, v ...interface{}) {
	suffix := fmt.Sprintf(format, v...)
	n.game.Println(n.path, suffix)
}

// this was in SetObject, but that's impossible....
// this could og in the adapter layer, but not here...
// case meta.ObjectProperty | meta.ArrayProperty:
// 	values := strings.Join(n.path, "").GetValues()
// 	if other, ok := object.(GameObject); !ok {
// 		if e := values.ClearValues(); e != nil {
// 			n.log("Set(): error clearing value: %s.",   e)
// 		}
// 	} else {
// 		if e := values.AppendObject(other.gobj.GetId()); e != nil {
// 			n.log("Set(): error appending value: %s.",   e)
// 		}
// 	}
// }