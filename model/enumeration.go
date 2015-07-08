package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

//
//
// An indexed set of values.
//
type Enumeration struct {
	choices map[ident.Id]int // index + 1, to handle zero/nil
	values  []StringPair     // 0-based index
}

//
//
//
func NewEnumeration(values []string) Enumeration {
	choices, outvalues := make(map[ident.Id]int), make([]StringPair, 0, len(values))
	for i, v := range values {
		s := MakeStringId(v)
		outvalues = append(outvalues, StringPair{s, v})
		choices[s] = i + 1
	}
	return Enumeration{choices, outvalues}
}

//
//
//
func CheckedEnumeration(values []string) (ret Enumeration, err error) {
	enum := NewEnumeration(values)
	switch {
	case enum.choices[""] != 0:
		err = EmptyValueError(values)
	case len(enum.choices) == 0:
		err = EmptyValueError(values)
	case len(enum.choices) != len(values):
		err = MultiplyDefinedError(values)
	}
	if err == nil {
		ret = enum
	}
	return ret, err
}

//
//
//
func (enum Enumeration) Values() []StringPair {
	return enum.values
}

//
//
//
func (enum Enumeration) IndexToValue(index int) (ret StringPair, err error) {
	inRange := index > 0 && index <= len(enum.values)
	if inRange {
		ret = enum.values[index-1]
	} else {
		err = OutOfRangeError(enum, index)
	}
	return ret, err
}

//
//
//
func (enum Enumeration) IndexToChoice(index int) (ret ident.Id, err error) {
	if value, e := enum.IndexToValue(index); e != nil {
		err = e
	} else {
		ret = value.id
	}
	return ret, err
}

//
//
//
func (enum Enumeration) StringToIndex(choice string) (index int, err error) {
	safer := MakeStringId(choice)
	return enum.ChoiceToIndex(safer)
}

//
//
//
func (enum Enumeration) ChoiceToIndex(choice ident.Id) (ret int, err error) {
	if idx, ok := enum.choices[choice]; !ok {
		err = OutOfRangeError(enum, choice)
	} else {
		ret = idx
	}
	return ret, err
}

//
//
//
func MultiplyDefinedError(values []string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("multiple values defined for enum %v", values)
	})
}

//
//
//
func EmptyValueError(values []string) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("empty values defined for enum %v", values)
	})
}

//
//
//
func OutOfRangeError(enum Enumeration, value interface{}) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%v out of range for %v", value, enum)
	})
}

//
//
//
func InvalidChoiceError(enum Enumeration, choice ident.Id) error {
	return errutil.Func(func() string {
		return fmt.Sprintf("%v is an disallowed choice for %v", choice, enum)
	})
}
