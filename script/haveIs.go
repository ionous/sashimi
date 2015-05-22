package script

import (
	"fmt"
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
		frag = FunctionFragment{func(b SubjectBlock) error {
			fields := S.KeyValueFields{b.subject, property, values[0]}
			return b.NewKeyValue(fields, origin.Code())
		}}
	default:
		frag = FunctionFragment{func(SubjectBlock) error {
			return fmt.Errorf("too many values specified %s", origin)
		}}
	}
	return frag
}

//
// Instance statement to select the value of an enumerations
// The enumeration must (eventually) be declared for the class via AreEither(), or AreOneOf()
//
func Is(choice string, choices ...string) IsPhrase {
	return IsPhrase{append(choices, choice)}
}

func (this IsPhrase) And(choice string) IsPhrase {
	this.choices = append(this.choices, choice)
	return this
}

type IsPhrase struct {
	choices []string
}

func (this IsPhrase) MakeStatement(b SubjectBlock) (err error) {
	for _, choice := range this.choices {
		fields := S.ChoiceFields{b.subject, choice}
		if e := b.NewChoice(fields, ""); e != nil {
			err = e
			break
		}
	}
	return err
}
