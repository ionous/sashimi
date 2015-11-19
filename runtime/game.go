package runtime

import (
	"fmt"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/runtime/internal"
	"github.com/ionous/sashimi/util/ident"
)

type Game struct {
	meta.Model
	game *internal.Game
}

func (g *Game) NewAdapter() G.Play {
	return internal.NewGameAdapter(g.game)
}
func (g *Game) NewGameObject(inst meta.Instance) G.IObject {
	return internal.NewGameObject(g.game, inst)
}

func NullObject(name string) G.IObject {
	return internal.NullObject(name)
}

func (g *Game) QueueActionInstances(act meta.Action, objects []meta.Instance) *internal.RuntimeAction {
	tgt := internal.NewObjectTarget(g.game, objects[0])
	data := internal.NewRuntimeAction(g.game, act, objects)
	g.game.Queue.QueueEvent(tgt, act.GetEvent().GetId(), data)
	return data
}

func (g *Game) QueueAction(action string, nouns ...ident.Id) (ret *internal.RuntimeAction, err error) {
	actionId := internal.MakeStringId(action)
	if act, ok := g.GetAction(actionId); !ok {
		err = fmt.Errorf("couldnt find action %s", action)
	} else if data, e := g.game.NewRuntimeAction(act, nouns...); e != nil {
		err = e
	} else {
		tgt := internal.NewObjectTarget(g.game, data.GetTarget())
		g.game.Queue.QueueEvent(tgt, act.GetEvent().GetId(), data)
		ret = data
	}
	return ret, err
}

// ProcessEvents in the event queue till empty, or errored.
func (g *Game) ProcessEvents() (err error) {
	for !g.game.Queue.Empty() {
		tgt, msg := g.game.Queue.Pop()
		if e := g.game.SendMessage(tgt, msg); e != nil {
			g.game.Println("error", e)
			err = e
			break
		}
	}
	return err
}
