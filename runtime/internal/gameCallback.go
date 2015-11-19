package internal

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// filter listeners to the events appropriate for this target
func NewGameListeners(game *Game, evt E.IEvent, target ident.Id, ls meta.Listeners) GameListeners {
	filtered := []meta.Listener{}
	for i := 0; i < ls.NumListener(); i++ {
		l := ls.ListenerNum(i)
		if target == l.GetInstance() || target == l.GetClass() {
			// callbacks from scripts can ask to be limited to the "target" phase,
			if trigger := !l.GetOptions().UseTargetOnly() || isTargetPhase(evt); trigger {
				filtered = append(filtered, l)
			}
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
	game *Game
	meta.Listener
}

// HandleEvent implements E.IListen.
func (cb GameCallback) HandleEvent(evt E.IEvent) (err error) {
	if act, ok := evt.Data().(*RuntimeAction); !ok {
		err = fmt.Errorf("unexpected game event data type %T", act)
	} else if fn, ok := cb.game.LookupCallback(cb.GetCallback()); !ok {
		err = fmt.Errorf("couldn't find callback for listener %s", cb.Listener)
	} else {
		if cb.GetOptions().UseAfterQueue() {
			call := QueuedCallback{cb.GetCallback(), fn}
			act.runAfterDefaults(call)
		} else {
			play := cb.game.newPlay(act, cb.GetClass())
			fn(play)
			if act.cancelled {
				evt.StopImmediatePropagation()
				evt.PreventDefault()
			}
		}
	}
	return err
}

// isTargetPhase expands the target phase to include the target instance and the instance's class.
func isTargetPhase(evt E.IEvent) bool {
	targetPhase := evt.Phase() == E.TargetPhase
	if !targetPhase {
		if clsTarget, ok := evt.CurrentTarget().(ClassTarget); ok {
			if clsTarget.host == evt.Target() {
				targetPhase = true
			}
		}
	}
	return targetPhase
}
