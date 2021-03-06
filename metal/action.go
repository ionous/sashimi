package metal

import (
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

type actionInfo struct {
	mdl *Metal
	*M.ActionModel
}

func (a actionInfo) GetId() ident.Id {
	return a.Id
}

func (a actionInfo) GetActionName() string {
	return a.Name
}

func (a actionInfo) GetEvent() (ret meta.Event) {
	if e, ok := a.mdl.Events[a.EventId]; !ok {
		panic(fmt.Sprintf("internal error, no event found for action %s", a.EventId))
	} else {
		ret = &eventInfo{a.mdl, e}
	}
	return
}

func (a actionInfo) GetNouns() meta.Nouns {
	ret := make(meta.Nouns, len(a.NounTypes))
	for i, c := range a.NounTypes {
		ret[i] = c
	}
	return ret
}

func (a actionInfo) GetCallbacks() (meta.Callbacks, bool) {
	return CallbackList{a.DefaultActions}, len(a.DefaultActions) > 0
}
