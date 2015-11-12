package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type actionInfo struct {
	mdl *MemoryModel
	*M.ActionInfo
}

func (a actionInfo) GetId() ident.Id {
	return a.Id
}

func (a actionInfo) GetActionName() string {
	return a.ActionName
}

func (a actionInfo) GetEvent() (ret api.Event) {
	if e, ok := a.mdl.Events[a.EventId]; !ok {
		panic(fmt.Sprintf("internal error, no event found for action %s", a.EventId))
	} else {
		ret = eventInfo{a.mdl, e}
	}
	return
}

func (a actionInfo) GetNouns() api.Nouns {
	ret := make(api.Nouns, len(a.NounTypes))
	for i, c := range a.NounTypes {
		ret[i] = c.Id
	}
	return ret
}

func (a actionInfo) GetCallbacks() (ret api.Callbacks, okay bool) {
	if cbs, ok := a.mdl.actions[a.Id]; ok {
		ret, okay = CallbackList{cbs}, true
	}
	return
}
