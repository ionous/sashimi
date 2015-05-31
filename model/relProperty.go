package model

type RelativeProperty struct {
	id       StringId
	name     string
	relation StringId // map property name to relation id
	relates  StringId // map to other class id
	isRev    bool
	isMany   bool
}

func NewRelative(id StringId, name string, relation StringId, relates StringId, isRev bool, isMany bool) *RelativeProperty {
	return &RelativeProperty{id, name, relation, relates, isRev, isMany}
}

func (this *RelativeProperty) Id() StringId {
	return this.id
}

func (this *RelativeProperty) Name() string {
	return this.name
}

// returns the relation, and true if it exited in the map
func (this *RelativeProperty) Relation(relations RelationMap) (Relation, bool) {
	r, ok := relations[this.relation]
	return r, ok
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the primary.
func (this *RelativeProperty) IsRev() bool {
	return this.isRev
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the many.
func (this *RelativeProperty) ToMany() bool {
	return this.isMany
}

// other class id
func (this *RelativeProperty) Relates() StringId {
	return this.relates
}
