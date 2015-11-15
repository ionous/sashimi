package runtime

import (
	"fmt"
	// E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	// M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/runtime/internal"
	//"github.com/ionous/sashimi/runtime/memory"
	"github.com/ionous/sashimi/util/ident"
	//"log"
	//"math/rand"
)

type Game struct {
	api.Model
	game *internal.Game
}

func (g *Game) NewAdapter() G.Play {
	return internal.NewGameAdapter(g.game)
}
func (g *Game) NewGameObject(inst api.Instance) G.IObject {
	return internal.NewGameObject(g.game, inst)
}

func NullObject(name string) G.IObject {
	return internal.NullObject(name)
}

func (g *Game) QueueAction(act api.Action, objects []api.Instance) *internal.RuntimeAction {
	tgt := internal.NewObjectTarget(g.game, objects[0])
	data := internal.NewAction(g.game, act, objects)
	g.game.Queue.QueueEvent(tgt, act.GetEvent().GetId(), data)
	return data
}

// mainly for testing; manual send of an event
func (g *Game) QueueEvent(event string, nouns ...ident.Id,
) (ret *internal.RuntimeAction, err error,
) {
	eventId := internal.MakeStringId(event)
	if event, ok := g.GetEvent(eventId); !ok {
		err = fmt.Errorf("couldnt find event %s", event)
	} else if act, e := g.game.NewRuntimeAction(event.GetAction(), nouns...); e != nil {
		err = e
	} else {
		tgt := internal.NewObjectTarget(g.game, act.GetTarget())
		g.game.Queue.QueueEvent(tgt, event.GetId(), act)
		ret = act
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
