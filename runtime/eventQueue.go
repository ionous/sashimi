package runtime

import (
	"fmt"
	E "github.com/ionous/sashimi/event"
)

// NewEventFrame returns a function for defer() end of event.
// FIX? does the event data need to be copied as well?
// FIX: this is bit uglys gross.
type EventFrame func(E.ITarget, *E.Message) func()

type EventQueue struct {
	*E.Queue
}

func (q EventQueue) QueueEvent(target E.ITarget, name string, data interface{}) {
	msg := E.Message{Name: name, Data: data}
	q.Enqueue(target, msg)
}

func (f EventFrame) SendMessage(tgt E.ITarget, msg *E.Message) (err error) {
	defer f(tgt, msg)()
	path := E.NewPathTo(tgt)

	// game.log.Printf("sending `%s` to: %s", msg.Name, path)
	if runDefault, e := msg.Send(path); e != nil {
		err = e
	} else {
		if runDefault {
			if act, ok := msg.Data.(*RuntimeAction); !ok {
				err = fmt.Errorf("unknown action data %T", msg.Data)
			} else {
				act.runDefaultActions()
			}
		}
	}
	return err
}

func DefaultEventFrame(E.ITarget, *E.Message) func() {
	return func() {}
}