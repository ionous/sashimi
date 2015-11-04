package runtime

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/model/table"
	"github.com/ionous/sashimi/util/ident"
)

// interface to access instance data.
// FIX: may need to split into hasone,hasmany subtypes
type RelativeValue struct {
	inst  ident.Id
	prop  M.RelativeProperty
	table *table.Table
}

// return a list of referenced instances
func (rel RelativeValue) List() (ret []ident.Id) {
	return rel.table.List(rel.inst, rel.prop.IsRev)
}

// FIX: where and how to validate table.style?
// returns item removed
func (rel RelativeValue) ClearReference() (ret ident.Id, err error) {
	if rel.prop.IsMany {
		err = fmt.Errorf("setting an object, but relation is a list")
	} else {
		// FIX: some sort of early return.
		isRev := rel.prop.IsRev
		rel.table.Remove(func(x, y ident.Id) (removed bool) {
			if !isRev && rel.inst == x {
				ret = y
				removed = true
			} else if rel.prop.IsRev && rel.inst == y {
				ret = x
				removed = true
			}
			return removed
		})
	}
	return ret, err
}

// FIX: table.style: where and how to validate style?
// returns previous value
func (rel RelativeValue) SetReference(other *M.InstanceInfo) (removed ident.Id, err error) {
	if !other.Class.CompatibleWith(rel.prop.Relates) {
		err = fmt.Errorf("%s not compatible with %+v", other, rel.prop)
	} else {
		if !rel.prop.IsMany {
			removed, err = rel.ClearReference()
		}
		if err == nil {
			src, dst := rel.inst, other.Id
			if rel.prop.IsRev {
				dst, src = src, dst
			}
			rel.table.Add(src, dst)
		}
	}
	return removed, err
}
