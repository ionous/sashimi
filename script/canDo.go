package script

import (
	S "github.com/ionous/sashimi/source"
)

// active verb:
func Can(verb string) CanDoPhrase {
	origin := NewOrigin(3)
	return CanDoPhrase{origin, verb}
}

//starts a requirements phrase for deciding how to provide nouns...
func (frag CanDoPhrase) And(doing string) RequiresWhatPhrase {
	return RequiresWhatPhrase{frag, doing}
}

// the target will be the same as the source
func (frag RequiresWhatPhrase) RequiresNothing() IFragment {
	return ActionAssertionFragment{RequiresWhatPhrase: frag}
}

// the target and the context will be input by the user, and will both be of the passed class
func (frag RequiresWhatPhrase) RequiresTwo(class string) IFragment {
	return ActionAssertionFragment{RequiresWhatPhrase: frag, target: class, context: class}
}

// the target will be input by the user, and will of the passed class
func (frag RequiresWhatPhrase) RequiresOne(class string) ActionAssertionFragment {
	return ActionAssertionFragment{RequiresWhatPhrase: frag, target: class}
}

// the context will be input by the user, and will of the passed class
func (frag ActionAssertionFragment) AndOne(class string) IFragment {
	frag.context = class
	return frag
}

//
type CanDoPhrase struct {
	Origin
	actionName string
}
type RequiresWhatPhrase struct {
	CanDoPhrase
	eventName string
}
type ActionAssertionFragment struct {
	RequiresWhatPhrase
	target, context string
}

func (frag ActionAssertionFragment) MakeStatement(b SubjectBlock) error {
	fields := S.ActionAssertionFields{
		frag.actionName, frag.eventName, b.subject, frag.target, frag.context}
	return b.NewActionAssertion(fields, frag.Code())
}
