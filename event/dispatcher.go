package event

import (
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/sbuf"
)

// EventListeners provides the helper code a way to traverse listeners.
type EventListeners interface {
	NumListener() int
	ListenerNum(int) IListen
}

func Capture(e IEvent, ls EventListeners) (err error) {
	if evt, ok := e.(*Proc); !ok {
		err = errutil.New("unknown event type", sbuf.Type{evt})
	} else {
		// capturing or targeting? trigger capture listeners
		if evt.Phase() != BubblingPhase && !evt.stopMore {
			for i := 0; i < ls.NumListener(); i++ {
				fn := ls.ListenerNum(i)
				if e := fn.HandleEvent(evt); e != nil || evt.stopNow {
					err = e
					break
				}
			}
		}
	}
	return
}

func Bubble(e IEvent, ls EventListeners) (err error) {
	if evt, ok := e.(*Proc); !ok {
		err = errutil.New("unknown event type", sbuf.Type{evt})
	} else {
		// bubbling or targeting? trigger bubble listeners
		if evt.Phase() != CapturingPhase && !evt.stopMore {
			for i := ls.NumListener(); i > 0; i-- {
				fn := ls.ListenerNum(i - 1)
				if e := fn.HandleEvent(evt); e != nil || evt.stopNow {
					err = e
					break
				}
			}
		}
	}
	return
}
