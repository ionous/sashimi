package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	M "github.com/ionous/sashimi/model"
)

//
// Adapt model listeners to game events.
// ( by implementing E.IListen )
//
type GameCallback struct {
	game *Game
	*M.ListenerCallback
	call CallbackPair
}

//
// Implementation of E.IListen.
//
func (cb GameCallback) HandleEvent(evt E.IEvent) (err error) {
	if act, okay := evt.Data().(*RuntimeAction); !okay {
		err = fmt.Errorf("unexpected game event data type %T", act)
	} else {
		// callbacks from scripts can ask to be limited to the "target" phase,
		// whereas event listeneres are either registered as part of the bubbling or capturing cycle.
		if trigger := !cb.UseTargetOnly() || isTargetPhase(evt); trigger {
			if cb.UseAfterQueue() {
				act.runAfterDefaults(cb.call)
			} else {
				if !act.runCallback(cb.call, cb.Class) {
					evt.StopImmediatePropagation()
					evt.PreventDefault()
				}
			}
		}
	}
	return err
}

//
// Expands the target phase to include the target instance and the instance's class.
//
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
