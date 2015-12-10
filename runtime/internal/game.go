package internal

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type Game struct {
	Model meta.Model
	RuntimeCore
	Queue ActionQueue
}

func NewGame(core RuntimeCore, m meta.Model) *Game {
	return &Game{
		m,
		core,
		NewActionQueue(),
	}
}

func (g *Game) newPlay(data *RuntimeAction, hint ident.Id) G.Play {
	return &GameEventAdapter{Game: g, data: data, hint: hint}
}

func (g *Game) Random(n int) int {
	return g.Rand.Intn(n)
}

// target: class or instance id
// note: we get dispatch multiple times for each event:
// on the capture, target, and bubble cycles.
func (g *Game) dispatch(evt E.IEvent, target ident.Id) (err error) {
	if src, ok := g.Model.GetEvent(evt.Id()); ok {
		if ls, ok := src.GetListeners(true); ok {
			err = E.Capture(evt, NewGameListeners(g, evt, target, ls))
		}
		if err == nil {
			if ls, ok := src.GetListeners(false); ok {
				err = E.Bubble(evt, NewGameListeners(g, evt, target, ls))
			}
		}
	}
	return
}

func (g *Game) QueueAction(data *RuntimeAction) {
	future := &QueuedAction{data: data}
	g.Queue.QueueFuture(future)
}

func (g *Game) ProcessActions() error {
	return g.Queue.ProcessActions(g)
}

// NewRuntimeAction: captures an action and bunch of nouns.
// TODO: unwind this... possibly now merege this with the go bits....
func (g *Game) NewRuntimeAction(action meta.Action, nouns ...ident.Id,
) (ret *RuntimeAction, err error,
) {
	types := action.GetNouns()
	switch diff := len(nouns) - len(types); {
	case diff < 0:
		err = fmt.Errorf("too few nouns specified for '%s', %d", action, diff)
	case diff > 0:
		err = fmt.Errorf("too many nouns specified for '%s', +%d", action, diff)
	default:
		objs := make([]meta.Instance, len(types))
		for i, class := range types {
			noun := nouns[i]
			if gobj, ok := g.Model.GetInstance(noun); !ok {
				err = InstanceNotFound(noun.String())
				break
			} else if !g.Model.AreCompatible(gobj.GetParentClass().GetId(), class) {
				err = TypeMismatch(noun.String(), class.String())
				break
			} else {
				objs[i] = gobj
			}
		}
		if err == nil {
			ret = NewRuntimeAction(action, objs)
		}
	}
	return ret, err
}
