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

// NewGameObject gives the passed game object the IObject interface.
// Public for testing.
func NewGameObjectFromId(game *Game, id ident.Id) (ret G.IObject) {
	var inst api.Instance
	if i, ok := game.ModelApi.GetInstance(id); ok {
		inst = i
	}
	return NewGameObject(game, inst)
}

func NewGameObject(game *Game, inst api.Instance) (ret G.IObject) {
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
	instances := []api.Instance{}
	clsid := StripStringId(class)
	for i := 0; i < ga.ModelApi.NumInstance(); i++ {
		gobj := ga.ModelApi.InstanceNum(i)
		if id := gobj.GetParentClass().GetId(); ga.ModelApi.AreCompatible(id, clsid) {
			instances = append(instances, gobj)
		}
	}
	return iList{ga.Game, NewPath(clsid), instances}
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

var DebugGet = false

// could make a map that implements IObject?
// could use special keys for $name, $fullname, $game, etc.
// FUTURE: use dependency injection instead
func (ga *GameEventAdapter) GetObject(name string) (ret G.IObject) {
	id := StripStringId(name)
	if gobj, ok := ga.ModelApi.GetInstance(id); ok {
		ret = NewGameObject(ga.Game, gobj)
	} else if ga.data != nil {
		// testing against ga.data b/c sometimes the adapter isnt invoked via an event.
		// to fix use different interfaces perhaps?
		if obj, ok := ga.data.findByParamName(name); ok {
			ret = obj
		} else {
			found := false
			clsid := MakeStringId(ga.ModelApi.Pluralize(lang.StripArticle(name)))
			if clsid == ga.hint {
				ret, found = ga.data.getObject(0)
			} else {
				ret, found = ga.data.findByClass(clsid)
			}
			if !found {
				ga.Log("couldnt find object", name, "including class", clsid)
				ret = NullObjectSource(RawPath(name), 3)
			}
		}
	}
	DebugGet = false
	return
}
