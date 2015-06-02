package model

import "encoding/json"

type RelativeProperty struct {
	fields RelativeFields
}

type RelativeFields struct {
	Id       StringId `json:"id"`       // property id
	Name     string   `json:"name"`     // property name
	Relation StringId `json:"relation"` // relation id
	Relates  StringId `json:"relates"`  // other class id
	IsRev    bool     `json:"rev"`
	ToMany   bool     `json:"many"`
}

func NewRelative(fields RelativeFields) *RelativeProperty {
	return &RelativeProperty{fields}
}

func (this *RelativeProperty) Id() StringId {
	return this.fields.Id
}

func (this *RelativeProperty) Name() string {
	return this.fields.Name
}

func (this *RelativeProperty) Relation() StringId {
	return this.fields.Relation
}

// returns the relation, and true if it exited in the map
func (this *RelativeProperty) FindRelation(relations RelationMap) (Relation, bool) {
	r, ok := relations[this.fields.Relation]
	return r, ok
}

// other class id
func (this *RelativeProperty) Relates() StringId {
	return this.fields.Relates
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the primary.
func (this *RelativeProperty) IsRev() bool {
	return this.fields.IsRev
}

// in the case where a one to many relation involves a class and itself,
// distingushes which side of the relation propery is the many.
func (this *RelativeProperty) ToMany() bool {
	return this.fields.ToMany
}

func (this *RelativeProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(&this.fields)
}
