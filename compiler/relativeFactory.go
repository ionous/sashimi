package compiler

import (
	M "github.com/ionous/sashimi/model"
)

//
// Creates relatives while building a picture of their relations.
//
type RelativeFactory struct {
	names     NameScope
	relations PendingRelations
}

type PendingRelations map[M.StringId]PendingRelation

func newRelativeFactory(names NameScope) *RelativeFactory {
	return &RelativeFactory{names, make(PendingRelations)}
}

//x
// Make a relative property by generating (or updating) a shared relation.
//
func (this *RelativeFactory) makeRelative(relative PendingRelative,
) (ret *M.RelativeProperty,
	err error,
) {
	if relId, e := this.names.addName(relative.viaRelation, "relation"); e != nil {
		err = e
	} else {
		// merge this relative with it's counter-part ( if any so far )
		rel := this.relations[relId]

		// FIX? the relation has split into two class relative properties
		// and now were merging them back together, verifying they match
		// it might have been better to keep the pair,
		// splitting them into their class haves at class creation time
		if e := rel.setRelative(relative); e != nil {
			err = e
		} else {
			// write merged data back
			this.relations[relId] = rel

			// create the relative property pointing to the generated relation data
			fields := M.RelativeFields{
				relative.class.id,
				relative.id,
				relative.name,
				relative.relatesTo.id,
				relId,
				relative.isRev,
				relative.toMany}
			ret = M.NewRelative(fields)
		}
	}
	return ret, err
}
