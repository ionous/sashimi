package model

type RelativeProperty struct {
	fields RelativeFields
}

type RelativeFields struct {
	Class    StringId `json:"id"`       // property id
	Id       StringId `json:"id"`       // property id
	Name     string   `json:"name"`     // property name
	Relates  StringId `json:"relates"`  // other class id
	Relation StringId `json:"relation"` // relation id
	IsRev    bool     `json:"rev"`
	ToMany   bool     `json:"many"`
}

func NewRelativeProperty(fields RelativeFields) *RelativeProperty {
	return &RelativeProperty{fields}
}

func (this *RelativeProperty) Class() StringId {
	return this.fields.Class
}

func (this *RelativeProperty) Id() StringId {
	return this.fields.Id
}

func (this *RelativeProperty) Name() string {
	return this.fields.Name
}

// id of the relation table
func (this *RelativeProperty) Relation() StringId {
	return this.fields.Relation
}

// id of the other class this property inolves
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
