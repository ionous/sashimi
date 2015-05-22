package event

import "fmt"

//
// Internal helper for default dispatcher:
// event name => list of listeners.
//
type EventListeners map[string]EventList

//
// List of listeners.
//
type EventList []IListen

//
// Add a new listener for the passed event name.
// NOTE: does not check for duplicates.
//
func (this EventListeners) AddEventListener(name string, fn IListen) {
	arr := this[name]
	arr = append(arr, fn)
	this[name] = arr
}

//
// Remove the listener for the passed event name.
//
func (this EventListeners) RemoveEventListener(name string, fn IListen) {
	if arr, ok := this[name]; ok {
		for i, f := range arr {
			if f == fn {
				// FIX: add unit test. may assert on end.
				arr = append(arr[:i], arr[i+1:]...)
				break
			}
		}
		this[name] = arr
	}
}

//
// Send the event to all listeners registered for the event's name.
// When forward is true, the order of handlers is in order of registration.
//
func (this EventListeners) HandleEvents(evt IEvent, forward bool) (err error) {
	// FIX: pass the proc in directly?
	if proc, ok := evt.(*Proc); !ok {
		err = fmt.Errorf("unknown event type %T", evt)
	} else {
		name := evt.Name()
		if src, ok := this[name]; ok {
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
					if proc.stopNow {
						break
					}
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
