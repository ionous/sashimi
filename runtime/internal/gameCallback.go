package internal

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

// filter listeners to the events appropriate for this target
func NewGameListeners(game *Game, evt E.IEvent, target ident.Id, ls meta.Listeners) GameListeners {
	filtered := []meta.Listener{}

	cls, isClassTarget := evt.CurrentTarget().(ClassTarget)
	isTargetPhase := evt.Phase() == E.TargetPhase
	//
	for i := 0; i < ls.NumListener(); i++ {
		l := ls.ListenerNum(i)
		trigger := false
		if isClassTarget {
			if l.GetInstance().Empty() && target == l.GetClass() {
				// expands the target phase to include the instance's class.
				isTargetPhase := isTargetPhase || cls.from == evt.Target()
				trigger = l.GetOptions().UseTargetOnly() == isTargetPhase
			}
		} else {
			if target == l.GetInstance() {
				trigger = l.GetOptions().UseTargetOnly() == isTargetPhase
			}
		}
		if trigger {
			filtered = append(filtered, l)
		}
	}
	return GameListeners{game, filtered}
}

// implements EventListeners
type GameListeners struct {
	game     *Game
	filtered []meta.Listener
}

func (gl GameListeners) NumListener() int {
	return len(gl.filtered)
}

func (gl GameListeners) ListenerNum(i int) E.IListen {
	l := gl.filtered[i]
	return GameCallback{gl.game, l}
}

// GameCallback adapts model listeners to game events.
// ( by implementing E.IListen )
type GameCallback struct {
	game *Game // FIX: switch to context?
	meta.Listener
}

// HandleEvent implements E.IListen.
func (cb GameCallback) HandleEvent(evt E.IEvent) (err error) {
	if act, ok := evt.Data().(*RuntimeAction); !ok {
		err = errutil.New("unexpected game event data type", sbuf.Type{act})
	} else {
		call := cb.GetCallback()
		if cb.GetOptions().UseAfterQueue() {
			act.runAfterDefaults(call) // FIX: switch to adding to adapter? i just dont like that the action changes...
		} else {
			rt := NewMars(cb.game, cb.game.newPlay(act, cb.GetClass()))
			if e := rt.Execute(call); e != nil {
				err = e
			} else if act.cancelled {
				evt.StopImmediatePropagation()
				evt.PreventDefault()
			}
		}
	}
	return err
}
