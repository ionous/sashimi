package appengine

import (
	"appengine"
	"appengine/datastore"
	DS "github.com/ionous/sashimi/appengine/datastore"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/net/ess"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

type AppSessions struct {
	ctx   appengine.Context
	model *M.Model
	calls api.LookupCallbacks
}

func NewSessions(
	ctx appengine.Context,
	model *M.Model,
	calls api.LookupCallbacks,
) AppSessions {
	return AppSessions{ctx, model, calls}
}

func (aps AppSessions) NewSession(doc resource.DocumentBuilder) (ret ess.ISession, err error) {
	id := ident.Dash(ident.MakeUniqueId())
	if s, e := aps.newSession(id); e != nil {
		err = e
	} else if e := s.FlushDocument(doc); e != nil {
		err = e
	} else {
		ret = s
	}
	return
}

func (aps AppSessions) GetSession(id string) (ret ess.ISession, okay bool) {
	if s, e := aps.newSession(id); e == nil {
		ret, okay = s, true
	}
	return
}

func (aps AppSessions) newSession(id string) (ret AppSession, err error) {
	// the model store stores changes to the passed model under the passed ancestor session key.
	ds := DS.NewModelStore(aps.ctx, aps.model, aps.newSessionKey(id))
	return NewAppSession(id, ds, aps.calls)
}

func (aps AppSessions) newSessionKey(id string) *datastore.Key {
	// FIX: you might consider parenting this to the app and data context
	// only one right now , but itd be nice to support multiple versions
	// multiple stories
	var kind string = "sessions"
	var stringID string = id
	var intID int64
	var parent *datastore.Key
	return datastore.NewKey(aps.ctx, kind, stringID, intID, parent)
}
