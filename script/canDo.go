package script

//Can("show").And("showing").RequiresOne("actor").AndOne("kind"),

// active verb:
func Can(verb string) CanDoPhrase {
	return CanDoPhrase{NewOrigin(1), verb}
}

//starts a requirements phrase for deciding how to provide nouns...
func (this CanDoPhrase) And(doing string) RequiresWhatPhrase {
	return RequiresWhatPhrase{this, doing}
}

// the target will be the same as the source
func (this RequiresWhatPhrase) RequiresNothing() IFragment {
	return ActionAssertionFragment{RequiresWhatPhrase: this}
}

// the target and the context will be input by the user, and will both be of the passed class
func (this RequiresWhatPhrase) RequiresTwo(class string) IFragment {
	return ActionAssertionFragment{RequiresWhatPhrase: this, target: class, context: class}
}

// the target will be input by the user, and will of the passed class
func (this RequiresWhatPhrase) RequiresOne(class string) ActionAssertionFragment {
	return ActionAssertionFragment{RequiresWhatPhrase: this, target: class}
}

// the context will be input by the user, and will of the passed class
func (this ActionAssertionFragment) AndOne(class string) IFragment {
	this.context = class
	return this
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

func (this ActionAssertionFragment) MakeStatement(b SubjectBlock) error {
	return b.NewActionAssertion(this.actionName, this.eventName,
		b.subject, this.target, this.context)
}
