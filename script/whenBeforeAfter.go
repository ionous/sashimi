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
	return EventPhrase{[]string{event}, origin, S.ListenTargetOnly | S.ListenCapture}
}

// a shortcut for meaning at the target
// ( queues the callback to run after the default actions have completed. )
func After(event string) EventPhrase {
	// FIX: I moved this to the capture phase so that closer to the instance is later.
	// good, bad? control?
	origin := NewOrigin(2)
	return EventPhrase{[]string{event}, origin, S.ListenTargetOnly | S.ListenCapture | S.ListenRunAfter}
}

// a shortcut for meaning at the target
// ( implemented as a bubbling event )
func When(event string) EventPhrase {
	origin := NewOrigin(2)
	return EventPhrase{[]string{event}, origin, S.ListenTargetOnly}
}

//
func (p EventPhrase) Or(event string) EventPhrase {
	p.events = append(p.events, event)
	return p
}

//
func (p EventPhrase) Always(cb G.Callback) EventFinalizer {
	return EventFinalizer{p, cb}
}

func (p EventPhrase) Go(phrase G.RuntimePhrase, phrases ...G.RuntimePhrase) EventFinalizer {
	return p.Always(func(g G.Play) {
		g.Go(phrase, phrases...)
	})
}

//
func (frag EventFinalizer) MakeStatement(b SubjectBlock) (err error) {
	for _, evt := range frag.events {
		fields := S.ListenFields{b.subject, evt, frag.cb, frag.options}
		if e := b.NewEventHandler(fields, frag.Code()); e != nil {
			err = e
			break
		}
	}
	return err
}

//
type EventPhrase struct {
	events []string // name of the event in question
	Origin
	options S.ListenOptions
}

//
type EventFinalizer struct {
	EventPhrase
	cb G.Callback
}
