package runtime

import (
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
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
	hint      *M.ClassInfo
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
func NewObjectAdapter(game *Game, gobj *GameObject) (ret G.IObject) {
	if gobj != nil {
		ret = ObjectAdapter{game, gobj}
	} else {
		ret = NullObjectSource("", 2)
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

// g.Go( Give(player).ToHe(2) )
func (ga *GameEventAdapter) Go(phrase G.RuntimePhrase) {
	phrase.Execute(ga)
}

//
func (ga *GameEventAdapter) Global(name string) interface{} {
	id := M.MakeStringId(name)
	val := ga.Globals[id]
	return val.Interface()
}

//
func (ga *GameEventAdapter) Add(data string) (ret G.IObject) {
	if cls, ok := ga.Model.Classes.FindClassBySingular(data); !ok {
		ret = NullObjectSource(data, 2)
	} else {
		id := ident.MakeUniqueId()
		gobj := &GameObject{id, cls, make(TemplateValues), make(TemplatePool), ga.Game.Model.Tables}
		for _, prop := range cls.Properties() {
			if e := gobj.setValue(prop, prop.Zero(cls.Constraints())); e != nil {
				panic(e)
			}
		}
		ret = NewObjectAdapter(ga.Game, gobj)
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
				if visits(NewObjectAdapter(ga.Game, gobj)) {
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
// FUTURE: use dependency injection instead
func (ga *GameEventAdapter) GetObject(name string) (obj G.IObject) {
	// asking by original name
	if gobj, ok := ga.FindObject(name); ok {
		obj = NewObjectAdapter(ga.Game, gobj)
	}
	// testing against data, because sometimes the adapter isnt invoked via an event
	// use different interfaces perhaps? maybe after injection works?
	if obj == nil && ga.data != nil {
		ctx := map[string]int{"action.Source": 0, "action.Target": 1, "action.Context": 2}
		// asking by action.something
		if index, okay := ctx[name]; okay {
			if gobj := ga.data.objs[index]; gobj != nil {
				obj = NewObjectAdapter(ga.Game, gobj)
			}
		}

		// asking by the class name,ex. The("story")
		if obj == nil {
			if namedCls, ok := ga.Model.Classes.FindClassBySingular(name); ok {
				index := -1
				// of the handler, or action parameter
				if namedCls == ga.hint {
					index = 0
				} else {
					// walk classes of the actions
					for i, srcCls := range ga.data.action.NounSlice() {
						if namedCls == srcCls {
							index = i
							break
						}
					}
				}
				if index >= 0 && index < len(ga.data.objs) {
					obj = NewObjectAdapter(ga.Game, ga.data.objs[index])
				} else {
					//backwards compat
					if obj == nil {
						src := ga.data.objs[0]
						if src.Class().CompatibleWith(namedCls.Id()) {
							obj = NewObjectAdapter(ga.Game, src)
						}
					}
				}
			}
		}
	}
	// logging and safety
	if obj == nil {
		obj = NullObjectSource(name, 3)
	}
	return obj
}
