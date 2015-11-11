package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/util/ident"
)

// GameCallback adapts model listeners to game events.
// ( by implementing E.IListen )
type GameCallback struct {
	c         ICreatePlay
	call      CallbackPair
	options   CallbackOptions
	classHint ident.Id
}

type CallbackOptions int

const (
	UseTargetOnly CallbackOptions = 1 << iota
	UseAfterQueue
)

func (opt CallbackOptions) UseTargetOnly() bool {
	return opt&UseTargetOnly != 0
}
func (opt CallbackOptions) UseAfterQueue() bool {
	return opt&UseAfterQueue != 0
}

// HandleEvent implements E.IListen.
func (cb GameCallback) HandleEvent(evt E.IEvent) (err error) {
	if act, okay := evt.Data().(*RuntimeAction); !okay {
		err = fmt.Errorf("unexpected game event data type %T", act)
	} else {
		// callbacks from scripts can ask to be limited to the "target" phase,
		// whereas event listeneres are either registered as part of the bubbling or capturing cycle.
		if trigger := !cb.options.UseTargetOnly() || isTargetPhase(evt); trigger {
			if cb.options.UseAfterQueue() {
				act.runAfterDefaults(cb.call)
			} else {
				play := cb.c.NewPlay(act, cb.classHint)
				cb.call.call(play)
				if act.cancelled {
					evt.StopImmediatePropagation()
					evt.PreventDefault()
				}
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
