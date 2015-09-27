package event

import "fmt"

//
// NewDispatcher creates a single "flat" set of event callbacks; it doesn't handle event hierarchy.
//
func NewDispatcher() Dispatcher {
	return Dispatcher{make(EventListeners), make(EventListeners)}
}

//
// Dispatcher provides a default event callback system.
//
type Dispatcher struct {
	bubble  EventListeners
	capture EventListeners
}

// Listen requests that the passed handler get called for the named event on the capture (or bubble) dispatch cycle,
func (d *Dispatcher) Listen(evt string, handler IListen, capture bool) {
	d.getMap(capture).AddEventListener(evt, handler)
}

// Silence cancels a previous listen request.
func (d *Dispatcher) Silence(evt string, handler IListen, capture bool) {
	d.getMap(capture).RemoveEventListener(evt, handler)
}

// Dispatch triggers capturing or bubbling handlers, depending on the phase of the passed event.
// The dispatcher can only handle events of "Proc" type.
func (d *Dispatcher) Dispatch(evt IEvent) (err error) {
	if proc, ok := evt.(*Proc); !ok {
		err = fmt.Errorf("unknown event type %T", evt)
	} else {
		phase := evt.Phase()
		// capturing or targeting? trigger capture listeners
		if phase != BubblingPhase && !proc.stopMore {
			err = d.capture.HandleEvents(proc, true)
		}
		// bubbling or targeting? trigger bubble listeners
		if err == nil && phase != CapturingPhase && !proc.stopMore {
			err = d.bubble.HandleEvents(proc, false)
		}
	}
	// FIX: does dispatch really need an error handling? i, personally, am not so sure.
	return err
}

//
func (d *Dispatcher) getMap(capture bool) (ret EventListeners) {
	if capture {
		ret = d.capture
	} else {
		ret = d.bubble
	}
	return ret
}
