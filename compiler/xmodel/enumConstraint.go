package xmodel

import (
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
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
	return
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
	return
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
	return
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
	return
}

//
func (cons *EnumConstraint) Copy() IConstrain {
	never := make(map[ident.Id]bool)
	for k, v := range cons.Never {
		never[k] = v
	}
	return &EnumConstraint{cons.Enumeration, cons.Only, never, cons.Usual, false}
}

func (cons *EnumConstraint) CheckChoice(choice ident.Id) (err error) {
	if index, e := cons.Enumeration.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	}
	return
}

//
func (cons *EnumConstraint) CheckIndex(index int) (err error) {
	if choice, e := cons.Enumeration.IndexToValue(index); e != nil {
		err = e
	} else if cons.Never[choice.Id] {
		err = InvalidChoiceError(cons.Enumeration, choice.Id)
	} else if !cons.Only.Empty() && cons.Only != choice.Id {
		err = InvalidChoiceError(cons.Enumeration, choice.Id)
	}
	return
}

// Always is forever.
func (cons *EnumConstraint) always(choice ident.Id) (err error) {
	if index, e := cons.Enumeration.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// first always assertion?
		case cons.Only.Empty():
			cons.Only = choice

		// ignore duplicate assertions
		case cons.Only == choice:
			break

		// some other always assertion
		default:
			err = errutil.New("always respecified", cons.Only, choice)
		}
	}
	return
}

// Usually is a loose recommendation, and can be overriden for each new owner
func (cons *EnumConstraint) usually(choice ident.Id) (err error) {
	if index, e := cons.Enumeration.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := cons.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// when something is first always 'x' and later usually 'x', that's okay.
		case !cons.Only.Empty() && cons.Only != choice:
			err = errutil.New("usually was always", choice, cons.Only)

		case cons.UsuallyLocal && cons.Usual != choice:
			err = errutil.New("usually was usually", choice, cons.Usual)

		default:
			cons.Usual = choice
			cons.UsuallyLocal = true
		}
	}
	return
}

// Seldom would be complex; instead limiting to the "inverse" of usually
// ( it really only makes sense if there are two choices,
// and the other can become usually. )
func (cons *EnumConstraint) seldom(choice ident.Id) (err error) {
	for other, _ := range cons.Enumeration.Choices {
		if other != choice {
			err = cons.Usually(other)
			if err != nil {
				break
			}
		}
	}
	return
}

// NOTE: cant exclude the final remaining choice
func (cons *EnumConstraint) exclude(choice ident.Id) (err error) {
	// ses raw index to allow the same choice to be excluded multiple times
	_, err = cons.Enumeration.ChoiceToIndex(choice)
	switch {
	// already an error, so do nothing
	case err != nil:
		break

	// "never" cant ever cancel an "always"
	case cons.Only == choice:
		err = errutil.New("never was always", cons.Only)

	// "never" cant cancel a "usually" from the same owner
	case cons.UsuallyLocal && cons.Usual == choice:
		err = errutil.New("usually respecified as never", cons.Usual)

	// okay, exclude:
	default:
		// ignore if the choice was already removed by another never
		if !cons.Never[choice] {
			olen := len(cons.Enumeration.Values)
			xlen := len(cons.Never)
			if olen-xlen-1 < 1 {
				err = errutil.New("never would remove all choices", choice)
			} else {
				cons.Never[choice] = true
			}
		}

		// never *can* cancel a usual from some earlier owner
		// but, note: don't turn it off unless we are actually excluding it.
		if err == nil && cons.Usual == choice {
			cons.Usual = ident.Empty()
		}
	}
	return
}
