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

func (g *Game) NewAdapter() *internal.GameEventAdapter {
	return internal.NewGameAdapter(g.game)
}

func NullObject(name string) G.IObject {
	return internal.NullObject(name)
}

// ex. after having parsed and matched raw player input.
func (g *Game) QueueActionInstances(act meta.Action, objects []meta.Instance) *internal.RuntimeAction {
	data := internal.NewRuntimeAction(act, objects)
	g.game.QueueAction(data)
	return data
}

// ex. parsing raw input, ending turns, testing.
func (g *Game) QueueAction(action string, nouns ...ident.Id) (ret *internal.RuntimeAction, err error) {
	actionId := internal.MakeStringId(action)
	if act, ok := g.GetAction(actionId); !ok {
		err = fmt.Errorf("couldnt find action %s", action)
	} else if data, e := g.game.NewRuntimeAction(act, nouns...); e != nil {
		err = e
	} else {
		g.game.QueueAction(data)
		ret = data
	}
	return ret, err
}

// Update game until the event queue till empty, or errored.
func (g *Game) ProcessActions() error {
	return g.game.ProcessActions()
}

func SaveGame(g G.Play, autosave bool) (string, error) {
	i := g.(*internal.GameEventAdapter)
	return i.Game.SaveGame(autosave)
}
