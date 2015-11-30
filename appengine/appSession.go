package appengine

import (
	A "appengine"
	D "appengine/datastore"
	DS "github.com/ionous/sashimi/appengine/datastore"
	"github.com/ionous/sashimi/net/app"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"io"
)

type AppSession struct {
	ctx A.Context
	ds  *DS.ModelStore
	*app.PartialSession
	key *D.Key // key of the session, used as the parent of all other queries
}

func NewAppSession(ctx A.Context, id string, ds *DS.ModelStore, calls api.LookupCallbacks) (ret AppSession, err error) {
	// FIX: you might consider parenting this to the app and data context
	// only one right now , but itd be nice to support multiple versions
	// multiple stories
	var parent *D.Key
	var stringId string = id
	var intId int64
	var kind string = "sessions"

	out := app.NewCommandOutput(id, app.NewObjectSerializer(AlwaysKnown{}))
	if partial, e := app.NewPartialSession(out, ds.Model(), calls); e != nil {
		err = e
	} else {
		ret = AppSession{ctx, ds, partial, D.NewKey(ctx, kind, stringId, intId, parent)}
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
