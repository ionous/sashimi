package mem

import (
	"github.com/ionous/sashimi/compiler/call"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/metal"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/ess"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/script"
	"github.com/ionous/sashimi/util/ident" // for generating unique ids
	"io/ioutil"
	"sync"
)

type essMap map[string]MemSession

// ess.ISessionResourceFactory
type MemSessions struct {
	sessions essMap
	model    meta.Model
	calls    call.MarkerStorage
	*sync.Mutex
}

type MemSession struct {
	*app.PartialSession // IResource
	*sync.RWMutex       // RLock, RUnlock, Lock, Unlock
}

func NewMemSessions() *MemSessions {
	return &MemSessions{sessions: make(essMap), Mutex: new(sync.Mutex)}
}

func (ess *MemSessions) NewSession(doc resource.DocumentBuilder) (ret ess.ISessionResource, err error) {
	if m, e := ess.getModel(); e != nil {
		err = e
	} else {
		id := ident.Dash(ident.MakeUniqueId())
		out := app.NewCommandOutput(id, app.NewObjectSerializer(make(app.KnownObjectMap)))
		if s, e := app.NewPartialSession(out, m, ess.calls); e != nil {
			err = e
		} else {
			out.FlushDocument(doc)
			mem := MemSession{s, new(sync.RWMutex)}
			//
			defer ess.Unlock()
			ess.Lock()
			ess.sessions[id] = mem
			ret = mem
		}
	}
	return
}

func (ess MemSessions) GetSession(id string) (ret ess.ISessionResource, okay bool) {
	defer ess.Unlock()
	ess.Lock()
	ret, okay = ess.sessions[id]
	return
}

func (ess *MemSessions) getModel() (ret meta.Model, err error) {
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
