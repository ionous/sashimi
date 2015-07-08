package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

// actions and events have slightly different purposes
// but for now, we put all of the info into action
// and just provide a way to find it by event name
type EventMap map[ident.Id]*ActionInfo

// FIX: i kind of think all maps with the string ids should have this
func (this EventMap) FindEventByName(name string) (ret *ActionInfo, err error) {
	id := MakeStringId(name)
	if act, ok := this[id]; !ok {
		err = EventNotFound{name}
	} else {
		ret = act
	}
	return ret, err
}

//
type EventNotFound struct {
	action string
}

//
func (this EventNotFound) Error() string {
	return fmt.Sprintf("unknown event requested %s", this.action)
}
