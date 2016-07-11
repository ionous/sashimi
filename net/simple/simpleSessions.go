package simple

import (
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident"
	"github.com/ionous/sashimi/util/uuid"
	"io/ioutil"
	"sync"
)

type essMap map[string]*SimpleSession

type Sessions struct {
	sessions essMap
	model    meta.Model
	calls    call.MarkerStorage
	*sync.Mutex
}

// FIX? merge with memSessions.go
// maybe making a separate cmd package for wrapping the memsessions with Session
func NewSessions() *Sessions {
	return &Sessions{sessions: make(essMap), Mutex: new(sync.Mutex)}
}

func (ess *Sessions) NewSession() (ret string, err error) {
	if m, e := ess.getModel(); e != nil {
		err = e
	} else if s, e := NewSimpleSession(m, ess.calls); e != nil {
		err = e
	} else {
		id := ident.Dash(uuid.MakeUniqueId())
		defer ess.Unlock()
		ess.Lock()
		ess.sessions[id] = s
		ret = id
	}
	return
}

func (ess Sessions) GetSession(id string) (ret *SimpleSession, okay bool) {
	defer ess.Unlock()
	ess.Lock()
	ret, okay = ess.sessions[id]
	return
}

func (ess *Sessions) getModel() (ret meta.Model, err error) {
	if ess.model != nil {
		ret = ess.model
	} else {
		calls := call.MakeMarkerStorage()
		if m, e := script.InitScripts().CompileCalls(ioutil.Discard, calls); e != nil {
			err = e
		} else {
			ess.model = metal.NewMetal(m, make(metal.ObjectValueMap))
			ess.calls = calls
			ret = ess.model
		}
	}
	return
}
