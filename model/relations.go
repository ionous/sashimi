package model

import "github.com/ionous/sashimi/util/ident"

type RelationStyle string

const (
	OneToOne   RelationStyle = "OneToOne"
	OneToMany                = "OneToMany"
	ManyToOne                = "ManyToOne"
	ManyToMany               = "ManyToMany"
)

func NewRelation(id ident.Id, name string, src, dst HalfRelation, style RelationStyle) Relation {
	return Relation{id, name, src, dst, style}
}

// Relation represents a property-pair.
// Currently, each relation becomes one table.
// This might always be the case, but it's also possible to imagine many property views of the same table.
type Relation struct {
	Id     ident.Id // unique id
	Name   string   // user specified name
	Source HalfRelation
	Dest   HalfRelation
	Style  RelationStyle
}

type HalfRelation struct {
	Class    ident.Id
	Property ident.Id
}

type RelationMap map[ident.Id]Relation

func (r Relation) GetOther(isRev bool) (other HalfRelation) {
	if isRev {
		other = r.Source
	} else {
		other = r.Dest
	}
	return
}
