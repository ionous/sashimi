package model

import "fmt"

//
//
// An indexed set of values.
//
type Enumeration struct {
	choices map[StringId]int // index + 1, to handle zero/nil
	values  []StringPair     // 0-based index
}

//
//
//
func NewEnumeration(values []string) Enumeration {
	choices, outvalues := make(map[StringId]int), make([]StringPair, 0, len(values))
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
func CheckedEnumeration(values []string) (ret *Enumeration, err error) {
	enum := NewEnumeration(values)
	switch {
	case enum.choices[""] != 0:
		err = EmptyValueError{values}
	case len(enum.choices) == 0:
		err = EmptyValueError{values}
	case len(enum.choices) != len(values):
		err = MultiplyDefinedError{values}
	}
	if err == nil {
		ret = &enum
	}
	return ret, err
}

//
//
//
func (this *Enumeration) Values() []StringPair {
	return this.values
}

//
//
//
func (this *Enumeration) IndexToValue(index int) (ret StringPair, err error) {
	inRange := index > 0 && index <= len(this.values)
	if inRange {
		ret = this.values[index-1]
	} else {
		err = OutOfRangeError{this, index}
	}
	return ret, err
}

//
//
//
func (this *Enumeration) IndexToChoice(index int) (ret StringId, err error) {
	if value, e := this.IndexToValue(index); e != nil {
		err = e
	} else {
		ret = value.id
	}
	return ret, err
}

//
//
//
func (this *Enumeration) StringToIndex(choice string) (index int, err error) {
	safer := MakeStringId(choice)
	return this.ChoiceToIndex(safer)
}

//
//
//
func (this *Enumeration) ChoiceToIndex(choice StringId) (index int, err error) {
	index = this.choices[choice]
	if !(index > 0) {
		err = OutOfRangeError{this, choice}
	}
	return index, err
}

//
//
//
type MultiplyDefinedError struct {
	values []string
}

func (this MultiplyDefinedError) Error() string {
	return fmt.Sprintf("multiple values defined for enum %v", this.values)
}

//
//
//
type EmptyValueError struct {
	values []string
}

func (this EmptyValueError) Error() string {
	return fmt.Sprintf("empty values defined for enum %v", this.values)
}

//
//
//
type OutOfRangeError struct {
	enum  *Enumeration
	value interface{}
}

func (this OutOfRangeError) Error() string {
	return fmt.Sprintf("%v out of range for %v", this.value, this.enum)
}

//
//
//
type InvalidChoiceError struct {
	enum   *Enumeration
	choice StringId
}

func (this InvalidChoiceError) Error() string {
	return fmt.Sprintf("%v is an disallowed choice for %v", this.choice, this.enum)
}
