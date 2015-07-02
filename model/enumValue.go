package model

import "fmt"

//
// Information about a particular choice.
//
type EnumValue struct {
	inst *InstanceInfo
	prop *EnumProperty
	cons *EnumConstraint
}

//
func (this *EnumValue) Property() IProperty {
	return this.prop
}

//
func (this *EnumValue) Index() (ret int, hadValue bool) {
	id := this.prop.id
	if index, ok := this.inst.values[id].(int); ok {
		ret = index
		hadValue = true
	} else {
		econ := this.Constraint()
		ret = econ.BestIndex()
	}
	return ret, hadValue
}

//
func (this *EnumValue) Constraint() *EnumConstraint {
	if this.cons == nil {
		id := this.prop.id
		con := this.inst.class.ConstraintById(id)
		if econ, ok := con.(*EnumConstraint); ok {
			this.cons = econ
		} else {
			this.cons = NewConstraint(this.prop.Enumeration)
		}

	}
	return this.cons
}

//
func (this *EnumValue) Any() (value interface{}, hasValue bool) {
	return this.Index()
}

// FIX: once this is re/moved from the model to the compiler level,
// then change the values from pointers to concrete types
func (this *EnumValue) SetAny(value interface{}) (err error) {
	switch choice := value.(type) {
	case int:
		err = this.SetByIndex(choice)
	case string:
		choiceId := MakeStringId(choice)
		if idx, e := this.prop.ChoiceToIndex(choiceId); e != nil {
			err = e
		} else {
			err = this.SetByIndex(idx)
		}
	default:
		err = fmt.Errorf(" %#v set any failed with %#v", this, value)
	}
	return err
}

func (this *EnumValue) SetByIndex(choice int) (err error) {
	constraints := this.Constraint()
	if e := constraints.CheckIndex(choice); e != nil {
		err = e
	} else {
		id := this.prop.id
		this.inst.values[id] = choice
	}
	return err
}

//
// Returns an impossible string if the choice is invalid.
//
func (this *EnumValue) String() string {
	index, _ := this.Index()
	v, _ := this.prop.IndexToValue(index)
	return v.str
}

//
// Returns an impossible string if the choice is invalid.
//
func (this *EnumValue) Choice() StringId {
	index, _ := this.Index()
	v, _ := this.prop.IndexToValue(index)
	return v.id
}
