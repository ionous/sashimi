package model

import (
	"fmt"
	"github.com/ionous/sashimi/util/ident"
)

type ActionInfo struct {
	Id                    ident.Id
	ActionName, EventName string
	NounTypes             []*ClassInfo
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

func (a *ActionInfo) String() string {
	return a.Id.String()
}
