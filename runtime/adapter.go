package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"strings"
)

//
// GameEventAdapter implements game.Play.
//
type GameEventAdapter struct {
	*Game
	data      *RuntimeAction
	cancelled bool
}

//
func (ga *GameEventAdapter) The(name string) G.IObject {
	return ga.GetObject(name)
}

//
func (ga *GameEventAdapter) Our(name string) G.IObject {
	return ga.GetObject(name)
}

//
func (ga *GameEventAdapter) A(name string) G.IObject {
	return ga.GetObject(name)
}

//
func (ga *GameEventAdapter) Add(data string) (ret G.IObject) {
	if _, ok := ga.Model.Classes.FindClassBySingular(data); !ok {
		ret = NullObject{}
	} else {

	}
	return nil
}

//
func (ga *GameEventAdapter) Remove(G.IObject) {
}

//
func (ga *GameEventAdapter) Any(name string) (obj G.IObject) {
	cls, _ := ga.Model.Classes.FindClassBySingular(name)
	if gobj := ga.FindFirstOf(cls); gobj != nil {
		obj = ObjectAdapter{ga.Game, gobj}
	} else {
		ga.log.Printf("no objects found of type requested `%s`", name)
		obj = NullObject{}
	}
	return obj
}

//
func (ga *GameEventAdapter) Say(texts ...string) {
	if len(texts) > 0 {
		for i, text := range texts {
			if strings.Contains(text, "{{") {
				// FIX? use source text as a cache index?
				// NOTE: cant easily use caller program counter index, because sub-functions might break that.
				if text, e := reallySlow(text, ga.data.values); e != nil {
					ga.log.Println(e)
				} else {
					texts[i] = text
				}
			}
		}
		// join the strings just like print would
		text := strings.Join(texts, " ")
		// find the new lines
		lines := strings.Split(text, "\n")
		ga.output.ScriptSays(lines)
	}
}

//
// func (ga *GameEventAdapter) Report(texts ...string) {
// 	if len(texts) > 0 {
// 		text := strings.Join(texts, " ")
// 		ga.output.Println(text)
// 	}
// }

//
func (ga *GameEventAdapter) Log(texts ...string) {
	if len(texts) > 0 {
		text := strings.Join(texts, " ")
		ga.output.Log(text)
	}
}

//
func (ga *GameEventAdapter) StopHere() {
	ga.cancelled = true
}

//
func (ga *GameEventAdapter) Rules() G.IGameRules {
	return ga.Game
}

//
// could make a map that implements IObject?
// could use special keys for $name, $fullname, $game, etc.
// the point would be, what exactly?
func (ga *GameEventAdapter) GetObject(name string) (obj G.IObject) {
	// asking by original name
	if gobj, ok := ga.FindObject(name); ok {
		obj = ObjectAdapter{ga.Game, gobj}
	}
	// testing against data, because sometimes the adapter isnt invoked via an event
	// use different interfaces perhaps? maybe after injection works?
	if obj == nil && ga.data != nil {
		ctx := map[string]int{"action.Source": 0, "action.Target": 1, "action.Context": 2}
		// asking by action.something
		if index, okay := ctx[name]; okay {
			if gobj := ga.data.objs[index]; gobj != nil {
				obj = ObjectAdapter{ga.Game, gobj}
			}
		}
		// asking by class name, ex. The("story")
		if obj == nil {
			src := ga.data.objs[0]
			cls, _ := ga.Model.Classes.FindClassBySingular(name)
			if src.inst.Class().CompatibleWith(cls.Id()) {
				obj = ObjectAdapter{ga.Game, src}
			}
		}
	}
	// logging and safety
	if obj == nil {
		ga.log.Output(3, fmt.Sprintf("unknown object requested `%s`", name))
		obj = NullObject{}
	}
	return obj
}
