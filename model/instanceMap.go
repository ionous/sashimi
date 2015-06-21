package model

import "fmt"

type InstanceMap map[StringId]*InstanceInfo

// helper to generate an escaped string and an error,
// FIX: return okay, and provide common not found error structures instead
func (this InstanceMap) FindInstance(name string) (*InstanceInfo, bool) {
	safe := MakeStringId(name)
	ret, okay := this[safe]
	return ret, okay
}

func (this InstanceMap) FindInstanceWithClass(name string, class *ClassInfo,
) (ret *InstanceInfo, err error) {
	if inst, ok := this.FindInstance(name); !ok {
		err = InstanceNotFound(name)
	} else if have := inst.Class(); have == class || have.HasParent(class) {
		ret = inst
	} else {
		err = fmt.Errorf("mismatched noun requested: %s,%s!=%s", name, have.Name(), class.Name())
	}

	return ret, err
}

type InstanceNotFound string

func (this InstanceNotFound) Error() string {
	return fmt.Sprintf("instance not found `%s`", string(this))
}
