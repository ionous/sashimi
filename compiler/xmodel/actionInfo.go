package xmodel

import (
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/sbuf"
)

type ActionInfo struct {
	Id         ident.Id
	ActionName string
	EventId    ident.Id
	NounTypes  []*ClassInfo
}

// FIX: a one to one action/event ratio isnt desirable
type EventInfo struct {
	Id        ident.Id
	EventName string
	ActionId  ident.Id
}

func NewAction(id ident.Id, action string, event ident.Id, classes ...*ClassInfo) (ret *ActionInfo, err error) {
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
		err = sbuf.NewString("bad nouns for").S(action).R(',').V(event).Error()
		// err = fmt.Errorf("bad nouns for %s,%s: %d, %s?", action, event, end, classes)
	} else {
		ret = &ActionInfo{id, action, event, classes[0:end]}
	}
	return ret, err
}

func (a ActionInfo) String() string {
	return a.ActionName
}

func (e EventInfo) String() string {
	return e.EventName
}
