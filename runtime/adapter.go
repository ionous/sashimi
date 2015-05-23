package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"strings"
)

//
// implements Game.Play
// used internally, mainly for the event return values
//
type GameEventAdapter struct {
	*Game
	data      *RuntimeAction
	cancelled bool
}

//
//(I)Play interface
//
func (this *GameEventAdapter) The(name string) G.IObject {
	return this.GetObject(name)
}

//
func (this *GameEventAdapter) Our(name string) G.IObject {
	return this.GetObject(name)
}

//
func (this *GameEventAdapter) A(name string) G.IObject {
	return this.GetObject(name)
}

//
func (this *GameEventAdapter) Any(name string) (obj G.IObject) {
	cls, _ := this.Model.Classes.FindClassBySingular(name)
	gobj := this.FindFirstOf(cls)
	if gobj != nil {
		obj = ObjectAdapter{this.Game, gobj}
	} else {
		this.log.Printf("no objects found of type requested `%s`", name)
		obj = this.nullobj
	}
	return obj
}

//
func (this *GameEventAdapter) Say(texts ...string) {
	this._print(texts...)
}

//
func (this *GameEventAdapter) Report(texts ...string) {
	if len(texts) > 0 {
		text := strings.Join(texts, " ")
		this.console.Println(text)
	}
}

//
func (this *GameEventAdapter) Log(texts ...string) {
	if len(texts) > 0 {
		text := strings.Join(texts, " ")
		this.console.Println(text)
	}
}

//
func (this *GameEventAdapter) StopHere() {
	this.cancelled = true
}

//
func (this *GameEventAdapter) Rules() G.IGameRules {
	return this.Game
}

func (this *GameEventAdapter) _print(texts ...string) {
	if len(texts) > 0 {
		for i, text := range texts {
			if strings.Contains(text, "{{") {
				// FIX? use source text as a cache index?
				// NOTE: cant easily use caller program counter index, because sub-functions might break that.
				if text, e := reallySlow(text, this.data.values); e != nil {
					this.log.Println(e)
				} else {
					texts[i] = text
				}
			}
		}
		// join the strings just like print would
		text := strings.Join(texts, " ")
		// find the new lines
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			this.console.Println(line)
		}
	}
}

//
// could make a map that implements IObject?
// could use special keys for $name, $fullname, $game, etc.
// the point would be, what exactly?
func (this *GameEventAdapter) GetObject(name string) (obj G.IObject) {
	// asking by original name
	if obj == nil {
		if gobj, ok := this.FindObject(name); ok {
			obj = ObjectAdapter{this.Game, gobj}
		}
	}
	// testing against data, because sometimes the adapter isnt invoked via an event
	// use different interfaces perhaps? maybe after injection works?
	if obj == nil && this.data != nil {
		ctx := map[string]int{"action.Source": 0, "action.Target": 1, "action.Context": 2}
		// asking by action.something
		if index, okay := ctx[name]; okay {
			if gobj := this.data.objs[index]; gobj != nil {
				obj = ObjectAdapter{this.Game, gobj}
			}
		}
		// asking by class name, ex. The("story")
		if obj == nil {
			src := this.data.objs[0]
			cls, _ := this.Model.Classes.FindClassBySingular(name)
			if src.info.CompatibleWith(cls) {
				obj = ObjectAdapter{this.Game, src}
			}
		}
	}
	// logging and safety
	if obj == nil {
		this.log.Output(3, fmt.Sprintf("unknown object requested `%s`", name))
		obj = this.nullobj
	}
	return obj
}
