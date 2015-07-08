package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

type ActionInfo struct {
	id            ident.Id
	action, event string
	nounTypes     []*ClassInfo
}

func NewAction(id ident.Id, action, event string, classes ...*ClassInfo) (ret *ActionInfo, err error) {
	end, found := len(classes), false
	for i := len(classes) - 1; i >= 0; i-- {
		c := classes[i]
		if c != nil {
			found = true
		} else {
			if found {
				end = -1
				break
			}
			end = i
		}
	}
	if end < 1 || end > 3 {
		err = fmt.Errorf("bad nouns for %s,%s: %d, %s?", action, event, end, classes)
	} else {
		ret = &ActionInfo{id, action, event, classes[0:end]}
	}
	return ret, err
}

// based on the action name
func (this *ActionInfo) Id() ident.Id {
	return this.id
}

// baed on id
func (this *ActionInfo) String() string {
	return this.id.String()
}

// action name
func (this *ActionInfo) Action() string {
	return this.action
}

// event name
func (this *ActionInfo) Event() string {
	return this.event
}

//
func (this *ActionInfo) Source() *ClassInfo {
	return this.neverSimple(0)
}

//
func (this *ActionInfo) Target() *ClassInfo {
	return this.neverSimple(1)
}

//
func (this *ActionInfo) Context() *ClassInfo {
	return this.neverSimple(2)
}

//
func (this *ActionInfo) neverSimple(i int) (ret *ClassInfo) {
	if i < len(this.nounTypes) {
		ret = this.nounTypes[i]
	}
	return ret
}

//
func (this *ActionInfo) NounSlice() (ret []*ClassInfo) {
	return this.nounTypes
}

//
func (this *ActionInfo) NumNouns() (ret int) {
	return len(this.nounTypes)
}
