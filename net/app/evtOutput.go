package app

import (
	"container/list"
	E "github.com/ionous/sashimi/event"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/util/ident"
)

// i wonder whether its possible to put ActorSays into an event --
// raise speak on the actor, and the action is "say"
// in that case the "x-rel" etc actions would become events
// and there would only be one action: print.
// [ unless we send the reports, logs, prints ]

// in fact, having seen the output now --
// id rather that x-rel etc were hierarchial events
// the "change-states" could be part of event data maybe
// and any text could be ... meta, or its own event data....
// its not super easy to change right now because the x-rel bits are so separate from the rest of the event/open/close system. that might be a deeper change, which would allow stories to watch for property changes too.

// currently, event data is not transfered; only action data.

// EventStream provides hierarchical event output.
type EventStream struct {
	list   *list.List //EventBlock
	events []*EventBlock
}

type EventBlock struct {
	Evt    string           `json:"evt,omitempty"`
	Tgt    *resource.Object `json:"tgt,omitempty"`
	Data   interface{}      `json:"data,omitempty"`
	Events []*EventBlock    `json:"events,omitempty"`
}

func NewEventStream() *EventStream {
	return &EventStream{list: list.New()}
}

func (evs *EventStream) Flush() (ret []*EventBlock) {
	ret, evs.events = evs.events, nil
	return ret
}

func (evs *EventStream) CurrentEvent() (ret *EventBlock) {
	if back := evs.list.Back(); back != nil {
		ret = (back.Value).(*EventBlock)
	}
	return ret
}

func (evs *EventStream) PushEvent(evt ident.Id, tgt E.ITarget, data interface{}) (ret *EventBlock) {
	// create a new event block, and add it ( as the current event )
	noRef := resource.ObjectList{}.NewObject(
		jsonId(tgt.Id()),
		jsonId(tgt.Class()))
	return evs.push(evt, noRef, data)
}

func (evs *EventStream) push(evt ident.Id, tgt *resource.Object, data interface{}) (ret *EventBlock) {
	// get parent before push
	parent := evs.CurrentEvent()
	block := &EventBlock{Evt: jsonId(evt), Tgt: tgt, Data: data}
	// link this event into its parent (if any)
	if parent != nil {
		parent.Events = append(parent.Events, block)
	} else {
		evs.events = append(evs.events, block)
	}
	// set this event as the new current event
	evs.list.PushBack(block)
	// return the new event
	return block
}

// AddAction adds an event without hierarchy
func (evs *EventStream) AddAction(act string, tgt *resource.Object, data interface{}) {
	evt := ident.MakeId(act)
	evs.push(evt, tgt, data)
	evs.PopEvent()
}

func (evs *EventStream) PopEvent() {
	evs.list.Remove(evs.list.Back())
}
