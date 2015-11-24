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
	Id     ident.Id      `json:"id"`   // unique id
	Name   string        `json:"name"` // user specified name
	Style  RelationStyle `json:"style"`
	Source ident.Id      `json:"src"` // property ids
	Target ident.Id      `json:"tgt"` // property ids
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
