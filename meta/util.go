package meta

import "github.com/ionous/sashimi/util/ident"

func FindFirstOf(mdl Model, cls ident.Id) (ret Instance, okay bool) {
	for i := 0; i < mdl.NumInstance(); i++ {
		inst := mdl.InstanceNum(i)
		if inst.GetParentClass() == cls {
			ret, okay = inst, true
			break
		}
	}
	return
}
