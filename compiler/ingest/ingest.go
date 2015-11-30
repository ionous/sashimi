package ingest

import (
	"encoding/json"
	"fmt"
	M "github.com/ionous/sashimi/compiler/model"
	G "github.com/ionous/sashimi/game"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// implements api.LookupCallbacks
type CodeCalls map[ident.Id]G.Callback

func (m CodeCalls) LookupCallback(id ident.Id) (ret G.Callback, okay bool) {
	if r, ok := m[id]; !ok {
		panic(fmt.Sprintf("couldnt find callback %s", id))
	} else {
		ret, okay = r, ok
	}
	return
}

func GetModelCode(data string, callbacks map[ident.Id]G.Callback) ModelCode {
	var model *M.Model
	if e := json.Unmarshal([]byte(data), &model); e != nil {
		panic(e)
	}
	return ModelCode{model, CodeCalls(callbacks)}
}

type ModelCode struct {
	Model *M.Model
	Code  api.LookupCallbacks
}
