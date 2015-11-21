package datastore

import (
	A "appengine"
	D "appengine/datastore"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// SimpleKeyGen implements sashimi/datastore/KeyGen.
type SimpleKeyGen struct {
	c      A.Context
	parent *D.Key
}

func NewKeyGen(c A.Context, parent *D.Key) SimpleKeyGen {
	return SimpleKeyGen{c, parent}
}

func (f SimpleKeyGen) NewKey(inst meta.Instance) *D.Key {
	kind := ident.Dash(inst.GetParentClass().GetId())
	stringID := ident.Dash(inst.GetId())
	intID := int64(0)
	return D.NewKey(f.c, kind, stringID, intID, f.parent)
}
