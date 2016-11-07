package internal

import (
	"github.com/ionous/mars/core"
	"github.com/ionous/mars/rtm"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// filter listeners to the events appropriate for this target
// MARS: make listener interface into an iterator so we can avoid the copy.
func NewGameListeners(act *rtm.ActionRuntime, evt E.IEvent, target ident.Id, ls meta.Listeners) GameListeners {
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
	return GameListeners{act, filtered}
}

// implements EventListeners
type GameListeners struct {
	act      *rtm.ActionRuntime
	filtered []meta.Listener
}

func (gls GameListeners) NumListener() int {
	return len(gls.filtered)
}

func (gls GameListeners) ListenerNum(i int) E.IListen {
	l := gls.filtered[i]
	return GameListener{gls.act, l}
}

// GameListener adapts model listeners to game event target E.IListener iteration
type GameListener struct {
	act    *rtm.ActionRuntime
	listen meta.Listener
}

// HandleEvent implements E.IListen.
func (gl GameListener) HandleEvent(evt E.IEvent) (err error) {
	call := gl.listen.GetCallback()
	if gl.listen.GetOptions().UseAfterQueue() {
		err = gl.act.RunLater(call, gl.listen.GetClass())
	} else {
		if e := gl.act.RunNow(call, gl.listen.GetClass()); e != nil {
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
