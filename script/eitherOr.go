package script

import (
	S "github.com/ionous/sashimi/source"
)

//
// Class statement to add a set of enumerated choices for all instances of the class
// TBD: turn into a function chain?
func AreOneOf(name string, or string, more string, rest ...string) ClassEnumFragment {
	origin := NewOrigin(2)
	return ClassEnumFragment{origin, append([]string{name, or, more}, rest...), nil}
}

//
// Class statement to add an either/or choice for all instances of the class
// ex. AreEither("this").Or("that")
func AreEither(firstChoice string) EitherOrPhrase {
	return EitherOrPhrase{firstChoice}
}

func (phrase EitherOrPhrase) Or(secondChoice string) ClassEnumFragment {
	origin := NewOrigin(2)
	return ClassEnumFragment{origin, []string{phrase.firstChoice, secondChoice}, nil}
}

type EitherOrPhrase struct {
	firstChoice string
}

//
// cascades
func (frag ClassEnumFragment) Usually(choice string) ClassEnumFragment {
	// note: its wrong to check for choice in choices: we want to be able to split the declaration
	frag.expects = append(frag.expects, S.NewExpectation(S.UsuallyExpect, choice))
	return frag
}
func (frag ClassEnumFragment) Always(choice string) ClassEnumFragment {
	frag.expects = append(frag.expects, S.NewExpectation(S.AlwaysExpect, choice))
	return frag
}
func (frag ClassEnumFragment) Seldom(choice string) ClassEnumFragment {
	frag.expects = append(frag.expects, S.NewExpectation(S.SeldomExpect, choice))
	return frag
}
func (frag ClassEnumFragment) Never(choice string) ClassEnumFragment {
	frag.expects = append(frag.expects, S.NewExpectation(S.NeverExpect, choice))
	return frag
}

//
// implementation:
//
type ClassEnumFragment struct {
	origin  Origin
	choices []string
	expects []S.PropertyExpectation
}

func (frag ClassEnumFragment) MakeStatement(b SubjectBlock) error {
	name := frag.choices[0] + "Property"
	enum := S.EnumFields{b.subject, name, frag.choices, frag.expects}
	return b.NewEnumeration(enum, frag.origin.Code())
}
