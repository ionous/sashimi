package play

import "github.com/ionous/sashimi/meta"

type noParents struct{}

func (noParents) LookupParent(meta.Instance) (inst meta.Instance, rel meta.Property, okay bool) {
	return
}
