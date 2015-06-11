package compiler

import (
	"fmt"
	M "github.com/ionous/sashimi/model"
)

//
// A property pointing to some other class.
// FIX? this is almost an exact duplicate of RelativeFields
//
type PendingRelative struct {
	class       *PendingClass
	id          M.StringId    // property id in class
	name        string        // original property name
	relatesTo   *PendingClass // the other side of the relation
	viaRelation string        // the name of the relation, which becomes an id.
	toMany      bool
	isRev       bool
}

//
//
//
func (this PendingRelative) String() string {
	return fmt.Sprintf("%s.%s", this.class, this.id)
}
