package memory

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type actionInfo struct {
	mdl *MemoryModel
	*M.ActionModel
}

func (a actionInfo) GetId() ident.Id {
	return a.Id
}

func (a actionInfo) GetActionName() string {
	return a.Name
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
		ret[i] = c
	}
	return ret
}

func (a actionInfo) GetCallbacks() (api.Callbacks, bool) {
	return CallbackList{a.DefaultActions}, len(a.DefaultActions) > 0
}
