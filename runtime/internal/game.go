package internal

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/runtime/api"
	//	"github.com/ionous/sashimi/runtime/memory"
	"github.com/ionous/sashimi/util/ident"
)

type Game struct {
	ModelApi api.Model
	RuntimeCore
	Queue  EventQueue
	Tables table.Tables
}

func NewGame(m api.Model, core RuntimeCore, tables table.Tables) *Game {
	return &Game{
		m,
		core,
		EventQueue{E.NewQueue()},
		tables,
	}
}

func (g *Game) newPlay(data interface{}, hint ident.Id) G.Play {
	adapter := NewGameAdapter(g)
	adapter.data = data.(*RuntimeAction)
	adapter.hint = hint
	return adapter
}

func (g *Game) Random(n int) int {
	return g.Rand.Intn(n)
}

// class or instance id
func (g *Game) dispatch(evt E.IEvent, target ident.Id) (err error) {
	if src, ok := g.ModelApi.GetEvent(evt.Id()); ok {
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

// TODO: unwind this.
func (g *Game) NewRuntimeAction(action api.Action, nouns ...ident.Id,
) (ret *RuntimeAction, err error,
) {
	types := action.GetNouns()
	switch diff := len(nouns) - len(types); {
	case diff < 0:
		err = fmt.Errorf("too few nouns specified for '%s', %d", action, diff)
	case diff > 0:
		err = fmt.Errorf("too many nouns specified for '%s', +%d", action, diff)
	default:
		objs := make([]api.Instance, len(types))
		for i, class := range types {
			noun := nouns[i]
			if gobj, ok := g.ModelApi.GetInstance(noun); !ok {
				err = M.InstanceNotFound(noun.String())
				break
			} else if !g.ModelApi.AreCompatible(gobj.GetParentClass().GetId(), class) {
				err = TypeMismatch(noun.String(), class.String())
				break
			} else {
				objs[i] = gobj
			}
		}
		if err == nil {
			ret = &RuntimeAction{game: g, action: action, objs: objs}
		}
	}
	return ret, err
}

// TODO: add interfaces for start and end
func (g *Game) SendMessage(tgt E.ITarget, msg *E.Message) (err error) {
	defer g.Frame(tgt, msg)()
	path := E.NewPathTo(tgt)

	// game.log.Printf("sending `%s` to: %s", msg.Name, path)
	if runDefault, e := msg.Send(path); e != nil {
		err = e
	} else {
		if runDefault {
			if act, ok := msg.Data.(*RuntimeAction); !ok {
				err = fmt.Errorf("unknown action data %T", msg.Data)
			} else {
				err = act.runDefaultActions()
			}
		}
	}
	return err
}
