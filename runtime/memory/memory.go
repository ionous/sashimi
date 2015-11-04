package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type MemoryModel struct {
	model *M.Model
}

func NewMemoryModel(model *M.Model) MemoryModel {
	return MemoryModel{model}
}

func (m MemoryModel) GetAction(id ident.Id) (ret api.Action, okay bool) {
	if a, ok := m.model.Actions[id]; ok {
		ret, okay = actionInfo{m.model, a}, true
	}
	return
}

func (m MemoryModel) GetParserActions(visit func(api.Action, []string) bool) {
	for _, p := range m.model.ParserActions {
		act := actionInfo{m.model, p.Action}
		if visit(act, p.Commands) {
			break
		}
	}
}

type actionInfo struct {
	model *M.Model
	*M.ActionInfo
}

func (a actionInfo) GetId() ident.Id {
	return a.Id
}

func (a actionInfo) GetActionName() string {
	return a.ActionName
}

func (a actionInfo) GetEventName() string {
	return a.EventName
}

func (a actionInfo) GetNouns() api.Nouns {
	ret := make(api.Nouns, len(a.NounTypes))
	for i, c := range a.NounTypes {
		ret[i] = c.Id
	}
	return ret
}
