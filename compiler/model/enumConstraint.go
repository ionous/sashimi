package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
	"math"
)

type EnumConstraint struct {
	Enumeration
	// FIX? would it be better to use index internally for these things?
	Only         ident.Id
	Never        map[ident.Id]bool
	Usual        ident.Id
	UsuallyLocal bool // usual set for cons constraint, or for some ancestor?
}

//
func NewConstraint(enum Enumeration) *EnumConstraint {
	never := make(map[ident.Id]bool)
	return &EnumConstraint{Enumeration: enum, Never: never}
}

//
func (cons *EnumConstraint) Always(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.always(choice)
	case string:
		err = cons.always(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.Enumeration, value)
	}
	return err
}

//
func (cons *EnumConstraint) Usually(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.usually(choice)
	case string:
		err = cons.usually(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.Enumeration, value)
	}
	return err
}

//
func (cons *EnumConstraint) Seldom(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.seldom(choice)
	case string:
		err = cons.seldom(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.Enumeration, value)
	}
	return err
}

//
func (cons *EnumConstraint) Exclude(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.exclude(choice)
	case string:
		err = cons.exclude(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.Enumeration, value)
	}
	return err
}

//
func (cons *EnumConstraint) Copy() IConstrain {
	never := make(map[ident.Id]bool)
	for k, v := range cons.Never {
		never[k] = v
	}
	return &EnumConstraint{cons.Enumeration, cons.Only, never, cons.Usual, false}
}

//
func (cons *EnumConstraint) CheckIndex(index int) (err error) {
	if choice, e := cons.Enumeration.IndexToValue(index); e != nil {
		err = e
	} else if cons.Never[choice.id] {
		err = InvalidChoiceError(cons.Enumeration, choice.id)
	} else if cons.Only != "" && cons.Only != choice.id {
		err = InvalidChoiceError(cons.Enumeration, choice.id)
	}
	return err
}

//
func (cons *EnumConstraint) BestIndex() (index int) {
	switch {
	case cons.Only != "":
		index = cons.Enumeration.Choices[cons.Only]

	case cons.Usual != "":
		index = cons.Enumeration.Choices[cons.Usual]

	default:
		smallestIndex := math.MaxInt32
		for k, i := range cons.Enumeration.Choices {
			if !cons.Never[k] && i < smallestIndex {
				smallestIndex = i
				index = i
			}
		}
	}
	return index
}

//
// Always is forever.
//
func (cons *EnumConstraint) always(choice ident.Id) (err error) {
	if index, e := cons.Enumeration.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// first always assertion?
		case cons.Only == "":
			cons.Only = choice

		// ignore duplicate assertions
		case cons.Only == choice:
			break

		// some other always assertion
		default:
			err = fmt.Errorf("always %v respecified as %v", cons.Only, choice)
		}
	}
	return err
}

//
// Usually is a loose recommendation, and can be overriden for each new owner
//
func (cons *EnumConstraint) usually(choice ident.Id) (err error) {
	if index, e := cons.Enumeration.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// when something is first always 'x' and later usually 'x', that's okay.
		case cons.Only != "" && cons.Only != choice:
			err = fmt.Errorf("usually '%v' was always '%v'", choice, cons.Only)

		case cons.UsuallyLocal && cons.Usual != choice:
			err = fmt.Errorf("usually `%v` was usually `%v`", choice, cons.Usual)

		default:
			cons.Usual = choice
			cons.UsuallyLocal = true
		}
	}
	return err
}

//
// Seldom would be complex; instead limiting to the "inverse" of usually
// ( it really only makes sense if there are two choices,
// and the other can become usually. )
//
func (cons *EnumConstraint) seldom(choice ident.Id) (err error) {
	for other, _ := range cons.Enumeration.Choices {
		if other != choice {
			err = cons.Usually(other)
			if err != nil {
				break
			}
		}
	}
	return err
}

// NOTE: cant exclude the final remaining choice
//
func (cons *EnumConstraint) exclude(choice ident.Id) (err error) {
	// ses raw index to allow the same choice to be excluded multiple times
	_, err = cons.Enumeration.ChoiceToIndex(choice)
	switch {
	// already an error, so do nothing
	case err != nil:
		break

	// "never" cant ever cancel an "always"
	case cons.Only == choice:
		err = fmt.Errorf("never %v was always", cons.Only)

	// "never" cant cancel a "usually" from the same owner
	case cons.UsuallyLocal && cons.Usual == choice:
		err = fmt.Errorf("usually %v respecified as never", cons.Usual)

	// okay, exclude:
	default:
		// ignore if the choice was already removed by another never
		if !cons.Never[choice] {
			olen := len(cons.Enumeration.Values)
			xlen := len(cons.Never)
			if olen-xlen-1 < 1 {
				err = fmt.Errorf("never %s would remove all choices", choice)
			} else {
				cons.Never[choice] = true
			}
		}

		// never *can* cancel a usual from some earlier owner
		// but, note: don't turn it off unless we are actually excluding it.
		if err == nil && cons.Usual == choice {
			cons.Usual = ""
		}
	}

	return err
}
