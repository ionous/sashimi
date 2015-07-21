package script

import (
	S "github.com/ionous/sashimi/source"
)

//
// Instance statement to set the value of a property
// The property must (eventually) be declared for the class via Have().
//
func Has(property string, values ...interface{}) (frag IFragment) {
	origin := NewOrigin(2)
	switch len(values) {
	case 0:
		frag = Is(property)
	case 1:
		frag = NewFunctionFragment(func(b SubjectBlock) error {
			fields := S.KeyValueFields{b.subject, property, values[0]}
			return b.NewKeyValue(fields, origin.Code())
		})
	default:
		frag = NewFunctionFragment(func(b SubjectBlock) error {
			fields := S.KeyValueFields{b.subject, property, values}
			return b.NewKeyValue(fields, origin.Code())
		})
	}
	return frag
}

//
// Instance statement to select the value of an enumerations
// The enumeration must (eventually) be declared for the class via AreEither(), or AreOneOf()
//
func Is(choice string, choices ...string) IsPhrase {
	origin := NewOrigin(2)
	return IsPhrase{origin, append(choices, choice)}
}

func (phrase IsPhrase) And(choice string) IsPhrase {
	phrase.choices = append(phrase.choices, choice)
	return phrase
}

type IsPhrase struct {
	origin  Origin
	choices []string
}

func (phrase IsPhrase) MakeStatement(b SubjectBlock) (err error) {
	for _, choice := range phrase.choices {
		fields := S.ChoiceFields{b.subject, choice}
		if e := b.NewChoice(fields, phrase.origin.Code()); e != nil {
			err = e
			break
		}
	}
	return err
}
