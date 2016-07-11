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
	// when we handle events from callbacks, we set this to the target's class to help resolve names specified by user code.
	hint ident.Id
}

// NewGameAdapter gives the passed game the IPlay interface
// Public for testing.
func NewGameAdapter(game *Game) *GameEventAdapter {
	return &GameEventAdapter{Game: game}
}

// NewGameObject gives the passed game object the IObject interface.
// Public for testing.
func (ga *GameEventAdapter) NewGameObjectFromId(id ident.Id) (ret G.IObject) {
	var inst meta.Instance
	if i, ok := ga.Model.GetInstance(id); ok {
		inst = i
	}
	return ga.NewGameObject(inst)
}

func (ga *GameEventAdapter) NewGameObject(inst meta.Instance) (ret G.IObject) {
	if inst != nil {
		ret = GameObject{ga, inst}
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

// g.Go( Move("the player").To("the pit of disrepair") )
func (ga *GameEventAdapter) Go(phrase G.RuntimePhrase, phrases ...G.RuntimePhrase) (ret G.IPromise) {
	if len(phrases) == 0 {
		future := &QueuedPhrase{data: ga.data, run: phrase}
		ga.Queue.QueueFuture(future)
		ret = NewPendingChain(ga, future)
	} else {
		phrases := append([]G.RuntimePhrase{phrase}, phrases...)
		future := &QueuedPhrases{data: ga.data, run: phrases}
		ga.Queue.QueueFuture(future)
		ret = NewPendingChain(ga, future)
	}
	return
}

func (ga *GameEventAdapter) List(class string) (ret G.IList) {
	instances := []meta.Instance{}
	clsid := StripStringId(class)
	if _, found := ga.Model.GetClass(clsid); !found {
		ga.Game.Println("List: couldnt find class", clsid)
	} else {
		for i := 0; i < ga.Model.NumInstance(); i++ {
			gobj := ga.Model.InstanceNum(i)
			if id := gobj.GetParentClass(); ga.Model.AreCompatible(id, clsid) {
				instances = append(instances, gobj)
			}
		}
	}
	return iList{ga, NewPath(clsid), instances}
}

//
func (ga *GameEventAdapter) Say(texts ...string) {
	if len(texts) > 0 {
		text := strings.Join(texts, " ")
		lines := strings.Split(text, lang.NewLine)
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
			ret = ga.NewGameObject(gobj)
		} else if ga.data != nil {
			// testing against ga.data b/c sometimes the adapter isnt invoked via an event.
			// to fix use different interfaces perhaps?
			if obj, ok := ga.data.findByName(ga.Model, name, ga.hint); ok {
				ret = ga.NewGameObject(obj)
			} else {
				msg := fmt.Sprintf("couldnt find object named '%s(%s)'", name, id)
				ga.Log(msg)
				//panic(msg)
			}
		}
	}
	if ret == nil {
		ret = NullObjectSource(RawPath(name), 3)
	}
	DebugGet = false
	return
}
