package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
	"math"
)

type EnumConstraint struct {
	enum Enumeration

	// FIX? would it be better to use index internally for these things?
	only         ident.Id
	never        map[ident.Id]bool
	usual        ident.Id
	usuallyLocal bool // usual set for cons constraint, or for some ancestor?
}

//
//
//
func NewConstraint(enum Enumeration) *EnumConstraint {
	never := make(map[ident.Id]bool)
	return &EnumConstraint{enum: enum, never: never}
}

//
//
//
func (cons *EnumConstraint) Always(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.always(choice)
	case string:
		err = cons.always(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.enum, value)
	}
	return err
}

//
//
//
func (cons *EnumConstraint) Usually(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.usually(choice)
	case string:
		err = cons.usually(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.enum, value)
	}
	return err
}

//
//
//
func (cons *EnumConstraint) Seldom(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.seldom(choice)
	case string:
		err = cons.seldom(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.enum, value)
	}
	return err
}

//
//
//
func (cons *EnumConstraint) Exclude(value interface{}) (err error) {
	switch choice := value.(type) {
	case ident.Id:
		err = cons.exclude(choice)
	case string:
		err = cons.exclude(MakeStringId(choice))
	default:
		err = OutOfRangeError(cons.enum, value)
	}
	return err
}

//
//
//
func (cons *EnumConstraint) Copy() IConstrain {
	never := make(map[ident.Id]bool)
	for k, v := range cons.never {
		never[k] = v
	}
	return &EnumConstraint{cons.enum, cons.only, never, cons.usual, false}
}

//
//
//
func (cons *EnumConstraint) CheckIndex(index int) (err error) {
	if choice, e := cons.enum.IndexToValue(index); e != nil {
		err = e
	} else if cons.never[choice.id] {
		err = InvalidChoiceError(cons.enum, choice.id)
	} else if cons.only != "" && cons.only != choice.id {
		err = InvalidChoiceError(cons.enum, choice.id)
	}
	return err
}

//
//
//
func (cons *EnumConstraint) BestIndex() (index int) {
	switch {
	case cons.only != "":
		index = cons.enum.choices[cons.only]

	case cons.usual != "":
		index = cons.enum.choices[cons.usual]

	default:
		smallestIndex := math.MaxInt32
		for k, i := range cons.enum.choices {
			if !cons.never[k] && i < smallestIndex {
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
	if index, e := cons.enum.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// first always assertion?
		case cons.only == "":
			cons.only = choice

		// ignore duplicate assertions
		case cons.only == choice:
			break

		// some other always assertion
		default:
			err = fmt.Errorf("always %v respecified as %v", cons.only, choice)
		}
	}
	return err
}

//
// Usually is a loose recommendation, and can be overriden for each new owner
//
func (cons *EnumConstraint) usually(choice ident.Id) (err error) {
	if index, e := cons.enum.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// when something is first always 'x' and later usually 'x', that's okay.
		case cons.only != "" && cons.only != choice:
			err = fmt.Errorf("usually '%v' was always '%v'", choice, cons.only)

		case cons.usuallyLocal && cons.usual != choice:
			err = fmt.Errorf("usually `%v` was usually `%v`", choice, cons.usual)

		default:
			cons.usual = choice
			cons.usuallyLocal = true
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
	for other, _ := range cons.enum.choices {
		if other != choice {
			err = cons.Usually(other)
			if err != nil {
				break
			}
		}
	}
	return err
}

//
//
// NOTE: cant exclude the final remaining choice
//
func (cons *EnumConstraint) exclude(choice ident.Id) (err error) {
	// ses raw index to allow the same choice to be excluded multiple times
	_, err = cons.enum.ChoiceToIndex(choice)
	switch {
	// already an error, so do nothing
	case err != nil:
		break

	// "never" cant ever cancel an "always"
	case cons.only == choice:
		err = fmt.Errorf("never %v was always", cons.only)

	// "never" cant cancel a "usually" from the same owner
	case cons.usuallyLocal && cons.usual == choice:
		err = fmt.Errorf("usually %v respecified as never", cons.usual)

	// okay, exclude:
	default:
		// ignore if the choice was already removed by another never
		if !cons.never[choice] {
			olen := len(cons.enum.values)
			xlen := len(cons.never)
			if olen-xlen-1 < 1 {
				err = fmt.Errorf("never %s would remove all choices", choice)
			} else {
				cons.never[choice] = true
			}
		}

		// never *can* cancel a usual from some earlier owner
		// but, note: don't turn it off unless we are actually excluding it.
		if err == nil && cons.usual == choice {
			cons.usual = ""
		}
	}

	return err
}
