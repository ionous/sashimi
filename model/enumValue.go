package model

//import "fmt"

//
// Information about a particular choice.
//
type EnumValue struct {
	inst *InstanceInfo
	prop *EnumProperty
}

//
func (enum *EnumValue) Index() (ret int, hadValue bool) {
	id, inst := enum.prop.id, enum.inst
	if index, ok := inst.values[id].(int); ok {
		ret = index
		hadValue = true
	} else if cons, ok := inst.class.PropertyConstraint(enum.prop); ok {
		if cons, ok := cons.(*EnumConstraint); ok {
			ret = cons.BestIndex()
		}
	} else {
		panic("couldnt get constraint")
	}
	return ret, hadValue
}

//
func (enum *EnumValue) Property() *EnumProperty {
	return enum.prop
}

//
func (enum *EnumValue) Any() (value interface{}, hasValue bool) {
	return enum.Index()
}

//
// Returns an impossible string if the choice is invalid.
//
func (enum *EnumValue) String() string {
	index, _ := enum.Index()
	v, _ := enum.prop.IndexToValue(index)
	return v.str
}

//
// Returns an impossible string if the choice is invalid.
//
func (enum *EnumValue) Choice() StringId {
	index, _ := enum.Index()
	v, _ := enum.prop.IndexToValue(index)
	return v.id
}
