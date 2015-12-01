package internal

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
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

// NewGameObject gives the passed game object the IObject interface.
// Public for testing.
func NewGameObjectFromId(game *Game, id ident.Id) (ret G.IObject) {
	var inst meta.Instance
	if i, ok := game.Model.GetInstance(id); ok {
		inst = i
	}
	return NewGameObject(game, inst)
}

func NewGameObject(game *Game, inst meta.Instance) (ret G.IObject) {
	if inst != nil {
		ret = GameObject{game, inst}
	} else {
		ret = NullObjectSource(PropertyPath{}, 2)
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

func (ga *GameEventAdapter) List(class string) (ret G.IList) {
	instances := []meta.Instance{}
	clsid := StripStringId(class)
	if _, found := ga.Model.GetClass(clsid); !found {
		ga.Game.Println("List: couldnt find class", clsid)
	} else {
		for i := 0; i < ga.Model.NumInstance(); i++ {
			gobj := ga.Model.InstanceNum(i)
			if id := gobj.GetParentClass().GetId(); ga.Model.AreCompatible(id, clsid) {
				instances = append(instances, gobj)
			}
		}
	}
	return iList{ga.Game, NewPath(clsid), instances}
}

//
func (ga *GameEventAdapter) Say(texts ...string) {
	if len(texts) > 0 {
		text := strings.Join(texts, " ")
		lines := strings.Split(text, "\n")
		ga.Output.ScriptSays(lines)
	}
}

//
func (ga *GameEventAdapter) Log(texts ...interface{}) {
	if len(texts) > 0 {
		text := fmt.Sprintln(texts...)
		ga.Output.Log(text)
	}
}

//
func (ga *GameEventAdapter) StopHere() {
	ga.data.cancelled = true
}

var DebugGet = false

// could make a map that implements IObject?
// could use special keys for $name, $fullname, $game, etc.
// FUTURE: use dependency injection instead
func (ga *GameEventAdapter) GetObject(name string) (ret G.IObject) {
	id := StripStringId(name)
	// empty names are possible from empty strings like Get("")
	if !id.Empty() {
		if gobj, ok := ga.Model.GetInstance(id); ok {
			ret = NewGameObject(ga.Game, gobj)
		} else if ga.data != nil {
			// testing against ga.data b/c sometimes the adapter isnt invoked via an event.
			// to fix use different interfaces perhaps?
			if obj, ok := ga.data.findByParamName(name); ok {
				ret = obj
			} else {
				clsid := MakeStringId(ga.Model.Pluralize(lang.StripArticle(name)))
				found := false
				if clsid == ga.hint {
					ret, found = ga.data.getObject(0)
				} else {
					ret, found = ga.data.findByClass(clsid)
				}
				if !found {
					ga.Log(fmt.Printf("couldnt find object '%s' name including class '%s'", name, clsid))
				}
			}
		}
	}
	if ret == nil {
		ret = NullObjectSource(RawPath(name), 3)
	}
	DebugGet = false
	return
}
