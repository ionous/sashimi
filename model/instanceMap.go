package model

import "fmt"

type InstanceMap map[StringId]*InstanceInfo

// helper to generate an escaped string and an error,
// FIX: return okay, and provide common not found error structures instead
func (this InstanceMap) FindInstance(name string) (ret *InstanceInfo, err error) {
	safe := MakeStringId(name)
	if inst, ok := this[safe]; ok {
		ret = inst
	} else {
		err = InstanceNotFound(name)
	}
	return ret, err
}

func (this InstanceMap) FindInstanceWithClass(name string, class *ClassInfo,
) (ret *InstanceInfo, err error) {
	inst, e := this.FindInstance(name)
	if e != nil {
		err = e
	} else {
		haveClass := inst.Class()
		if haveClass == class || haveClass.HasParent(class) {
			ret = inst
		} else {
			err = fmt.Errorf("mismatched noun requested: %s,%s!=%s", name, haveClass.Name(), class.Name())
		}
	}
	return ret, err
}

type InstanceNotFound string

func (this InstanceNotFound) Error() string {
	return fmt.Sprintf("instance not found `%s`", string(this))
}
