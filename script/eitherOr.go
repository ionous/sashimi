package script

import (
	"fmt"
	S "github.com/ionous/sashimi/source"
	"strings"
)

//
// Class statement to add a set of enumerated choices for all instances of the class
// TBD: turn into a function chain?
func AreOneOf(name string, or string, more string, rest ...string) ClassEnumFragment {
	origin := NewOrigin(2)
	return ClassEnumFragment{origin: origin, choices: append([]string{name, or, more}, rest...)}
}

//
// Class statement to add an either/or choice for all instances of the class
// ex. AreEither("this").Or("that")
func AreEither(firstChoice string) EitherOrPhrase {
	return EitherOrPhrase{firstChoice}
}

func (phrase EitherOrPhrase) Or(secondChoice string) ClassEnumFragment {
	origin := NewOrigin(2)
	return ClassEnumFragment{origin: origin, choices: []string{phrase.firstChoice, secondChoice}}
}

type EitherOrPhrase struct {
	firstChoice string
}

//
// cascades
func (frag ClassEnumFragment) Usually(choice string) ClassEnumFragment {
	// note: its wrong to check for choice in choices: we want to be able to split the declaration
	//internal.expects = append(frag.expects, S.NewExpectation(S.UsuallyExpect, choice))
	frag.usually = choice
	return frag
}

// func (frag ClassEnumFragment) Always(choice string) ClassEnumFragment {
// 	frag.expects = append(frag.expects, S.NewExpectation(S.AlwaysExpect, choice))
// 	return frag
// }
// func (frag ClassEnumFragment) Seldom(choice string) ClassEnumFragment {
// 	frag.expects = append(frag.expects, S.NewExpectation(S.SeldomExpect, choice))
// 	return frag
// }
// func (frag ClassEnumFragment) Never(choice string) ClassEnumFragment {
// 	frag.expects = append(frag.expects, S.NewExpectation(S.NeverExpect, choice))
// 	return frag
// }

//
// implementation:
//
type ClassEnumFragment struct {
	origin  Origin
	choices []string
	usually string
	//expects []S.PropertyExpectation
}

func (frag ClassEnumFragment) MakeStatement(b SubjectBlock) error {
	name := frag.choices[0] //-property
	enum := S.EnumFields{b.subject, name, frag.choices}
	if frag.usually != "" {
		for i, v := range frag.choices {
			if strings.EqualFold(v, frag.usually) {
				frag.choices[0], frag.choices[i] = frag.choices[i], frag.choices[0]
				break
			}
		}
		if frag.choices[0] != frag.usually {
			return fmt.Errorf("usually not found %s", frag.usually)
		}
	}

	return b.NewEnumeration(enum, frag.origin.Code())
}
