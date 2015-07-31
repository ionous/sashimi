package event

//
// EvenListeners contains a mapping of event name to a list of event listeners -- used internally by the default dispatcher.
//
type EventListeners map[string]EventList

//
// EvenList contains a list of IListen callbacks
//
type EventList []IListen

//
// AddEventListener to register a new function callback for the passed event name.
// NOTE: does not check for duplicates.
//
func (l EventListeners) AddEventListener(name string, fn IListen) {
	arr := l[name]
	arr = append(arr, fn)
	l[name] = arr
}

//
// RemoveEventListener for the passed event-function pair which was previously registered with AddEventListener
//
func (l EventListeners) RemoveEventListener(name string, fn IListen) {
	if arr, ok := l[name]; ok {
		for i, f := range arr {
			if f == fn {
				// FIX: add unit test. may assert on end.
				arr = append(arr[:i], arr[i+1:]...)
				break
			}
		}
		l[name] = arr
	}
}

//
// HandleEvents by sending the event to all listeners registered for the event's name.
// When forward is true, the order of handlers is in order of registration.
//
func (l EventListeners) HandleEvents(evt *Proc, forward bool) (err error) {
	name := evt.Name()
	if src, ok := l[name]; ok {
		if cnt := len(src); cnt > 0 {
			// clone the array to avoid add/remove shenanigans in event callbacks.
			// FIX? delaying add / remove commands might be better,
			// flushing those commands here in dispatch.
			temp := make([]IListen, cnt)
			if forward {
				copy(temp, src)
			} else {
				reverseCopy(temp, src)
			}
			for _, fn := range temp {
				if e := fn.HandleEvent(evt); e != nil {
					err = e
					break
				}
				if evt.stopNow {
					break
				}
			}
		}
	}
	return err
}

// helper for EventListeners.dispatch()
func reverseCopy(dst, src []IListen) {
	cnt := len(src)
	for i := 0; i < cnt; i++ {
		dst[(cnt-1)-i] = src[i]
	}
}
