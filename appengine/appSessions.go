package appengine

import (
	"appengine"
	"appengine/datastore"
	DS "github.com/ionous/sashimi/appengine/datastore"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/net/app"
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
	key, out := aps.newSessionKey(id), app.NewCommandOutput(id, app.NewObjectSerializer(AlwaysKnown{}))
	ret, err = aps.newSession(key, out)
	out.FlushDocument(doc)
	return
}

func (aps AppSessions) GetSession(id string) (ret ess.ISession, okay bool) {
	key, out := aps.newSessionKey(id), app.NewCommandOutput(id, app.NewObjectSerializer(AlwaysKnown{}))
	if s, e := aps.newSession(key, out); e == nil {
		ret, okay = s, true
	}
	return
}

func (aps AppSessions) newSession(key *datastore.Key, out *app.CommandOutput) (ret AppSession, err error) {
	// the model store stores changes to the passed model under the passed ancestor session key.
	ds := DS.NewModelStore(aps.ctx, aps.model, key)
	return NewAppSession(out, ds, aps.calls)
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
