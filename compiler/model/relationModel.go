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
	Id             ident.Id // unique id
	Name           string   // user specified name
	Style          RelationStyle
	Source, Target ident.Id // property ids
}

func (r RelationModel) String() string {
	return r.Name
}

func (r RelationModel) GetOther(fromProp ident.Id) (other ident.Id) {
	if fromProp == r.Source {
		other = r.Target
	} else {
		other = r.Source
	}
	return
}
