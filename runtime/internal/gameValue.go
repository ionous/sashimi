package internal

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// FIX? because we store property path, we cant test for equality directly.
// ( therefore changed to pointer value for gopherjs optimization )
type gameValue struct {
	game  *GameEventAdapter
	path  PropertyPath
	ptype meta.PropertyType
	value meta.Value
}

func (n *gameValue) Num() (ret float64) {
	if n.ptype != meta.NumProperty {
		n.log("Num(): property is not a number.")
	} else {
		ret = n.value.GetNum()
	}
	return
}

func (n *gameValue) SetNum(value float64) {
	if n.ptype != meta.NumProperty {
		n.log("SetNum(): property is not a number.")
	} else if e := n.value.SetNum(value); e != nil {
		n.log("SetNum(): error setting value: %v.", e)
	}
}

func (n *gameValue) Text() (ret string) {
	if n.ptype != meta.TextProperty {
		n.log("Text(): property is not text.")
	} else {
		ret = n.value.GetText()
	}
	return
}

func (n *gameValue) SetText(text string) {
	if n.ptype != meta.TextProperty {
		n.log("SetText(): property is not text.")
	} else if e := n.value.SetText(text); e != nil {
		n.log("SetText(): error setting value: %v.", e)
	}
}

// TBD: should these be logged? its sure nice to have be able to test objects generically for properties
func (n *gameValue) Object() G.IObject {
	var res ident.Id
	if n.ptype != meta.ObjectProperty {
		n.log("Object(): property is not an object.")
	} else {
		res = n.value.GetObject()
	}
	return n.game.NewGameObjectFromId(res)
}

func (n *gameValue) SetObject(obj G.IObject) {
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

func (n *gameValue) State() (ret ident.Id) {
	if n.ptype != meta.StateProperty {
		n.log("State(): property is not a state.")
	} else {
		ret = n.value.GetState()
	}
	return
}

func (n *gameValue) SetState(val ident.Id) {
	if n.ptype != meta.StateProperty {
		n.log("SetState(): property is not a state.")
	} else if e := n.value.SetState(val); e != nil {
		n.log("SetState(): error setting value: %v.", e)
	}
}

func (n *gameValue) log(format string, v ...interface{}) {
	suffix := fmt.Sprintf(format, v...)
	n.game.Println(n.path, suffix)
}
