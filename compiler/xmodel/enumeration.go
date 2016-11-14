package xmodel

import (
	"github.com/ionous/sashimi/util/errutil"
	"github.com/ionous/sashimi/util/ident"
)

// An indexed set of values.
type Enumeration struct {
	Choices map[ident.Id]int // index + 1, to handle zero/nil
	Values  []StringPair     // 0-based index
}

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
func CheckedEnumeration(values []string) (ret Enumeration, err error) {
	enum := NewEnumeration(values)
	switch {
	case enum.Choices[ident.Empty()] != 0:
		err = EmptyValueError(values)
	case len(enum.Choices) == 0:
		err = EmptyValueError(values)
	case len(enum.Choices) != len(values):
		err = MultiplyDefinedError(values)
	}
	if err == nil {
		ret = enum
	}
	return
}

//
func (enum Enumeration) IndexToValue(index int) (ret StringPair, err error) {
	inRange := index > 0 && index <= len(enum.Values)
	if inRange {
		ret = enum.Values[index-1]
	} else {
		err = OutOfRangeError(enum, index)
	}
	return
}

//
func (enum Enumeration) IndexToChoice(index int) (ret ident.Id, err error) {
	if value, e := enum.IndexToValue(index); e != nil {
		err = e
	} else {
		ret = value.Id
	}
	return
}

//
func (enum Enumeration) ChoiceToIndex(choice ident.Id) (ret int, err error) {
	if idx, ok := enum.Choices[choice]; !ok {
		err = OutOfRangeError(enum, choice)
	} else {
		ret = idx
	}
	return
}

//
func MultiplyDefinedError(values []string) error {
	return errutil.New("multiple values defined for enum", values)
}

//
func EmptyValueError(values []string) error {
	return errutil.New("empty values defined for enum", values)
}

//
func OutOfRangeError(enum Enumeration, value interface{}) error {
	return errutil.New(value, "out of range for", enum)
}

//
func InvalidChoiceError(enum Enumeration, choice ident.Id) error {
	return errutil.New(choice, "is an disallowed choice for", enum)
}
