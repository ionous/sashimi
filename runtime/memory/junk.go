package memory

import (
	M "github.com/ionous/sashimi/model"
	"github.com/ionous/sashimi/util/ident"
)

// single/plural properties need some fixing
var singular = ident.MakeId("singular")
var plural = ident.MakeId("plural")

// junkProperty imitates properties for "single" and "plural" names in a class
type junkProperty struct {
	id  ident.Id
	val string
}

func (p junkProperty) GetId() ident.Id {
	return p.id
}

func (p junkProperty) GetName() string {
	return p.id.String()
}

func (p junkProperty) GetZero(_ M.ConstraintSet) interface{} {
	return p.val
}
