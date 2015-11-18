package model

import "github.com/ionous/sashimi/util/ident"

type RelationStyle string

const (
	OneToOne   RelationStyle = "OneToOne"
	OneToMany                = "OneToMany"
	ManyToOne                = "ManyToOne"
	ManyToMany               = "ManyToMany"
)

// Relation represents a property-pair.
// Each relation becomes one table.
type RelationModel struct {
	Id   ident.Id // unique id
	Name string   // user specified name
	// Type might be the relation by value data....
	Source HalfRelation
	Dest   HalfRelation
	Style  RelationStyle
}

// Half relation exists for backwards compat with client
// see also: GetRelative()
type HalfRelation struct {
	Class    ident.Id
	Property ident.Id
}

func (r RelationModel) String() string {
	return r.Name
}

func (r RelationModel) GetOther(isRev bool) (other HalfRelation) {
	if isRev {
		other = r.Source
	} else {
		other = r.Dest
	}
	return
}
