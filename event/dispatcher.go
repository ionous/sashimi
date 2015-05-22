package event

//
// Subscription management for event handlers.
// A single dispatcher is "flat", it doesn't itself handle capturing or bubbling.
//
func NewDispatcher() Dispatcher {
	d := Dispatcher{make(EventListeners), make(EventListeners)}
	return d
}

//
// The default dispatcher: see NewDispatcher()
//
type Dispatcher struct {
	bubble  EventListeners
	capture EventListeners
}

//
func (this *Dispatcher) Listen(evt string, handler IListen, capture bool) {
	this.getMap(capture).AddEventListener(evt, handler)
}

//
func (this *Dispatcher) Silence(evt string, handler IListen, capture bool) {
	this.getMap(capture).RemoveEventListener(evt, handler)
}

//
func (this *Dispatcher) Dispatch(evt IEvent) (err error) {
	phase := evt.Phase()
	// capturing or targeting? trigger capture listeners
	if phase != BubblingPhase {
		err = this.capture.HandleEvents(evt, true)
	}
	// bubbling or targeting? trigger bubble listeners
	if err == nil && phase != CapturingPhase {
		err = this.bubble.HandleEvents(evt, false)
	}
	// FIX: does dispatch really need an error handling? i, personally, am not so sure.
	return err
}

//
func (this *Dispatcher) getMap(capture bool) (ret EventListeners) {
	if capture {
		ret = this.capture
	} else {
		ret = this.bubble
	}
	return ret
}
