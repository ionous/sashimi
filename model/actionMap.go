package model

import "fmt"

type ActionMap map[StringId]*ActionInfo

// FIX: i kind of think all maps with the string ids should have this
func (this ActionMap) FindActionByName(name string) (ret *ActionInfo, err error) {
	id := MakeStringId(name)
	if act, ok := this[id]; !ok {
		err = ActionNotFound{name}
	} else {
		ret = act
	}
	return ret, err
}

//
type ActionNotFound struct {
	action string
}

//
func (this ActionNotFound) Error() string {
	return fmt.Sprintf("unknown action requested %s", this.action)
}
