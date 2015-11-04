package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
	"strings"
)

// GameEventAdapter implements game.Play.
type GameEventAdapter struct {
	*Game
	data      *RuntimeAction
	cancelled bool
	hint      ident.Id
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
func (ga *GameEventAdapter) Global(name string) (ret interface{}, okay bool) {
	id := M.MakeStringId(name)
	if val, ok := ga.Globals[id]; !ok {
		ga.Log("no such global", name)
	} else {
		ret = val.Interface()
		okay = true
	}
	return ret, okay
}

//
func (ga *GameEventAdapter) Add(data string) (ret G.IObject) {
	if cls, ok := ga.Model.Classes.FindClassBySingular(data); !ok {
		ret = NullObjectSource(data, 2)
	} else {
		id := ident.MakeUniqueId()
		gobj := &GameObject{id, cls, make(RuntimeValues), ga.Game.Tables}
		for _, prop := range cls.Properties {
			if e := gobj.setValue(prop, prop.GetZero(cls.Constraints)); e != nil {
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
			if gobj.Class().CompatibleWith(cls.Id) {
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
		text := strings.Join(texts, " ")
		lines := strings.Split(text, "\n")
		ga.output.ScriptSays(lines)
	}
}

//
func (ga *GameEventAdapter) Log(texts ...interface{}) {
	if len(texts) > 0 {
		text := fmt.Sprintln(texts...)
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

var DebugGet = false

// could make a map that implements IObject?
// could use special keys for $name, $fullname, $game, etc.
// FUTURE: use dependency injection instead
func (ga *GameEventAdapter) GetObject(name string) (ret G.IObject) {
	if obj, ok := ga.getObject(name); ok {
		ret = obj
	} else {
		ret = NullObjectSource(name, 3)
	}
	DebugGet = false
	return ret
}

func (ga *GameEventAdapter) getObject(name string) (ret G.IObject, okay bool) {
	// asking by object name
	if gobj, ok := ga.Game.FindObject(name); ok {
		ret, okay = NewObjectAdapter(ga.Game, gobj), true
	} else if ga.data != nil {
		// testing against ga.data b/c sometimes the adapter isnt invoked via an event.
		// to fix use different interfaces perhaps?
		if obj, ok := ga.data.findByParamName(name); ok {
			ret, okay = obj, true
		} else if cls, ok := ga.Model.Classes.FindClassBySingular(name); ok {
			// FIX?hint tries to find targeted classln, not sure if its working
			if DebugGet {
				ga.Log("DebugGet: cls:", cls, "hint:", ga.hint)
			}
			if cls.Id == ga.hint {
				ret, okay = ga.data.getObject(0)
			} else {
				ret, okay = ga.data.findByClass(cls)
			}
		}
	}

	return ret, okay
}
