package appengine

import (
	A "appengine"
	DS "github.com/ionous/sashimi/appengine/datastore"
	M "github.com/ionous/sashimi/compiler/model"
	"github.com/ionous/sashimi/net/ess"
	"github.com/ionous/sashimi/net/resource"
	"github.com/ionous/sashimi/runtime/api"
	"github.com/ionous/sashimi/util/ident"
)

// data, and code constants: model *M.Model, calls api.LookupCallbacks
// net.HandleResource(app.GameResource(mem.NewMemSessions())))

type AppSessions struct {
	ctx   A.Context
	ds    *DS.ModelStore
	calls api.LookupCallbacks
}

func NewSessions(ctx A.Context,
	model *M.Model,
	calls api.LookupCallbacks,
) AppSessions {
	ds := DS.NewModelStore(ctx, model)
	return AppSessions{ctx, ds, calls}
}

func (aps AppSessions) NewSession(doc resource.DocumentBuilder) (ret ess.ISessionResource, err error) {
	id := ident.Dash(ident.MakeUniqueId())
	return aps.newSession(id)
}

func (aps AppSessions) GetSession(id string) (ret ess.ISessionResource, okay bool) {
	if s, e := aps.newSession(id); e == nil {
		ret, okay = s, true
	}
	return
}

func (aps AppSessions) newSession(id string) (ret AppSession, err error) {
	return NewAppSession(aps.ctx, id, aps.ds, aps.calls)
}
