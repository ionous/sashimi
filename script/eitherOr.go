package script

import (
	S "github.com/ionous/sashimi/source"
)

//
// Class statement to add a set of enumerated choices for all instances of the class
// TBD: turn into a function chain?
func AreOneOf(name string, or string, more string, rest ...string) ClassEnumFragment {
	return ClassEnumFragment{NewOrigin(1), append([]string{name, or, more}, rest...), nil}
}

//
// Class statement to add an either/or choice for all instances of the class
// ex. AreEither("this").Or("that")
func AreEither(firstChoice string) EitherOrPhrase {
	return EitherOrPhrase{firstChoice}
}

func (this EitherOrPhrase) Or(secondChoice string) ClassEnumFragment {
	return ClassEnumFragment{NewOrigin(1), []string{this.firstChoice, secondChoice}, nil}
}

type EitherOrPhrase struct {
	firstChoice string
}

//
// cascades
func (this ClassEnumFragment) Usually(choice string) ClassEnumFragment {
	// note: its wrong to check for choice in choices: we want to be able to split the declaration
	this.expects = append(this.expects, S.NewExpectation(S.UsuallyExpect, choice))
	return this
}
func (this ClassEnumFragment) Always(choice string) ClassEnumFragment {
	this.expects = append(this.expects, S.NewExpectation(S.AlwaysExpect, choice))
	return this
}
func (this ClassEnumFragment) Seldom(choice string) ClassEnumFragment {
	this.expects = append(this.expects, S.NewExpectation(S.SeldomExpect, choice))
	return this
}
func (this ClassEnumFragment) Never(choice string) ClassEnumFragment {
	this.expects = append(this.expects, S.NewExpectation(S.NeverExpect, choice))
	return this
}

//
// implementation:
//
type ClassEnumFragment struct {
	origin  Origin
	choices []string
	expects []S.PropertyExpectation
}

func (this ClassEnumFragment) MakeStatement(b SubjectBlock) error {
	name := this.choices[0] + "Property"
	enum := S.EnumFields{b.subject, name, this.choices, this.expects}
	return b.NewEnumeration(enum, "")
}
