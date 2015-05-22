package model

import (
	"fmt"
	"math"
)

type EnumConstraint struct {
	enum Enumeration

	// FIX? would it be better to use index internally for these things?
	only         StringId
	never        map[StringId]bool
	usual        StringId
	usuallyLocal bool // usual set for this constraint, or for some ancestor?
}

//
//
//
func NewConstraint(enum Enumeration) *EnumConstraint {
	never := make(map[StringId]bool)
	return &EnumConstraint{enum: enum, never: never}
}

//
//
//
func (this *EnumConstraint) Always(value interface{}) (err error) {
	switch choice := value.(type) {
	case StringId:
		err = this.always(choice)
	case string:
		err = this.always(MakeStringId(choice))
	default:
		err = OutOfRangeError{&this.enum, value}
	}
	return err
}

//
//
//
func (this *EnumConstraint) Usually(value interface{}) (err error) {
	switch choice := value.(type) {
	case StringId:
		err = this.usually(choice)
	case string:
		err = this.usually(MakeStringId(choice))
	default:
		err = OutOfRangeError{&this.enum, value}
	}
	return err
}

//
//
//
func (this *EnumConstraint) Seldom(value interface{}) (err error) {
	switch choice := value.(type) {
	case StringId:
		err = this.seldom(choice)
	case string:
		err = this.seldom(MakeStringId(choice))
	default:
		err = OutOfRangeError{&this.enum, value}
	}
	return err
}

//
//
//
func (this *EnumConstraint) Exclude(value interface{}) (err error) {
	switch choice := value.(type) {
	case StringId:
		err = this.exclude(choice)
	case string:
		err = this.exclude(MakeStringId(choice))
	default:
		err = OutOfRangeError{&this.enum, value}
	}
	return err
}

//
//
//
func (this *EnumConstraint) Copy() IConstrain {
	never := make(map[StringId]bool)
	for k, v := range this.never {
		never[k] = v
	}
	return &EnumConstraint{this.enum, this.only, never, this.usual, false}
}

//
//
//
func (this *EnumConstraint) CheckIndex(index int) (err error) {
	if choice, e := this.enum.IndexToValue(index); e != nil {
		err = e
	} else if this.never[choice.id] {
		err = InvalidChoiceError{&this.enum, choice.id}
	} else if this.only != "" && this.only != choice.id {
		err = InvalidChoiceError{&this.enum, choice.id}
	}
	return err
}

//
//
//
func (this *EnumConstraint) BestIndex() (index int) {
	switch {
	case this.only != "":
		index = this.enum.choices[this.only]

	case this.usual != "":
		index = this.enum.choices[this.usual]

	default:
		smallestIndex := math.MaxInt32
		for k, i := range this.enum.choices {
			if !this.never[k] && i < smallestIndex {
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
func (this *EnumConstraint) always(choice StringId) (err error) {
	if index, e := this.enum.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := this.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// first always assertion?
		case this.only == "":
			this.only = choice

		// ignore duplicate assertions
		case this.only == choice:
			break

		// some other always assertion
		default:
			err = fmt.Errorf("always %v respecified as %v", this.only, choice)
		}
	}
	return err
}

//
// Usually is a loose recommendation, and can be overriden for each new owner
//
func (this *EnumConstraint) usually(choice StringId) (err error) {
	if index, e := this.enum.ChoiceToIndex(choice); e != nil {
		err = e
	} else if e := this.CheckIndex(index); e != nil {
		err = e
	} else {
		switch {
		// when something is first always 'x' and later usually 'x', that's okay.
		case this.only != "" && this.only != choice:
			err = fmt.Errorf("usually '%v' was always '%v'", choice, this.only)

		case this.usuallyLocal && this.usual != choice:
			err = fmt.Errorf("usually `%v` was usually `%v`", choice, this.usual)

		default:
			this.usual = choice
			this.usuallyLocal = true
		}
	}
	return err
}

//
// Seldom would be complex; instead limiting to the "inverse" of usually
// ( it really only makes sense if there are two choices,
// and the other can become usually. )
//
func (this *EnumConstraint) seldom(choice StringId) (err error) {
	for other, _ := range this.enum.choices {
		if other != choice {
			err = this.Usually(other)
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
func (this *EnumConstraint) exclude(choice StringId) (err error) {
	// ses raw index to allow the same choice to be excluded multiple times
	_, err = this.enum.ChoiceToIndex(choice)
	switch {
	// already an error, so do nothing
	case err != nil:
		break

	// "never" cant ever cancel an "always"
	case this.only == choice:
		err = fmt.Errorf("never %v was always", this.only)

	// "never" cant cancel a "usually" from the same owner
	case this.usuallyLocal && this.usual == choice:
		err = fmt.Errorf("usually %v respecified as never", this.usual)

	// okay, exclude:
	default:
		// ignore if the choice was already removed by another never
		if !this.never[choice] {
			olen := len(this.enum.values)
			xlen := len(this.never)
			if olen-xlen-1 < 1 {
				err = fmt.Errorf("never %s would remove all choices", choice)
			} else {
				this.never[choice] = true
			}
		}

		// never *can* cancel a usual from some earlier owner
		// but, note: don't turn it off unless we are actually excluding it.
		if err == nil && this.usual == choice {
			this.usual = ""
		}
	}

	return err
}
