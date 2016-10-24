package internal

import (
	"github.com/ionous/mars/core"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// filter listeners to the events appropriate for this target
// MARS: make listener interface into an iterator so we can avoid the copy.
func NewGameListeners(d *Dispatch, evt E.IEvent, target ident.Id, ls meta.Listeners) GameListeners {
	filtered := []meta.Listener{}

	cls, isClassTarget := evt.CurrentTarget().(ClassTarget)
	isTargetPhase := evt.Phase() == E.TargetPhase
	//
	for i := 0; i < ls.NumListener(); i++ {
		l := ls.ListenerNum(i)
		trigger := false
		// MARS: these should probably be two different iterator/functions
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
	return GameListeners{d, filtered}
}

// implements EventListeners
type GameListeners struct {
	d        *Dispatch
	filtered []meta.Listener
}

func (gls GameListeners) NumListener() int {
	return len(gls.filtered)
}

func (gls GameListeners) ListenerNum(i int) E.IListen {
	l := gls.filtered[i]
	return GameListener{gls.d, l}
}

// GameListener adapts model listeners to game event target E.IListener iteration
type GameListener struct {
	dispatch *Dispatch
	listen   meta.Listener
}

// HandleEvent implements E.IListen.
func (gl GameListener) HandleEvent(evt E.IEvent) (err error) {
	call := gl.listen.GetCallback()
	if gl.listen.GetOptions().UseAfterQueue() {
		gl.dispatch.RunActionLater(call, gl.listen.GetClass())
	} else {
		if e := gl.dispatch.RunActionNow(call, gl.listen.GetClass()); e != nil {
			if _, cancelled := e.(core.StopNow); !cancelled {
				err = e
			} else {
				evt.StopImmediatePropagation()
				evt.PreventDefault()
			}
		}
	}
	return err
}
