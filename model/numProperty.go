package model

type NumProperty struct {
	id   StringId
	name string
}

func NewNumProperty(id StringId, name string) *NumProperty {
	return &NumProperty{id, name}
}

func (this *NumProperty) Id() StringId {
	return this.id
}

func (this *NumProperty) Name() string {
	return this.name
}
