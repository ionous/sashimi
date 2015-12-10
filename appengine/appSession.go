package appengine

import (
	//"appengine"
	DS "github.com/ionous/sashimi/appengine/datastore"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/standard/framework"
	"io"
)

type AppSession struct {
	//ctx appengine.Context
	ds *DS.ModelStore
	*app.PartialSession
}

func NewAppSession(
	id string,
	ds *DS.ModelStore,
	calls api.LookupCallbacks,
) (
	ret AppSession, err error,
) {
	mdl := ds.Model()
	out := app.NewCommandOutput(id, mdl, framework.NewStandardView(mdl))
	if partial, e := app.NewPartialSession(mdl, calls, out); e != nil {
		err = e
	} else {
		ret = AppSession{ds, partial}
	}
	return
}

// FIX? dont have to worry about memory stomping, but may have to worrry about consistency of data.
func (AppSession) RLock()   {}
func (AppSession) RUnlock() {}
func (AppSession) Lock()    {}
func (AppSession) Unlock()  {}

func (s AppSession) Post(reader io.Reader) (ret resource.Document, err error) {
	if d, e := s.PartialSession.Post(reader); e != nil {
		err = e
	} else if e := s.ds.Flush(); e != nil {
		err = e
	} else {
		ret = d
	}
	return
}
