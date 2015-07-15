package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/util/ident"
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
// NewGameAdapter gives the passed game the IPlay interface
// Public for testing.
//
func NewGameAdapter(game *Game) *GameEventAdapter {
	return &GameEventAdapter{Game: game}
}

//
// NewObjectAdapter gives the passed game object the IObject interface.
// Public for testing.
//
func (ga *GameEventAdapter) NewObjectAdapter(gobj *GameObject) (ret G.IObject) {
	if gobj != nil {
		ret = ObjectAdapter{ga, gobj}
	} else {
		ret = NullObject()
	}
	return ret
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
	if cls, ok := ga.Model.Classes.FindClassBySingular(data); !ok {
		ret = NullObject()
	} else {
		id := ident.MakeUniqueId()
		gobj := &GameObject{id, cls, make(TemplateValues), make(TemplatePool), ga.Game.Model.Tables}
		for _, prop := range cls.Properties() {
			if e := gobj.setValue(prop, prop.Zero(cls.Constraints())); e != nil {
				panic(e)
			}
		}
		ret = ga.NewObjectAdapter(gobj)
		ga.Objects[id] = gobj
	}
	return ret
}

//
func (ga *GameEventAdapter) Visit(class string, visits func(G.IObject) bool) (okay bool) {
	if cls, ok := ga.Model.Classes.FindClass(class); !ok {
		ga.log.Printf("no objects found of type requested `%s`", class)
	} else {
		for _, gobj := range ga.Objects {
			if gobj.Class().CompatibleWith(cls.Id()) {
				if visits(ga.NewObjectAdapter(gobj)) {
					okay = true
					break
				}
			}
		}
	}
	return okay
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
		obj = ga.NewObjectAdapter(gobj)
	}
	// testing against data, because sometimes the adapter isnt invoked via an event
	// use different interfaces perhaps? maybe after injection works?
	if obj == nil && ga.data != nil {
		ctx := map[string]int{"action.Source": 0, "action.Target": 1, "action.Context": 2}
		// asking by action.something
		if index, okay := ctx[name]; okay {
			if gobj := ga.data.objs[index]; gobj != nil {
				obj = ga.NewObjectAdapter(gobj)
			}
		}
		// asking by class name, ex. The("story")
		if obj == nil {
			for _, src := range ga.data.objs {
				cls, _ := ga.Model.Classes.FindClassBySingular(name)
				if src.Class().CompatibleWith(cls.Id()) {
					obj = ga.NewObjectAdapter(src)
					if src.Class() == cls {
						break // best match
					}
				}
			}
		}
	}
	// logging and safety
	if obj == nil {
		ga.log.Output(3, fmt.Sprintf("unknown object requested `%s`", name))
		obj = NullObject()
	}
	return obj
}
