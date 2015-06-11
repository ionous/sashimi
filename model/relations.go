package model

type RelationStyle string

const (
	OneToOne   RelationStyle = "OneToOne"
	OneToMany                = "OneToMany"
	ManyToOne                = "ManyToOne"
	ManyToMany               = "ManyToMany"
)

func NewRelation(id StringId, name string, src, dst HalfRelation, style RelationStyle) Relation {
	return Relation{id, name, src, dst, style}
}

//
// A relation represents a property-pair.
// Currently, each relation becomes one table.
// This might always be the case, but it's also possible to imagine many property views of the same table.
//
type Relation struct {
	id    StringId // unique id
	name  string   // user specified name
	src   HalfRelation
	dst   HalfRelation
	style RelationStyle
}

type HalfRelation struct {
	Class    StringId
	Property StringId
}

type RelationMap map[StringId]Relation

func (this Relation) Id() StringId {
	return this.id
}

func (this Relation) Name() string {
	return this.name
}

func (this Relation) Style() RelationStyle {
	return this.style
}

func (this Relation) Source() HalfRelation {
	return this.src
}

func (this Relation) Destination() HalfRelation {
	return this.dst
}

func (this Relation) Other(class StringId, property StringId) (other HalfRelation, okay bool) {
	relative := HalfRelation{class, property}
	if relative == this.src {
		other = this.dst
		okay = true
	} else if relative == this.dst {
		other = this.src
		okay = true
	}
	return
}
