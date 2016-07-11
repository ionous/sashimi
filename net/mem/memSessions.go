package mem

import (
	"github.com/ionous/sashimi/compiler/call"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/ess"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/standard/framework"
	"github.com/ionous/sashimi/util/ident" // for generating unique ids
	"github.com/ionous/sashimi/util/uuid"
	"io/ioutil"
	"sync"
)

type essMap map[string]MemSession

// ess.SessionFactory
type MemSessions struct {
	sessions essMap
	model    *M.Model
	calls    call.MarkerStorage
	*sync.Mutex
}

type MemSession struct {
	*app.PartialSession // IResource
	*sync.RWMutex       // RLock, RUnlock, Lock, Unlock
}

func NewSessions() *MemSessions {
	return &MemSessions{sessions: make(essMap), Mutex: new(sync.Mutex)}
}

func (ess *MemSessions) NewSession(doc resource.DocumentBuilder) (ret ess.Session, err error) {
	if e := ess.compile(); e != nil {
		err = e
	} else {
		id := ident.Dash(uuid.MakeUniqueId())
		// FIX? load?
		meta := metal.NewMetal(ess.model, make(metal.ObjectValueMap))
		out := app.NewCommandOutput(id, meta, framework.NewStandardView(meta))
		if s, e := app.NewPartialSession(meta, ess.calls, out); e != nil {
			err = e
		} else {
			out.FlushDocument(doc)
			//
			defer ess.Unlock()
			mem := MemSession{s, new(sync.RWMutex)}
			ess.Lock()
			ess.sessions[id] = mem
			ret = mem
		}
	}
	return
}

func (ess MemSessions) GetSession(id string) (ret ess.Session, okay bool) {
	defer ess.Unlock()
	ess.Lock()
	ret, okay = ess.sessions[id]
	return
}

func (ess *MemSessions) compile() (err error) {
	if ess.model == nil {
		calls := call.MakeMarkerStorage()
		if m, e := script.InitScripts().CompileCalls(ioutil.Discard, calls); e != nil {
			err = e
		} else {
			ess.model, ess.calls = m, calls
		}
	}
	return
}
