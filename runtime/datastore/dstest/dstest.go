package dstest

import (
	A "appengine"
	D "appengine/datastore"
	"github.com/ionous/sashimi/meta"
	"github.com/ionous/sashimi/util/ident"
)

// MockKeyGen implements sashimi/datastore/KeyGen.
type MockKeyGen struct {
	c      A.Context
	parent *D.Key
}

func NewKeyGen(c A.Context, parent *D.Key) MockKeyGen {
	return MockKeyGen{c, parent}
}

func (f MockKeyGen) NewKey(inst meta.Instance) *D.Key {
	// maybe....
	// h:= fnv.New64()
	// io.WriteString(h, inst.GetId())
	kind := ident.Dash(inst.GetParentClass().GetId())
	stringID := ident.Dash(inst.GetId())
	intID := int64(0) //h.Sum64()
	return D.NewKey(f.c, kind, stringID, intID, f.parent)
}
