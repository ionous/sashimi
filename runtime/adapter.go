package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/lang"
	"strings"
)

// GameEventAdapter implements game.Play.
type GameEventAdapter struct {
	*Game
	data *RuntimeAction
	hint ident.Id
}

//
// NewGameAdapter gives the passed game the IPlay interface
// Public for testing.
//
func NewGameAdapter(game *Game) *GameEventAdapter {
	return &GameEventAdapter{Game: game}
}

// NewObjectAdapter gives the passed game object the IObject interface.
// Public for testing.
func NewObjectAdapterFromId(game *Game, id ident.Id) (ret G.IObject) {
	if inst, ok := game.ModelApi.GetInstance(id); ok {
		ret = ObjectAdapter{game, inst}
	} else {
		ret = NullObjectSource("", 2)
	}
	return ret
}

func NewObjectAdapter(game *Game, inst api.Instance) (ret G.IObject) {
	if inst != nil {
		ret = ObjectAdapter{game, inst}
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
	id := MakeStringId(name)
	if val, ok := ga.Globals[id]; !ok {
		ga.Log("no such global", name)
	} else {
		ret = val.Interface()
		okay = true
	}
	return ret, okay
}

// Add creates a new (data) object.
func (ga *GameEventAdapter) NewFrom(class string) G.IObject {
	panic("not supported")
	// clsid := StripStringId(class)
	// if cls, ok := ga.ModelApi.GetClass(clsid); !ok {
	// 	ret = NullObjectSource(class, 2)
	// } else {
	// 	id := ident.MakeUniqueId()
	// 	if gobj, e := NewGameObject(ga.ModelApi, id, cls, cls, ga.Game.Tables); e != nil {
	// 		ga.Log(e)
	// 		ret = NullObjectSource(class, 2)
	// 	} else {
	// 		ret = NewObjectAdapter(ga.Game, gobj)
	// 		ga.Objects[id] = gobj
	// 	}
	// }
	// return ret
}

//
func (ga *GameEventAdapter) Visit(class string, visits func(G.IObject) bool) (okay bool) {
	clsid := StripStringId(class)
	for i := 0; i < ga.ModelApi.NumInstance(); i++ {
		gobj := ga.ModelApi.InstanceNum(i)
		if id := gobj.GetParentClass().GetId(); ga.ModelApi.AreCompatible(id, clsid) {
			if visits(NewObjectAdapter(ga.Game, gobj)) {
				okay = true
				break
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
	ga.data.cancelled = true
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
	if gobj, ok := ga.ModelApi.GetInstance(StripStringId(name)); ok {
		ret, okay = NewObjectAdapter(ga.Game, gobj), true
	} else if ga.data != nil {
		// testing against ga.data b/c sometimes the adapter isnt invoked via an event.
		// to fix use different interfaces perhaps?
		if obj, ok := ga.data.findByParamName(name); ok {
			ret, okay = obj, true
		} else {
			clsid := MakeStringId(ga.ModelApi.Pluralize(lang.StripArticle(name)))
			if clsid == ga.hint {
				ret, okay = ga.data.getObject(0)
			} else {
				ret, okay = ga.data.findByClass(clsid)
			}
			if !okay {
				ga.Log("couldnt find object", name, "including class", clsid)
				fmt.Println("!!!", clsid)
			}
		}
	}
	return ret, okay
}
