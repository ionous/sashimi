package internal

import (
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/util/ident"
)

type EventQueue struct {
	*E.Queue
}

func (q EventQueue) QueueEvent(target E.ITarget, id ident.Id, data interface{}) {
	msg := E.Message{Id: id, Data: data}
	q.Enqueue(target, msg)
}
