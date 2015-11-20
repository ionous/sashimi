package datastore

import (
	D "appengine/datastore"
	"github.com/ionous/sashimi/meta"
)

type KeyGen interface {
	NewKey(meta.Instance) *D.Key
}
