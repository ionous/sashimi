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
	if index, ok := this.inst.enum[id]; ok {
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

//
func (this *EnumValue) SetAny(value interface{}) (err error) {
	if choice, ok := value.(int); !ok {
		err = fmt.Errorf(" %#v set any failed with %#v", this, value)
	} else {
		id := this.prop.id
		econ := this.Constraint()
		if e := econ.CheckIndex(choice); e != nil {
			err = e
		} else {
			this.inst.enum[id] = choice
		}
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
