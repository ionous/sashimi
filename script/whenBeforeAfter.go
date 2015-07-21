package script

import (
	E "github.com/ionous/sashimi/event"
	G "github.com/ionous/sashimi/game"
	S "github.com/ionous/sashimi/source"
)

// statement to declare an default action handler
func To(action string, c G.Callback) IFragment {
	source := NewOrigin(2)
	return NewFunctionFragment(func(b SubjectBlock) error {
		fields := S.RunFields{b.subject, action, c, E.TargetPhase}
		return b.NewActionHandler(fields, source.Code())
	})
}

//
// FIX: itd be nice to have some sort of wrapper to detect if they are used outside of,
// or rather not consumed by, the(). the wrapper would error at script compile time.

// a shortcut for meaning at the target
// ( implemented as a capturing event )
func Before(event string) EventPhrase {
	origin := NewOrigin(2)
	return EventPhrase{event, origin, S.ListenTargetOnly | S.ListenCapture}
}

// a shortcut for meaning at the target
// ( queues the callback to run after the default actions have completed. )
func After(event string) EventPhrase {
	// FIX: I moved this to the capture phase so that closer to the instance is later.
	// good, bad? control?
	origin := NewOrigin(2)
	return EventPhrase{event, origin, S.ListenTargetOnly | S.ListenCapture | S.ListenRunAfter}
}

// a shortcut for meaning at the target
// ( implemented as a bubbling event )
func When(event string) EventPhrase {
	origin := NewOrigin(2)
	return EventPhrase{event, origin, S.ListenTargetOnly}
}

//
func WhenBubbling(event string, cb G.Callback) EventFinalizer {
	origin := NewOrigin(2)
	return EventPhrase{event, origin, S.ListenBubble}.Always(cb)
}

//
func WhenCapturing(event string, cb G.Callback) EventFinalizer {
	origin := NewOrigin(2)
	return EventPhrase{event, origin, S.ListenCapture}.Always(cb)
}

//
func (phrase EventPhrase) Always(cb G.Callback) EventFinalizer {
	return EventFinalizer{phrase, cb}
}

//
func (frag EventFinalizer) MakeStatement(b SubjectBlock) error {
	fields := S.ListenFields{b.subject, frag.event, frag.cb, frag.options}
	return b.NewEventHandler(fields, frag.Code())
}

//
type EventPhrase struct {
	event string // name of the event in question
	Origin
	options S.ListenOptions
}

//
type EventFinalizer struct {
	EventPhrase
	cb G.Callback
}
